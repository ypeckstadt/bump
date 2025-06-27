# Installation Guide

This guide covers different ways to install Bump on your system.

## Installation Methods

### Homebrew (Recommended for macOS/Linux)

The easiest way to install Bump:

```bash
# Add the custom tap
brew tap ypeckstadt/tap

# Install bump
brew install bump

# Verify installation
bump --version
```

### Go Install

If you have Go 1.22+ installed:

```bash
go install github.com/ypeckstadt/bump/cmd/bump@latest
```

**Note:** Make sure your `$GOPATH/bin` is in your PATH.

### Download Binary

Download the latest release from GitHub:

1. Visit the [Releases page](https://github.com/ypeckstadt/bump/releases)
2. Download the appropriate binary for your platform
3. Extract and move to your PATH:

```bash
# macOS/Linux example
tar -xzf bump-darwin-amd64.tar.gz
sudo mv bump /usr/local/bin/

# Or add to personal bin directory
mkdir -p ~/bin
mv bump ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Build from Source

For developers or custom builds:

```bash
# Clone the repository
git clone https://github.com/ypeckstadt/bump.git
cd bump

# Install development dependencies
make dev-deps

# Build the binary
make build

# Install to your PATH
make install
```

## Verify Installation

```bash
# Check if bump is installed
bump version

# Show build information
bump version --build-info
```

## System Requirements

- **Operating System**: Linux, macOS, Windows
- **Git**: Required for git operations
- **Go**: 1.22+ (if building from source or using `go install`)

**Note:** Bump is tested on Go 1.22 and 1.23 in CI.

## Optional Dependencies

### Development Tools

For the best experience, install these optional tools:

```bash
# Linting (optional but recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Security scanning (optional)
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Make sure your `GOPATH/bin` is in your PATH
2. Reload your shell: `source ~/.bashrc` or `source ~/.zshrc`
3. Verify the binary location: `which bump`

### Permission Issues

On macOS, you might need to allow the binary:

```bash
# If you get a security warning
xattr -d com.apple.quarantine /path/to/bump
```

### Build Issues

If building from source fails:

1. Ensure you have Go 1.24 or later: `go version`
2. Check your GOPATH is set correctly
3. Try cleaning and rebuilding:
   ```bash
   make clean
   make build
   ```

## Next Steps

- Read the [Usage Guide](usage.md) to get started
- Check out [Examples](examples.md) for common workflows
- Review [Configuration](configuration.md) for advanced setup