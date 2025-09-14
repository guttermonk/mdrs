# NixOS Setup Guide for mdrs

## Overview
This guide explains how to build, run, and develop `mdrs` (Markdown Renderer & Search) on NixOS using Nix flakes.

## Prerequisites
- NixOS or Nix package manager installed
- Flakes enabled (for flake.nix usage)

### Enabling Flakes
If you haven't enabled flakes yet, add this to your NixOS configuration (`/etc/nixos/configuration.nix`):
```nix
nix.settings.experimental-features = [ "nix-command" "flakes" ];
```

Or for non-NixOS systems, add to `~/.config/nix/nix.conf`:
```
experimental-features = nix-command flakes
```

## Quick Start

### Running mdrs directly from GitHub
```bash
# Run the latest version directly
nix run github:guttermonk/mdrs -- README.md

# Or from your local fork
nix run github:yourusername/mdrs -- README.md
```

### Running from local repository
```bash
# Clone the repository
git clone https://github.com/guttermonk/mdrs.git
cd mdrs

# Run directly
nix run . -- README.md

# Or build and run
nix build
./result/bin/mdrs README.md
```

## Building mdrs

### Using Nix Flakes (Recommended)
```bash
# Build the package
nix build

# The binary will be available at
./result/bin/mdrs

# Install to your profile
nix profile install .
```

### Using Legacy Nix
```bash
# Enter development shell
nix-shell

# Build with go
go build

# Or use make
make build
```

## Development Environment

### Using Nix Flakes
```bash
# Enter the development shell
nix develop

# Or with direnv (recommended)
echo "use flake" > .envrc
direnv allow
```

### Development Shell Features
The development environment includes:
- Go compiler and tools
- gopls (Go language server)
- go-tools (gofmt, goimports, etc.)
- golangci-lint (linting)
- make (build automation)
- git (version control)
- gox (cross-compilation)
- entr (file watching)
- ripgrep (fast searching)

### Building in Development Shell
```bash
# Enter development shell
nix develop

# Build the project
go build

# Run tests
go test ./...

# Run with a file
go run . README.md

# Build with make
make build

# Install locally
make install
```

## Vendor Hash Updates

When dependencies change (go.mod/go.sum), you'll need to update the vendor hash in `flake.nix`:

1. First, try to build with an incorrect hash:
```bash
nix build
```

2. Nix will fail and show you the correct hash:
```
error: hash mismatch in fixed-output derivation '/nix/store/...':
  specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
  got:       sha256-BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=
```

3. Copy the "got" hash and update `vendorHash` in `flake.nix`:
```nix
vendorHash = "sha256-BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=";
```

## Integration with NixOS Configuration

### System-wide Installation
Add to your NixOS configuration:
```nix
{ pkgs, ... }:
{
  environment.systemPackages = [
    (pkgs.callPackage /path/to/mdrs {})
  ];
}
```

### Home Manager Installation
If using Home Manager:
```nix
{ pkgs, ... }:
{
  home.packages = [
    (pkgs.callPackage /path/to/mdrs {})
  ];
}
```

### As a Flake Input
In your system flake:
```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    mdrs.url = "github:guttermonk/mdrs";
  };

  outputs = { self, nixpkgs, mdrs, ... }: {
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        {
          environment.systemPackages = [
            mdrs.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

## Cross-Compilation

The flake supports cross-compilation for multiple architectures:
```bash
# Build for different systems
nix build .#packages.x86_64-linux.mdrs
nix build .#packages.aarch64-linux.mdrs
nix build .#packages.x86_64-darwin.mdrs
nix build .#packages.aarch64-darwin.mdrs
```

## Direnv Integration

For automatic environment loading:

1. Install direnv:
```bash
nix-env -iA nixpkgs.direnv
```

2. Add to your shell configuration:
```bash
# For bash
eval "$(direnv hook bash)"

# For zsh
eval "$(direnv hook zsh)"

# For fish
direnv hook fish | source
```

3. Allow direnv in the project:
```bash
cd mdrs
direnv allow
```

Now the development environment will automatically load when you enter the directory.

## Troubleshooting

### Flakes Not Enabled
If you get an error about experimental features:
```bash
# Run with experimental features enabled
nix --experimental-features 'nix-command flakes' build
```

### Go Module Download Issues
If you encounter issues with Go module downloads:
```bash
# Clear Go module cache
go clean -modcache

# Update dependencies
go mod download
go mod tidy
```

### Permission Denied
If you get permission errors:
```bash
# Ensure proper ownership
chown -R $USER:$USER .

# Clean build artifacts
rm -rf result result-*
nix build
```

## Advanced Usage

### Custom Build Flags
Modify the `ldflags` in `flake.nix` to add custom build information:
```nix
ldflags = [
  "-s"
  "-w"
  "-X main.GitCommit=${gitCommit}"
  "-X main.Version=${version}"
  "-X main.BuildDate=${builtins.toString builtins.currentTime}"
];
```

### Running Tests in Nix Build
Enable tests by setting `doCheck = true` in `flake.nix`:
```nix
doCheck = true;
checkPhase = ''
  go test -v ./...
'';
```

## Contributing

When contributing to mdrs on NixOS:

1. Use the development shell:
```bash
nix develop
```

2. Make your changes and test:
```bash
go test ./...
go build
./mdrs README.md
```

3. Update vendor hash if dependencies changed
4. Test the Nix build:
```bash
nix build
```

5. Submit your pull request

## Resources

- [Nix Flakes Documentation](https://nixos.wiki/wiki/Flakes)
- [Go Modules in Nix](https://nixos.wiki/wiki/Go)
- [mdrs Repository](https://github.com/guttermonk/mdrs)