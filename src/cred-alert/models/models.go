package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pivotal-golang/lager"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DiffScan struct {
	Model
	Org             string
	Repo            string
	FromCommit      string
	ToCommit        string
	Timestamp       time.Time
	TaskID          string
	CredentialFound bool
}

type Commit struct {
	Model
	SHA       string
	Timestamp time.Time
	Org       string
	Repo      string
}

//go:generate counterfeiter . CommitRepository

type CommitRepository interface {
	RegisterCommit(logger lager.Logger, commit *Commit) error
	IsCommitRegistered(logger lager.Logger, sha string) (bool, error)
}

type commitRepository struct {
	db *gorm.DB
}

func NewCommitRepository(db *gorm.DB) *commitRepository {
	return &commitRepository{
		db: db,
	}
}

func (c *commitRepository) RegisterCommit(logger lager.Logger, commit *Commit) error {
	logger = logger.Session("registering-commit", lager.Data{
		"commit-timestamp": commit.Timestamp.Unix(),
		"org":              commit.Org,
		"repo":             commit.Repo,
		"sha":              commit.SHA,
	})

	err := c.db.Save(commit).Error
	if err != nil {
		logger.Error("error-registering-commit", err)
	} else {
		logger.Info("successfully-registered-commit")
	}

	return err
}

func (c *commitRepository) IsCommitRegistered(logger lager.Logger, sha string) (bool, error) {
	logger = logger.Session("finding-commit", lager.Data{
		"sha": sha,
	})

	var commits []Commit
	err := c.db.Where("SHA = ?", sha).First(&commits).Error
	if err != nil {
		logger.Error("error-finding-commit", err)
	}

	return len(commits) == 1, err
}

//go:generate counterfeiter . DiffScanRepository

type DiffScanRepository interface {
	SaveDiffScan(lager.Logger, *DiffScan) error
}

type diffScanRepository struct {
	db *gorm.DB
}

func NewDiffScanRepository(db *gorm.DB) *diffScanRepository {
	return &diffScanRepository{db: db}
}

func (d *diffScanRepository) SaveDiffScan(logger lager.Logger, diffScan *DiffScan) error {
	logger = logger.Session("saving-diffscan", lager.Data{
		"org":              diffScan.Org,
		"repo":             diffScan.Repo,
		"from-commit":      diffScan.FromCommit,
		"to-commit":        diffScan.ToCommit,
		"scan-timestamp":   diffScan.Timestamp.Unix(),
		"task-id":          diffScan.TaskID,
		"credential-found": diffScan.CredentialFound,
	})
	err := d.db.Save(diffScan).Error

	if err != nil {
		logger.Error("error-saving-diffscan", err)
	} else {
		logger.Info("successfully-saved-diffscan")
	}

	return err
}
