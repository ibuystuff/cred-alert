package ingestor

import (
	"cred-alert/db"
	"cred-alert/metrics"
	"cred-alert/queue"

	"github.com/pivotal-golang/lager"
)

//go:generate counterfeiter . Ingestor

type Ingestor interface {
	IngestPushScan(lager.Logger, PushScan) error
}

type ingestor struct {
	whitelist        *Whitelist
	taskQueue        queue.Queue
	generator        queue.UUIDGenerator
	commitRepository db.CommitRepository

	requestCounter      metrics.Counter
	ignoredEventCounter metrics.Counter
}

func NewIngestor(taskQueue queue.Queue, emitter metrics.Emitter, whitelist *Whitelist, generator queue.UUIDGenerator, commitRepository db.CommitRepository) *ingestor {
	requestCounter := emitter.Counter("cred_alert.ingestor_requests")
	ignoredEventCounter := emitter.Counter("cred_alert.ignored_events")

	handler := &ingestor{
		taskQueue:           taskQueue,
		whitelist:           whitelist,
		generator:           generator,
		commitRepository:    commitRepository,
		requestCounter:      requestCounter,
		ignoredEventCounter: ignoredEventCounter,
	}

	return handler
}

func (s *ingestor) IngestPushScan(logger lager.Logger, scan PushScan) error {
	logger = logger.Session("handle-event")

	if s.whitelist.IsIgnored(scan.Repository) {
		logger.Info("ignored-repo", lager.Data{
			"repo": scan.Repository,
		})

		s.ignoredEventCounter.Inc(logger)

		return nil
	}

	// Check if from commit is registered, if not queue ref-scan
	repoIsRegistered, err := s.commitRepository.IsRepoRegistered(logger, scan.Owner, scan.Repository)
	if err != nil {
		logger.Error("Error checking for repo: ", err)
		repoIsRegistered = false
	}

	if repoIsRegistered != true {
		id := s.generator.Generate()
		task := queue.RefScanPlan{
			Owner:      scan.Owner,
			Repository: scan.Repository,
			Ref:        scan.Diffs[0].From,
		}.Task(id)

		sessionName := "enqueuing-ref-scan-for-new-repo"

		err := s.taskQueue.Enqueue(task)
		if err != nil {
			logger.Session(sessionName).Error("enqueue-ref-scan-failed", err)
			return err
		}

		logger.Session(sessionName).Info("enqueue-ref-scan-succeeded")
	}

	s.requestCounter.Inc(logger)

	for _, scanDiff := range scan.Diffs {
		id := s.generator.Generate()

		var task queue.Task
		if scanDiff.From == initialCommitParentHash {
			logger = logger.Session("enqueuing-ref-scan", lager.Data{
				"task-id": id,
			})

			task = queue.RefScanPlan{
				Owner:      scan.Owner,
				Repository: scan.Repository,
				Ref:        scanDiff.To,
			}.Task(id)
		} else {
			logger = logger.Session("enqueuing-diff-scan", lager.Data{
				"task-id": id,
			})

			task = queue.DiffScanPlan{
				Owner:      scan.Owner,
				Repository: scan.Repository,
				From:       scanDiff.From,
				To:         scanDiff.To,
			}.Task(id)
		}

		logger.Debug("starting")

		err := s.taskQueue.Enqueue(task)
		if err != nil {
			logger.Error("failed to enqueue scan", err)
			return err
		}

		s.commitRepository.RegisterCommit(logger, &db.Commit{
			Repository: scan.Repository,
			Owner:      scan.Owner,
			SHA:        scanDiff.To,
		})

		logger.Info("done")
	}

	return nil
}
