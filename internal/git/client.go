package git

import (
	"bump/internal/config"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
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
		// Validate tag format to prevent command injection
		if !isValidGitTag(tag) {
			return nil, fmt.Errorf("invalid git tag format: %s", tag)
		}
		// Use git log with explicit revision range
		// Input is validated by isValidGitTag() to prevent command injection
		revRange := tag + "..HEAD"
		cmd = exec.Command("git", "log", "--oneline", revRange) // #nosec G204
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

func (g *Client) CheckoutBranch(branch string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would checkout branch: %s\n", branch)
		return nil
	}

	cmd := exec.Command("git", "checkout", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
	}

	return nil
}

func (g *Client) GetDefaultBranch() (string, error) {
	// Try to get the default branch from remote origin
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(output)), "/")
		if len(parts) > 1 {
			return parts[1], nil
		}
	}

	// Fallback: check for common default branch names
	branches := []string{"main", "master"}
	for _, branch := range branches {
		if g.BranchExists(branch) {
			return branch, nil
		}
	}

	return "main", nil
}

func (g *Client) BranchExists(branch string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	err := cmd.Run()
	return err == nil
}

func (g *Client) CreateBranch(branch, sourceBranch string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would create branch: %s from %s\n", branch, sourceBranch)
		return nil
	}

	// Checkout source branch first
	cmd := exec.Command("git", "checkout", sourceBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout source branch %s: %w", sourceBranch, err)
	}

	// Create and checkout new branch
	cmd = exec.Command("git", "checkout", "-b", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %s: %w", branch, err)
	}

	return nil
}

func (g *Client) MergeBranch(sourceBranch, targetBranch string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would merge branch: %s into %s\n", sourceBranch, targetBranch)
		return nil
	}

	// Checkout target branch
	cmd := exec.Command("git", "checkout", targetBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout target branch %s: %w", targetBranch, err)
	}

	// Merge source branch
	cmd = exec.Command("git", "merge", sourceBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge branch %s into %s: %w", sourceBranch, targetBranch, err)
	}

	return nil
}

func (g *Client) PushBranch(branch string) error {
	if g.cfg.DryRun {
		fmt.Printf("[DRY RUN] Would push branch: %s\n", branch)
		return nil
	}

	cmd := exec.Command("git", "push", "origin", branch)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push branch %s: %w", branch, err)
	}

	return nil
}

func (g *Client) GetAllTags() ([]string, error) {
	cmd := exec.Command("git", "for-each-ref", "--sort=-creatordate", "--format=%(refname:short) %(creatordate:iso)", "refs/tags")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	return lines, nil
}

// isValidGitTag validates that a git tag contains only safe characters
// to prevent command injection attacks
func isValidGitTag(tag string) bool {
	// Git tags should only contain alphanumeric, dots, hyphens, underscores, and forward slashes
	// This regex allows semantic versioning tags like v1.2.3, v1.2.3-alpha, etc.
	validTagPattern := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)

	// Additional length check to prevent excessively long inputs
	if len(tag) > 100 {
		return false
	}

	// Check for dangerous characters that could be used for command injection
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]", "<", ">", "\\", "\"", "'", " "}
	for _, char := range dangerousChars {
		if strings.Contains(tag, char) {
			return false
		}
	}

	return validTagPattern.MatchString(tag)
}
