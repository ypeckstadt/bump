# Bump

A Go CLI tool for semantic version management and releases, inspired by shell release scripts.

## Features

- 🚀 **Interactive Release Mode** - Full-featured with prompts, validation, and pre-release checks
- ⚡ **Quick Release Mode** - Fast command-line releases perfect for CI/CD
- 📋 **Pre-release Checks** - Automatic build, test, and lint validation
- 🏷️ **Git Integration** - Tag creation, validation, and pushing
- 🎨 **Colored Output** - Clear visual feedback with progress indicators
- 🔍 **Dry Run Mode** - Preview changes without making them

## Installation

```bash
go install github.com/yourusername/bump@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/bump
cd bump
go build -o bump
```

## Usage

### Interactive Mode (Recommended)

```bash
# Start interactive release process
bump

# Or explicitly use interactive mode
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

```bash
# Quick patch release (v1.2.3 → v1.2.4)
bump quick patch

# Quick minor release (v1.2.3 → v1.3.0)  
bump quick minor

# Quick major release (v1.2.3 → v2.0.0)
bump quick major
```

### Other Commands

```bash
# Show current version
bump version

# Dry run mode (preview without changes)
bump --dry-run

# Verbose output
bump --verbose

# Help
bump --help
```

## Version Types

| Type | When to Use | Example |
|------|-------------|---------|
| **Patch** | Bug fixes, security patches | v1.2.3 → v1.2.4 |
| **Minor** | New features, backward compatible | v1.2.3 → v1.3.0 |
| **Major** | Breaking changes, major rewrites | v1.2.3 → v2.0.0 |

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

## Pre-release Checks

The tool automatically runs these checks before creating releases:

- **Build** - `go build ./...`
- **Tests** - `go test ./...` 
- **Lint** - `golangci-lint run` (if available)
- **Go mod tidy** - `go mod tidy`

## Integration with CI/CD

### GitHub Actions

```yaml
- name: Quick Release
  run: |
    ./bump quick patch
```

### Manual Workflow

```bash
# 1. Finish feature work
git add . && git commit -m "Add new feature"

# 2. Run interactive release  
bump

# 3. Choose version type
# 4. Let CI/CD handle the rest
```

## Safety Features

- Working directory cleanliness validation
- Tag existence checking
- Multiple confirmation prompts (interactive mode)
- Dry-run mode for testing
- Graceful error handling

## Requirements

- Git repository with write access
- Go 1.24 or later
- Optional: golangci-lint for linting checks

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `bump` to create a release (dogfooding!)

## License

MIT License - see LICENSE file for details