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
            inherit (pkgs) callPackage;
            gomod2nixPkgs = gomod2nix.legacyPackages.${system};
            inherit system;
            pkgs = nixpkgs.legacyPackages.${system};
          }
        );
    in
    {

      packages = eachSystem (
        { callPackage, gomod2nixPkgs, ... }:
        {

          default = callPackage ./. { inherit (gomod2nixPkgs) buildGoApplication; };
        }
      );

      devShells = eachSystem (
        { callPackage, gomod2nixPkgs, ... }:
        {
          default = callPackage ./shell.nix { inherit (gomod2nixPkgs) gomod2nix; };
        }
      );

      checks = eachSystem (
        { pkgs, system, ... }:
        let
          inherit (pkgs) lib;
          direnv = self.packages.${system}.default;
          direnvNoMan = direnv.override { __includeMan = false; };

          testSrc =
            files:
            lib.fileset.toSource {
              root = ./.;
              fileset = lib.fileset.unions files;
            };

          runTest =
            name:
            {
              files ? [ ./test ],
              inputs ? [ ],
              command,
            }:
            pkgs.runCommand "direnv-${name}"
              {
                src = testSrc files;
                nativeBuildInputs = [ pkgs.bashInteractive ] ++ inputs;
                # Fixes pwsh tests on (sandboxed) macOS
                sandboxProfile = ''
                  (allow file-read* (subpath "/usr/share/icu"))
                '';
              }
              ''
                # Keep the old store mtimes: some tests expect `touch .envrc` to change the mtime without a sleep.
                cp -pR $src/. .
                chmod -R u+w .
                patchShebangs test/
                ln -s ${direnv}/bin/direnv direnv
                export HOME=$TMPDIR/home
                mkdir -p $HOME
                ${command}
                touch $out
              '';

          shellTest =
            name: shell: command:
            runTest name {
              inputs = [
                pkgs.python3
                pkgs.ruby
                shell
              ];
              inherit command;
            };
        in
        {
          package = direnv;

          test-shellcheck = runTest "test-shellcheck" {
            files = [
              ./stdlib.sh
              ./test/stdlib.bash
            ];
            inputs = [ pkgs.shellcheck ];
            command = ''
              shellcheck stdlib.sh
              shellcheck ./test/stdlib.bash
            '';
          };

          test-stdlib = runTest "test-stdlib" {
            files = [
              ./README.md
              ./stdlib.sh
              ./test
            ];
            command = "./test/stdlib.bash";
          };

          test-go-lint = direnvNoMan.overrideAttrs (old: {
            name = "direnv-test-go-lint";
            nativeBuildInputs = (old.nativeBuildInputs or [ ]) ++ [ pkgs.golangci-lint ];
            buildPhase = ''
              export HOME=$TMPDIR/home
              export GOLANGCI_LINT_CACHE=$TMPDIR/golangci-cache
              cp ${./.golangci.yml} .golangci.yml
              golangci-lint run
            '';
            doCheck = false;
            installPhase = "touch $out";
          });

          test-bash = shellTest "test-bash" pkgs.bashInteractive "bash ./test/direnv-test.bash";
          test-elvish = shellTest "test-elvish" pkgs.elvish "elvish ./test/direnv-test.elv";
          test-fish = shellTest "test-fish" pkgs.fish "fish ./test/direnv-test.fish";
          test-tcsh = shellTest "test-tcsh" pkgs.tcsh "tcsh -e ./test/direnv-test.tcsh";
          test-zsh = shellTest "test-zsh" pkgs.zsh "zsh ./test/direnv-test.zsh";
          test-pwsh = shellTest "test-pwsh" pkgs.powershell "pwsh ./test/direnv-test.ps1";
          test-mx = shellTest "test-mx" pkgs.murex "murex -trypipe ./test/direnv-test.mx";
        }
        // lib.optionalAttrs (system == "x86_64-linux") {
          dist = direnvNoMan.overrideAttrs {
            name = "direnv-dist";
            buildPhase = "make -j$NIX_BUILD_CORES dist";
            doCheck = false;
            installPhase = ''
              mkdir -p $out
              cp -r dist/* $out/
            '';
          };
        }
      );
    };
}
