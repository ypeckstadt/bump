package main

import (
	"fmt"
	"log"

	"github.com/ypeckstadt/bump/internal/bump"
	"github.com/ypeckstadt/bump/internal/config"
	"github.com/ypeckstadt/bump/pkg/version"

	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

func main() {
	cfg = config.New()

	rootCmd := &cobra.Command{
		Use:   "bump",
		Short: "A version bumping tool for semantic versioning",
		Long: `Bump is a CLI tool for managing semantic versions in git repositories.
It provides both interactive and quick release modes with git integration.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runInteractiveMode()
			} else {
				runQuickMode(args[0])
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "Show what would happen without making changes")
	rootCmd.PersistentFlags().BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&cfg.NoBranch, "nobranch", false, "Skip branch creation prompt entirely")
	rootCmd.PersistentFlags().BoolVar(&cfg.CreateBranch, "create-branch", false, "Create a branch for the tag")
	rootCmd.PersistentFlags().StringVar(&cfg.SourceBranch, "source-branch", "", "Source branch for creating the new branch (default: main/master)")
	rootCmd.PersistentFlags().StringVar(&cfg.BranchName, "branch-name", "", "Name for the new branch (default: tag name without 'v' prefix)")
	rootCmd.PersistentFlags().BoolVar(&cfg.AutoMerge, "auto-merge", false, "Automatically merge if branch exists")
	rootCmd.PersistentFlags().BoolVar(&cfg.AutoPush, "auto-push", false, "Automatically push the branch")

	// Add standard --version flag for CI compatibility
	var showVersion bool
	rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "Show version and exit")

	// Override the default run function to handle --version flag
	originalRun := rootCmd.Run
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if showVersion {
			buildInfo := version.Get()
			fmt.Printf("%s\n", buildInfo.Version)
			return
		}
		originalRun(cmd, args)
	}

	var showBuildInfo bool
	var showRepo bool
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show bump tool version",
		Run: func(cmd *cobra.Command, args []string) {
			if showBuildInfo {
				buildInfo := version.Get()
				fmt.Printf("Bump Version: %s\n", buildInfo.Version)
				fmt.Printf("Git Commit: %s\n", buildInfo.GitCommit)
				fmt.Printf("Build Date: %s\n", buildInfo.BuildDate)
				fmt.Printf("Go Version: %s\n", buildInfo.GoVersion)
			} else if showRepo {
				currentVersion := bump.GetCurrentVersion()
				fmt.Printf("Repository version: %s\n", currentVersion)
			} else {
				// Show tool version by default (for CI compatibility)
				buildInfo := version.Get()
				fmt.Printf("%s\n", buildInfo.Version)
			}
		},
	}
	versionCmd.Flags().BoolVar(&showBuildInfo, "build-info", false, "Show detailed build information")
	versionCmd.Flags().BoolVar(&showRepo, "repo", false, "Show current repository version")

	quickCmd := &cobra.Command{
		Use:   "quick [patch|minor|major]",
		Short: "Quick release without prompts",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runQuickMode(args[0])
		},
	}

	interactiveCmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive release with prompts and checks",
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveMode()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show current repository version and status",
		Run: func(cmd *cobra.Command, args []string) {
			currentVersion := bump.GetCurrentVersion()
			fmt.Printf("Current repository version: %s\n", currentVersion)
		},
	}

	tagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "List all tags sorted by creation date (newest first)",
		Run: func(cmd *cobra.Command, args []string) {
			release := bump.NewRelease(cfg)
			if err := release.ListTags(); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(versionCmd, quickCmd, interactiveCmd, statusCmd, tagsCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runInteractiveMode() {
	fmt.Println("Starting interactive release mode...")
	release := bump.NewRelease(cfg)
	if err := release.RunInteractive(); err != nil {
		log.Fatal(err)
	}
}

func runQuickMode(versionType string) {
	fmt.Printf("Running quick %s release...\n", versionType)
	release := bump.NewRelease(cfg)
	if err := release.RunQuick(versionType); err != nil {
		log.Fatal(err)
	}
}
