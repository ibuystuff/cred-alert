package webhook_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cred-alert/github/githubfakes"
	"cred-alert/metrics"
	"cred-alert/metrics/metricsfakes"
	"cred-alert/notifications/notificationsfakes"
	"cred-alert/sniff"
	"cred-alert/webhook"

	"github.com/google/go-github/github"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Extract", func() {
	var (
		logger *lagertest.TestLogger
		event  github.PushEvent
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("extract")

		event = github.PushEvent{
			Before: github.String("abc123bef04e"),
			Repo: &github.PushEventRepository{
				Name: github.String("repository-name"),
				Owner: &github.PushEventRepoOwner{
					Name: github.String("repository-owner"),
				},
			},
			Commits: []github.PushEventCommit{
				{ID: github.String("commit-sha-1")},
				{ID: github.String("commit-sha-2")},
				{ID: github.String("commit-sha-3")},
				{ID: github.String("commit-sha-4")},
				{ID: github.String("commit-sha-5")},
			},
		}
	})

	It("can extract a value object from a github push event", func() {
		scan, valid := webhook.Extract(logger, event)
		Expect(valid).To(BeTrue())

		Expect(scan.Owner).To(Equal("repository-owner"))
		Expect(scan.Repository).To(Equal("repository-name"))
		Expect(scan.Diffs).To(Equal([]webhook.PushScanDiff{
			{Start: "abc123bef04e", End: "commit-sha-1"},
			{Start: "commit-sha-1", End: "commit-sha-2"},
			{Start: "commit-sha-2", End: "commit-sha-3"},
			{Start: "commit-sha-3", End: "commit-sha-4"},
			{Start: "commit-sha-4", End: "commit-sha-5"},
		}))
	})

	It("can have a full repository name", func() {
		scan, valid := webhook.Extract(logger, event)
		Expect(valid).To(BeTrue())

		Expect(scan.Owner).To(Equal("repository-owner"))
		Expect(scan.Repository).To(Equal("repository-name"))

		Expect(scan.FullRepoName()).To(Equal("repository-owner/repository-name"))
	})

	It("can handle if there are no commits in a push (may not even be possible)", func() {
		event.Commits = []github.PushEventCommit{}

		_, valid := webhook.Extract(logger, event)
		Expect(valid).To(BeFalse())
	})
})

var _ = Describe("EventHandler", func() {
	var (
		eventHandler     webhook.EventHandler
		logger           *lagertest.TestLogger
		emitter          *metricsfakes.FakeEmitter
		notifier         *notificationsfakes.FakeNotifier
		fakeGithubClient *githubfakes.FakeClient

		orgName      string
		repoName     string
		repoFullName string

		sniffFunc func(lager.Logger, sniff.Scanner, func(sniff.Line))

		requestCounter      *metricsfakes.FakeCounter
		credentialCounter   *metricsfakes.FakeCounter
		ignoredEventCounter *metricsfakes.FakeCounter

		whitelist *webhook.Whitelist
		event     github.PushEvent
	)

	BeforeEach(func() {
		orgName = "rad-co"
		repoName = "my-awesome-repo"
		repoFullName = fmt.Sprintf("%s/%s", orgName, repoName)

		sniffFunc = func(lager.Logger, sniff.Scanner, func(sniff.Line)) {}

		emitter = &metricsfakes.FakeEmitter{}
		notifier = &notificationsfakes.FakeNotifier{}
		requestCounter = &metricsfakes.FakeCounter{}
		credentialCounter = &metricsfakes.FakeCounter{}
		ignoredEventCounter = &metricsfakes.FakeCounter{}

		whitelist = webhook.BuildWhitelist()

		emitter.CounterStub = func(name string) metrics.Counter {
			switch name {
			case "cred_alert.webhook_requests":
				return requestCounter
			case "cred_alert.violations":
				return credentialCounter
			case "cred_alert.ignored_events":
				return ignoredEventCounter
			default:
				panic("unexpected counter name! " + name)
			}
		}

		logger = lagertest.NewTestLogger("event-handler")
		fakeGithubClient = new(githubfakes.FakeClient)

		event = github.PushEvent{
			Repo: &github.PushEventRepository{
				FullName: github.String(repoFullName),
				Name:     github.String(repoName),
				Owner: &github.PushEventRepoOwner{
					Name: github.String(orgName),
				},
			},
			Before: github.String("sha0"),
			Commits: []github.PushEventCommit{
				github.PushEventCommit{ID: github.String("def456")},
			},
		}
	})

	JustBeforeEach(func() {
		eventHandler = webhook.NewEventHandler(fakeGithubClient, sniffFunc, emitter, notifier, whitelist)
	})

	Context("when there are multiple commits in a single event", func() {
		var before string = "before"
		var id0, id1, id2 string = "a", "b", "c"

		BeforeEach(func() {
			commit0 := github.PushEventCommit{ID: &id0}
			commit1 := github.PushEventCommit{ID: &id1}
			commit2 := github.PushEventCommit{ID: &id2}
			commits := []github.PushEventCommit{commit0, commit1, commit2}

			event.Before = &before
			event.Commits = commits
		})

		It("compares each commit individually", func() {
			eventHandler.HandleEvent(logger, event)

			fakeGithubClient.CompareRefsReturns("", errors.New("disaster"))
			Expect(fakeGithubClient.CompareRefsCallCount()).To(Equal(3))
			_, _, _, sha0, sha1 := fakeGithubClient.CompareRefsArgsForCall(0)
			Expect(sha0).To(Equal(before))
			Expect(sha1).To(Equal(id0))
			_, _, _, sha0, sha1 = fakeGithubClient.CompareRefsArgsForCall(1)
			Expect(sha0).To(Equal(id0))
			Expect(sha1).To(Equal(id1))
			_, _, _, sha0, sha1 = fakeGithubClient.CompareRefsArgsForCall(2)
			Expect(sha0).To(Equal(id1))
			Expect(sha1).To(Equal(id2))
		})
	})

	It("emits count when it is invoked", func() {
		eventHandler.HandleEvent(logger, event)

		Expect(requestCounter.IncCallCount()).To(Equal(1))
	})

	Context("It has a whitelist of ignored repos", func() {
		var scanCount int

		BeforeEach(func() {
			repoName = "some-credentials"

			scanCount = 0
			sniffFunc = func(lager.Logger, sniff.Scanner, func(sniff.Line)) {
				scanCount++
			}
			whitelist = webhook.BuildWhitelist(repoName)
			event.Repo.Name = &repoName
		})

		It("ignores patterns in whitelist", func() {
			eventHandler.HandleEvent(logger, event)

			Expect(scanCount).To(BeZero())
			Expect(len(logger.LogMessages())).To(Equal(1))
			Expect(logger.LogMessages()[0]).To(ContainSubstring("ignored-repo"))
			Expect(logger.Logs()[0].Data["repo"]).To(Equal(repoName))
		})

		It("emits a count of ignored push events", func() {
			eventHandler.HandleEvent(logger, event)
			Expect(ignoredEventCounter.IncCallCount()).To(Equal(1))
		})
	})

	Context("when a credential is found", func() {
		var filePath string
		var sha0 string = "sha0"

		BeforeEach(func() {
			filePath = "some/file/path"

			sniffFunc = func(logger lager.Logger, scanner sniff.Scanner, handleViolation func(sniff.Line)) {
				handleViolation(sniff.Line{
					Path:       filePath,
					LineNumber: 1,
					Content:    "content",
				})
			}

			event.Commits[0].ID = &sha0
		})

		It("emits count of the credentials it has found", func() {
			eventHandler.HandleEvent(logger, event)

			Expect(credentialCounter.IncNCallCount()).To(Equal(1))
		})

		It("sends a notification", func() {
			eventHandler.HandleEvent(logger, event)

			Expect(notifier.SendNotificationCallCount()).To(Equal(1))

			_, repo, sha, line := notifier.SendNotificationArgsForCall(0)

			Expect(repo).To(Equal(repoFullName))
			Expect(sha).To(Equal(sha0))
			Expect(line).To(Equal(sniff.Line{
				Path:       "some/file/path",
				LineNumber: 1,
				Content:    "content",
			}))
		})
	})

	Context("when we fail to fetch the diff", func() {
		var wasScanned bool

		BeforeEach(func() {
			wasScanned = false

			fakeGithubClient.CompareRefsReturns("", errors.New("disaster"))

			sniffFunc = func(lager.Logger, sniff.Scanner, func(sniff.Line)) {
				wasScanned = true
			}
		})

		It("does not try to scan the diff", func() {
			eventHandler.HandleEvent(logger, event)

			Expect(wasScanned).To(BeFalse())
			Expect(credentialCounter.IncNCallCount()).To(Equal(0))
		})
	})
})
