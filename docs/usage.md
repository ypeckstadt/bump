# Usage Guide

This guide covers how to use Bump for semantic version management in your projects.

## Quick Start

```bash
# Show current version
bump version

# Interactive release (recommended)
bump

# Quick release
bump quick patch
```

## Commands Overview

### Interactive Mode

The interactive mode provides a guided experience with validation and checks:

```bash
bump
# or
bump interactive
```

**Features:**
- Shows current version and recent commits
- Interactive version type selection
- Pre-release checks (build, test, lint)
- Working directory validation
- Confirmation prompts

### Quick Mode

For automated releases and CI/CD:

```bash
bump quick patch    # v1.2.3 → v1.2.4
bump quick minor    # v1.2.3 → v1.3.0
bump quick major    # v1.2.3 → v2.0.0
```

### Version Information

```bash
# Show bump tool version (for CI/scripts)
bump version                    # Output: 1.0.0

# Show detailed build information
bump version --build-info

# Show current repository version from git tags
bump version --repo             # Output: Repository version: v1.2.3
bump status                     # Output: Current repository version: v1.2.3
```

## Global Flags

### Dry Run Mode

Preview what would happen without making changes:

```bash
bump --dry-run
bump quick minor --dry-run
```

### Verbose Output

Enable detailed output for debugging:

```bash
bump --verbose
bump quick patch --verbose
```

## Version Types

### Patch Release (Bug Fixes)

Use for:
- Bug fixes
- Security patches  
- Documentation updates
- Small improvements

```bash
bump quick patch
```

### Minor Release (New Features)

Use for:
- New features
- Backward compatible changes
- Performance improvements
- New functionality

```bash
bump quick minor
```

### Major Release (Breaking Changes)

Use for:
- Breaking API changes
- Major rewrites
- Incompatible changes
- New major versions

```bash
bump quick major
```

## Workflow Examples

### Development Workflow

```bash
# 1. Finish your feature work
git add .
git commit -m "Add new encryption feature"

# 2. Run interactive release
bump

# 3. Select version type (minor for new feature)
# 4. Confirm and push

# 5. GitHub Actions handles the rest
```

### CI/CD Integration

```bash
# In your CI pipeline
bump quick patch
```

### Hotfix Workflow

```bash
# 1. Create hotfix branch
git checkout -b hotfix/security-fix

# 2. Make your fix
git commit -m "Fix security vulnerability"

# 3. Quick patch release
bump quick patch

# 4. Merge back to main
```

## Pre-release Checks

Bump automatically runs these checks before creating releases:

1. **Build Check**: `go build ./...`
2. **Test Check**: `go test ./...`
3. **Lint Check**: `golangci-lint run` (if available)
4. **Go Mod Tidy**: `go mod tidy`

### Skipping Checks

Use dry-run mode to skip actual execution:

```bash
bump --dry-run
```

## Git Integration

### Requirements

- Must be in a git repository
- Must have git write access for tag pushing
- Clean working directory (interactive mode warns if dirty)

### What Bump Does

1. **Validates** git repository and working directory
2. **Gets** current version from latest git tag
3. **Calculates** new version based on type
4. **Creates** annotated git tag with release message
5. **Pushes** tag to origin remote

### Tag Format

Bump uses semantic versioning with `v` prefix:

- `v1.0.0` - Major release
- `v1.1.0` - Minor release  
- `v1.1.1` - Patch release

## Error Handling

### Common Errors

**Not a git repository:**
```bash
Error: not a git repository
```
→ Run `git init` or navigate to a git repository

**Tag already exists:**
```bash
Error: tag v1.2.3 already exists
```
→ Delete the tag or use a different version

**Working directory not clean:**
```bash
Warning: Working directory is not clean
```
→ Commit or stash changes, or continue anyway

**No existing tags:**
```bash
Current version: v0.0.0
```
→ Bump will start from v0.0.0 for first release

## Integration with GitHub Actions

Bump creates git tags that trigger GitHub Actions workflows:

1. **Tag created** → GitHub Actions triggered
2. **Build binaries** for multiple platforms
3. **Create GitHub release** with assets
4. **Update package registries** (Homebrew, etc.)

### Example Workflow

```yaml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
      - run: make build-all
      - uses: softprops/action-gh-release@v1
```

## Best Practices

### Version Selection

- **Patch**: Only for fixes, no new functionality
- **Minor**: New features, but backward compatible
- **Major**: Breaking changes or major milestones

### Commit Messages

Use clear, descriptive commit messages:

```bash
# Good
git commit -m "Add AES-256 encryption support"
git commit -m "Fix memory leak in file parser"

# Less helpful
git commit -m "Update code"
git commit -m "Fix bug"
```

### Release Messages

Interactive mode prompts for release messages:

```
Release message: (Release v1.2.0) Add encryption and improve performance
```

### Testing

Always test before releases:

```bash
# Run tests manually
go test ./...

# Or let bump run pre-release checks
bump --verbose
```

## Next Steps

- Learn about [Configuration](configuration.md) options
- See [Examples](examples.md) for real-world scenarios
- Review [Architecture](architecture.md) for technical details