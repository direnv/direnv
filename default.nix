{
  buildGoApplication,
  lib,
  stdenv,
  bash,
  __includeMan ? true,
}:
buildGoApplication {
  pname = "direnv";
  version = lib.fileContents ./version.txt;
  subPackages = [ "." ];

  src = lib.fileset.toSource {
    root = ./.;
    fileset = lib.fileset.unions (
      [
        ./go.mod
        ./go.sum
        ./gomod2nix.toml
        ./GNUmakefile
        ./stdlib.sh
        ./version.txt
        ./README.md
        (lib.fileset.fileFilter (file: file.hasExt "go") ./.)
        ./test
        ./internal
        ./pkg
        (lib.fileset.fileFilter (file: file.name == ".envrc") ./.)
      ]
      ++ lib.optional __includeMan ./man
    );
  };

  modules = ./gomod2nix.toml;

  # we have no bash at the moment for windows
  BASH_PATH = lib.optionalString (!stdenv.hostPlatform.isWindows) "${bash}/bin/bash";

  # replace the build phase to use the GNUMakefile instead
  buildPhase = ''
    ls -la ./vendor
    make BASH_PATH=$BASH_PATH
  '';

  installPhase = ''
    echo $GOCACHE
    make install PREFIX=$out
  '';

  checkPhase = ''
    make test-go
  '';

  meta = {
    description = "A shell extension that manages your environment";
    homepage = "https://direnv.net";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.zimbatm ];
    mainProgram = "direnv";
  };
}
