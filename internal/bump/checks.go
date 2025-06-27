package bump

import (
	"fmt"
	"os/exec"
	"strings"
	"bump/internal/config"
)

type Checker struct {
	cfg *config.Config
}

func NewChecker(cfg *config.Config) *Checker {
	return &Checker{
		cfg: cfg,
	}
}

func (c *Checker) RunAll() error {
	checks := []struct {
		name string
		fn   func() error
	}{
		{"Build", c.checkBuild},
		{"Tests", c.checkTests},
		{"Lint", c.checkLint},
		{"Go mod tidy", c.checkGoModTidy},
	}
	
	for _, check := range checks {
		if c.cfg.Verbose {
			printInfo(fmt.Sprintf("Running %s check...", check.name))
		}
		
		if err := check.fn(); err != nil {
			if c.cfg.Verbose {
				printError(fmt.Sprintf("❌ %s check failed: %v", check.name, err))
			}
			continue
		}
		
		if c.cfg.Verbose {
			printSuccess(fmt.Sprintf("✅ %s check passed", check.name))
		}
	}
	
	return nil
}

func (c *Checker) checkBuild() error {
	if c.cfg.DryRun {
		printInfo("[DRY RUN] Would run: go build")
		return nil
	}
	
	cmd := exec.Command("go", "build", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %s", string(output))
	}
	return nil
}

func (c *Checker) checkTests() error {
	if c.cfg.DryRun {
		printInfo("[DRY RUN] Would run: go test")
		return nil
	}
	
	cmd := exec.Command("go", "test", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "no test files") {
			return nil
		}
		return fmt.Errorf("tests failed: %s", string(output))
	}
	return nil
}

func (c *Checker) checkLint() error {
	if c.cfg.DryRun {
		printInfo("[DRY RUN] Would run: golangci-lint run")
		return nil
	}
	
	cmd := exec.Command("golangci-lint", "run")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return nil
		}
		return fmt.Errorf("lint failed: %s", string(output))
	}
	return nil
}

func (c *Checker) checkGoModTidy() error {
	if c.cfg.DryRun {
		printInfo("[DRY RUN] Would run: go mod tidy")
		return nil
	}
	
	cmd := exec.Command("go", "mod", "tidy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}
	return nil
}