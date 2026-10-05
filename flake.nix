{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.systems.url = "github:nix-systems/default";

  outputs =
    {
      self,
      nixpkgs,
      gomod2nix,
      systems,
    }:
    let
      eachSystem =
        f:
        nixpkgs.lib.genAttrs (import systems) (
          system:
          f rec {
            callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
            gomod2nixPkgs = gomod2nix.legacyPackages.${system};
            inherit system;
            pkgs = nixpkgs.legacyPackages.${system};
          }
        );
    in
    {

      packages = eachSystem ({ callPackage, gomod2nixPkgs, ... }: {

          default = callPackage ./. { inherit (gomod2nixPkgs) buildGoApplication; };
        }
      );

      devShells = eachSystem ({ callPackage, gomod2nixPkgs, ... }: {
          default = callPackage ./shell.nix { inherit (gomod2nixPkgs) mkGoEnv gomod2nix; };
        }
      );

      checks = eachSystem ({ pkgs, system, ... }:
        let
          sourceFiles = pkgs.lib.fileset.toSource {
            root = ./.;
            fileset = pkgs.lib.fileset.unions [
              ./go.mod
              ./go.sum
              ./GNUmakefile
              ./stdlib.sh
              ./version.txt
              ./README.md
              (pkgs.lib.fileset.fileFilter (file: file.hasExt "go") ./.)
              ./test
              ./internal
              ./pkg
              (pkgs.lib.fileset.fileFilter (file: file.name == ".envrc") ./.)
            ];
          };
        in {
          package = self.packages.${system}.default;
          tests = self.packages.${system}.default.overrideAttrs (old: {
            src = sourceFiles;
            nativeBuildInputs = (old.nativeBuildInputs or []) ++ self.devShells.${system}.default.nativeBuildInputs;
            buildPhase = ''
              export GOLANGCI_LINT_CACHE=$TMPDIR/golangci-cache
              export XDG_CACHE_HOME=$TMPDIR/cache
              export HOME=$TMPDIR/home
              mkdir -p $GOLANGCI_LINT_CACHE $XDG_CACHE_HOME $HOME

              # Patch shebangs in test files
              patchShebangs test/

              make test
            '';
            installPhase = ''
              mkdir -p $out
              touch $out/tests-passed
            '';
          });
          dist = self.packages.${system}.default.overrideAttrs (old: {
            src = sourceFiles;
            nativeBuildInputs = (old.nativeBuildInputs or []) ++ self.devShells.${system}.default.nativeBuildInputs;
            buildPhase = ''
              make dist
            '';
            installPhase = ''
              mkdir -p $out
              cp -r dist/* $out/
            '';
          });
        }
      );
    };
}
