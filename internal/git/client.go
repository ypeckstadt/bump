package git

import (
	"fmt"
	"os/exec"
	"strings"
	"bump/internal/config"
)

type Client struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (g *Client) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func (g *Client) IsWorkingDirectoryClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	
	return len(strings.TrimSpace(string(output))) == 0, nil
}

func (g *Client) GetLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("no tags found")
	}
	
	return strings.TrimSpace(string(output)), nil
}

func (g *Client) GetCommitsSinceTag(tag string) ([]string, error) {
	var cmd *exec.Cmd
	if tag == "" || tag == "v0.0.0" {
		cmd = exec.Command("git", "log", "--oneline", "-10")
	} else {
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", tag), "--oneline")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	
	return lines, nil
}

func (g *Client) CreateTag(tag, message string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would create tag: %s with message: %s\n", tag, message)
		return nil
	}
	
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag %s: %w", tag, err)
	}
	
	return nil
}

func (g *Client) PushTag(tag string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would push tag: %s\n", tag)
		return nil
	}
	
	cmd := exec.Command("git", "push", "origin", tag)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push tag %s: %w", tag, err)
	}
	
	return nil
}

func (g *Client) TagExists(tag string) bool {
	cmd := exec.Command("git", "tag", "-l", tag)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) == tag
}

func (g *Client) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	
	return strings.TrimSpace(string(output)), nil
}