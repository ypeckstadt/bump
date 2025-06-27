package config

type Config struct {
	DryRun  bool
	Verbose bool
}

func New() *Config {
	return &Config{
		DryRun:  false,
		Verbose: false,
	}
}