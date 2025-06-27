package bump

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"bump/internal/config"
	"bump/internal/git"
	"bump/internal/version"
)

type Release struct {
	cfg     *config.Config
	git     *git.Client
	version *version.Version
}

func NewRelease(cfg *config.Config) *Release {
	gitClient := git.NewClient(cfg)
	currentVersionStr := GetCurrentVersion()
	ver := version.NewFromString(currentVersionStr)
	
	return &Release{
		cfg:     cfg,
		git:     gitClient,
		version: ver,
	}
}

func (r *Release) RunInteractive() error {
	printInfo("🚀 Interactive Release Mode")
	
	if !r.git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}
	
	clean, err := r.git.IsWorkingDirectoryClean()
	if err != nil {
		return fmt.Errorf("failed to check working directory: %w", err)
	}
	
	if !clean {
		printWarning("⚠️  Working directory is not clean")
		if !r.confirmProceed("Continue anyway?") {
			return fmt.Errorf("release cancelled")
		}
	}
	
	printInfo(fmt.Sprintf("Current version: %s", r.version.String()))
	
	commits, err := r.git.GetCommitsSinceTag(r.version.Raw)
	if err != nil {
		printWarning("Could not get commits since last tag")
	} else if len(commits) > 0 {
		printInfo("Recent commits:")
		for i, commit := range commits {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(commits)-5)
				break
			}
			fmt.Printf("  %s\n", commit)
		}
	}
	
	versionType, err := r.promptVersionType()
	if err != nil {
		return err
	}
	
	newVersion, err := r.version.Bump(versionType)
	if err != nil {
		return err
	}
	
	printInfo(fmt.Sprintf("New version will be: %s", newVersion.String()))
	
	if r.git.TagExists(newVersion.String()) {
		return fmt.Errorf("tag %s already exists", newVersion.String())
	}
	
	if err := r.runPreReleaseChecks(); err != nil {
		return err
	}
	
	message, err := r.promptReleaseMessage(newVersion.String())
	if err != nil {
		return err
	}
	
	if !r.confirmProceed(fmt.Sprintf("Create and push tag %s?", newVersion.String())) {
		return fmt.Errorf("release cancelled")
	}
	
	return r.createAndPushTag(newVersion.String(), message)
}

func (r *Release) RunQuick(versionType string) error {
	if !r.git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}
	
	newVersion, err := r.version.Bump(versionType)
	if err != nil {
		return err
	}
	
	if r.git.TagExists(newVersion.String()) {
		return fmt.Errorf("tag %s already exists", newVersion.String())
	}
	
	message := fmt.Sprintf("Release %s", newVersion.String())
	
	printInfo(fmt.Sprintf("Creating %s release: %s → %s", versionType, r.version.String(), newVersion.String()))
	
	return r.createAndPushTag(newVersion.String(), message)
}

func (r *Release) promptVersionType() (string, error) {
	patchVersion := r.version.BumpPatch()
	minorVersion := r.version.BumpMinor()
	majorVersion := r.version.BumpMajor()
	
	prompt := promptui.Select{
		Label: "Select version type",
		Items: []string{
			fmt.Sprintf("patch (%s) - bug fixes", patchVersion.String()),
			fmt.Sprintf("minor (%s) - new features", minorVersion.String()),
			fmt.Sprintf("major (%s) - breaking changes", majorVersion.String()),
		},
	}
	
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	
	if strings.Contains(result, "patch") {
		return "patch", nil
	} else if strings.Contains(result, "minor") {
		return "minor", nil
	} else {
		return "major", nil
	}
}

func (r *Release) promptReleaseMessage(version string) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Release message",
		Default: fmt.Sprintf("Release %s", version),
	}
	
	return prompt.Run()
}

func (r *Release) confirmProceed(message string) bool {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	
	result, err := prompt.Run()
	return err == nil && (result == "y" || result == "yes")
}

func (r *Release) runPreReleaseChecks() error {
	printInfo("Running pre-release checks...")
	
	checker := NewChecker(r.cfg)
	if err := checker.RunAll(); err != nil {
		return fmt.Errorf("pre-release checks failed: %w", err)
	}
	
	printSuccess("✅ All checks passed")
	return nil
}

func (r *Release) createAndPushTag(tag, message string) error {
	printInfo(fmt.Sprintf("Creating tag %s...", tag))
	if err := r.git.CreateTag(tag, message); err != nil {
		return err
	}
	
	printInfo(fmt.Sprintf("Pushing tag %s...", tag))
	if err := r.git.PushTag(tag); err != nil {
		return err
	}
	
	printSuccess(fmt.Sprintf("✅ Successfully created and pushed tag %s", tag))
	printInfo("GitHub Actions should now trigger the release workflow")
	return nil
}

func GetCurrentVersion() string {
	cfg := config.New()
	gitClient := git.NewClient(cfg)
	version, err := gitClient.GetLatestTag()
	if err != nil {
		return "v0.0.0"
	}
	return version
}

func printInfo(message string) {
	color.Blue(message)
}

func printSuccess(message string) {
	color.Green(message)
}

func printWarning(message string) {
	color.Yellow(message)
}

func printError(message string) {
	color.Red(message)
}