{
  description = "mdrs - Markdown Renderer & Search for the terminal";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        
        # Version information
        version = "0.1.0";
        gitCommit = if (self ? rev) then self.rev else "dirty";
        
      in
      {
        packages = {
          default = self.packages.${system}.mdrs;
          
          mdrs = pkgs.buildGoModule {
            pname = "mdrs";
            inherit version;
            
            src = ./.;
            
            # Generate vendor hash with: nix run nixpkgs#nix-prefetch-git -- --url . --fetch-submodules
            # Or let nix tell you the correct hash on first build
            vendorHash = "sha256-rCd3Zccr2B1kGXnRqUDuBYpe78PP9LiDuAMMr0jvkvI=";
            
            # Add version information as build flags
            ldflags = [
              "-s"
              "-w"
              "-X main.GitCommit=${gitCommit}"
              "-X main.GitLastTag=v${version}"
              "-X main.GitExactTag=v${version}"
            ];
            
            # Rename binary from mdr to mdrs if needed
            postInstall = ''
              if [ -f $out/bin/mdr ]; then
                mv $out/bin/mdr $out/bin/mdrs
              fi
            '';
            
            # Disable tests that might require network access
            doCheck = false;
            
            meta = with pkgs.lib; {
              description = "A standalone Markdown renderer for the terminal with search functionality";
              homepage = "https://github.com/guttermonk/mdrs";
              license = licenses.mit;
              maintainers = [];
              mainProgram = "mdrs";
            };
          };
        };
        
        apps = {
          default = self.apps.${system}.mdrs;
          
          mdrs = {
            type = "app";
            program = "${self.packages.${system}.mdrs}/bin/mdrs";
          };
        };
        
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              # Go development
              go
              gopls
              go-tools
              golangci-lint
              
              # Build tools
              gnumake
              git
              
              # Optional: for cross-compilation
              gox
              
              # Helpful tools
              entr  # for file watching
              ripgrep  # for searching
            ];
            
            shellHook = ''
              echo "mdrs development environment"
              echo "Available commands:"
              echo "  go build       - Build the project"
              echo "  go run . FILE  - Run mdrs with a markdown file"
              echo "  make build     - Build using Makefile"
              echo "  nix build      - Build with nix"
              echo "  nix run        - Run the built version"
              echo ""
              echo "Go version: $(go version)"
            '';
          };
        };
        
        # Legacy support for nix-shell
        devShell = self.devShells.${system}.default;
      });
}