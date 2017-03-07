package db

type Repository struct {
	Model

	Cloned bool

	Name          string
	Owner         string
	Path          string
	SSHURL        string `gorm:"column:ssh_url"`
	Private       bool
	DefaultBranch string
	RawJSON       []byte `gorm:"column:raw_json"`

	FailedFetches int `gorm:"column:failed_fetches"`
	Disabled      bool
}
