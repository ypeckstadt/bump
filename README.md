# Bump

A Go CLI tool for semantic version management and releases, inspired by shell release scripts.

## Quick Start

### Install via Homebrew (recommended)
```bash
brew install ypeckstadt/homebrew-tap/bump
```

### Or install with Go
```bash
go install github.com/ypeckstadt/bump/cmd/bump@latest
```

### Start using bump
```bash
bump --version
```

```bash
bump
```

## Features

- 🚀 **Interactive Release Mode** - Full-featured with prompts, validation, and pre-release checks
- ⚡ **Quick Release Mode** - Fast command-line releases perfect for CI/CD
- 📋 **Pre-release Checks** - Automatic build, test, and lint validation
- 🏷️ **Git Integration** - Tag creation, validation, and pushing
- 🌿 **Branch Management** - Create release branches with tag creation
- 📊 **Tag Listing** - View all tags sorted by creation date
- 🎨 **Colored Output** - Clear visual feedback with progress indicators
- 🔍 **Dry Run Mode** - Preview changes without making them

## Installation

### Option 1: Homebrew (Recommended for macOS/Linux)

Add the custom tap:
```bash
brew tap ypeckstadt/tap
```

Install bump:
```bash
brew install bump
```

### Option 2: Go Install

Install directly with Go (requires Go 1.22+):
```bash
go install github.com/ypeckstadt/bump/cmd/bump@latest
```

### Option 3: Download Binary

1. Visit the [Releases page](https://github.com/ypeckstadt/bump/releases)
2. Download the appropriate binary for your platform
3. Extract and move to your PATH:

For macOS/Linux:
```bash
sudo mv bump /usr/local/bin/
```

Or add to your personal bin directory:
```bash
mv bump ~/bin/
```

### Option 4: Build from Source

Clone and build:
```bash
git clone https://github.com/ypeckstadt/bump
cd bump
make build
```

Install to your PATH:
```bash
make install
```

### Verify Installation

Check installation:
```bash
bump --version
```

Show detailed build info:
```bash
bump version --build-info
```

### Updating Bump

Update via Homebrew:
```bash
brew upgrade bump
```

Update via Go install:
```bash
go install github.com/ypeckstadt/bump/cmd/bump@latest
```

## Usage

### Interactive Mode (Recommended)

Start interactive release process:
```bash
bump
```

Or explicitly use interactive mode:
```bash
bump interactive
```

Features:
- Shows current version and recent commits
- Interactive version type selection (patch/minor/major)
- Pre-release checks (build, lint, tests)
- Custom release messages
- Git working directory validation
- Confirmation prompts

### Quick Mode

Quick patch release (v1.2.3 → v1.2.4):
```bash
bump quick patch
```

Quick minor release (v1.2.3 → v1.3.0):
```bash
bump quick minor
```

Quick major release (v1.2.3 → v2.0.0):
```bash
bump quick major
```

### Other Commands

Show bump tool version (CI-friendly):
```bash
bump version
```

```bash
bump --version
```

Show repository version from git tags:
```bash
bump status
```

```bash
bump version --repo
```

Show detailed build info:
```bash
bump version --build-info
```

List all tags sorted by creation date (newest first):
```bash
bump tags
```

Dry run mode (preview without changes):
```bash
bump --dry-run
```

Verbose output:
```bash
bump --verbose
```

Help:
```bash
bump --help
```

### Branch Management (New!)

#### Interactive Branch Creation
After creating a tag, bump will ask if you want to create a branch:
- Choose source branch (defaults to main/master)
- Choose target branch name (defaults to tag name without 'v' prefix)
- Merge automatically if branch exists
- Push branch to origin
- Automatically returns to your original branch after operations

#### Non-Interactive Branch Creation
For CI/CD pipelines, use CLI flags for automatic branch creation:

Create branch with default settings:
```bash
bump quick patch --create-branch
```

Specify all branch options:
```bash
bump quick minor \
  --create-branch \
  --source-branch main \
  --branch-name release/v1.3.0 \
  --auto-merge \
  --auto-push
```

Full example: Create patch release with automatic branch:
```bash
bump quick patch --create-branch --auto-push
```

Skip branch creation entirely:
```bash
bump --nobranch
```

```bash
bump quick patch --nobranch
```

#### Branch Options
- `--nobranch` - Skip branch creation prompt entirely
- `--create-branch` - Create a branch for the tag
- `--source-branch <name>` - Source branch (default: main/master)
- `--branch-name <name>` - Target branch name (default: tag name without 'v' prefix)
- `--auto-merge` - Automatically merge if branch exists
- `--auto-push` - Automatically push the branch

## Version Types

| Type | When to Use | Example |
|------|-------------|---------|
| **Patch** | Bug fixes, security patches | v1.2.3 → v1.2.4 |
| **Minor** | New features, backward compatible | v1.2.3 → v1.3.0 |
| **Major** | Breaking changes, major rewrites | v1.2.3 → v2.0.0 |

## Troubleshooting

### Homebrew Installation Issues

If you encounter issues with the Homebrew installation:

```bash
# Update Homebrew first
brew update

# If tap already exists, update it
brew tap ypeckstadt/tap --force

# Reinstall if needed
brew uninstall bump
brew install bump
```

### Go Install Issues

If `go install` fails:

```bash
# Check Go version (needs 1.22+)
go version

# Clear module cache if needed
go clean -modcache
go install github.com/ypeckstadt/bump/cmd/bump@latest
```

### Binary Not Found

If `bump` command is not found after installation:

```bash
# For Go install, add GOPATH/bin to PATH
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# For Homebrew, ensure Homebrew bin is in PATH
echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.bashrc  # Apple Silicon
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc     # Intel Mac
```

## Examples

### Interactive Release

```bash
$ bump
🚀 Interactive Release Mode
Current version: v1.2.3

Recent commits:
  abc123d Fix memory leak in parser
  def456e Add new encryption feature
  ghi789a Update dependencies

? Select version type:
  ▸ patch (v1.2.4) - bug fixes
    minor (v1.3.0) - new features  
    major (v2.0.0) - breaking changes

New version will be: v1.2.4
Running pre-release checks...
✅ All checks passed

? Release message: (Release v1.2.4) Fix memory leak and security updates
? Create and push tag v1.2.4? (y/N) y

Creating tag v1.2.4...
Pushing tag v1.2.4...
✅ Successfully created and pushed tag v1.2.4

? Do you want to create a branch for this tag? (y/N) y
? Source branch: (main) main
? Target branch name: (1.2.4) release/1.2.4
✅ Successfully created branch release/1.2.4 from main
? Do you want to push branch release/1.2.4 to origin? (y/N) y
✅ Successfully pushed branch release/1.2.4
Returned to branch main

GitHub Actions should now trigger the release workflow
```

### Quick Release

```bash
$ bump quick minor
Running quick minor release...
Creating minor release: v1.2.3 → v1.3.0
Creating tag v1.3.0...
Pushing tag v1.3.0...
✅ Successfully created and pushed tag v1.3.0
```

### List Tags

```bash
$ bump tags
Found 5 tags (sorted by creation date, newest first):

v1.3.0 2024-01-15 14:30:00 +0000
v1.2.4 2024-01-14 10:15:30 +0000
v1.2.3 2024-01-10 09:45:00 +0000
v1.2.2 2024-01-08 16:20:00 +0000
v1.2.1 2024-01-05 11:00:00 +0000
```

## License

MIT License - see LICENSE file for details
