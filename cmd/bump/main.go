package main

import (
	"fmt"
	"log"

	"bump/internal/bump"
	"bump/internal/config"
	"bump/pkg/version"

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

	var showBuildInfo bool
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		Run: func(cmd *cobra.Command, args []string) {
			if showBuildInfo {
				buildInfo := version.Get()
				fmt.Printf("Bump Version: %s\n", buildInfo.Version)
				fmt.Printf("Git Commit: %s\n", buildInfo.GitCommit)
				fmt.Printf("Build Date: %s\n", buildInfo.BuildDate)
				fmt.Printf("Go Version: %s\n", buildInfo.GoVersion)
			} else {
				currentVersion := bump.GetCurrentVersion()
				fmt.Printf("Current version: %s\n", currentVersion)
			}
		},
	}
	versionCmd.Flags().BoolVar(&showBuildInfo, "build-info", false, "Show build information")

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

	rootCmd.AddCommand(versionCmd, quickCmd, interactiveCmd)

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
