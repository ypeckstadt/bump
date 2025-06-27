package config

type Config struct {
	DryRun       bool
	Verbose      bool
	NoBranch     bool
	CreateBranch bool
	SourceBranch string
	BranchName   string
	AutoMerge    bool
	AutoPush     bool
}

func New() *Config {
	return &Config{
		DryRun:       false,
		Verbose:      false,
		NoBranch:     false,
		CreateBranch: false,
		SourceBranch: "",
		BranchName:   "",
		AutoMerge:    false,
		AutoPush:     false,
	}
}
