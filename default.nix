{ buildGoApplication, lib, stdenv, bash }:
buildGoApplication {
  pname = "direnv";
  version = lib.fileContents ./version.txt;
  subPackages = [ "." ];

  src = ./.;
  pwd = ./.;
  modules = ./gomod2nix.toml;

  # we have no bash at the moment for windows
  BASH_PATH =
    lib.optionalString (!stdenv.hostPlatform.isWindows)
      "${bash}/bin/bash";

  # replace the build phase to use the GNUMakefile instead
  buildPhase = ''
    ls -la ./vendor
    make BASH_PATH=$BASH_PATH
  '';

  installPhase = ''
    echo $GOCACHE
    make install PREFIX=$out
  '';

  meta = {
    description = "A shell extension that manages your environment";
    homepage = "https://direnv.net";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.zimbatm ];
    mainProgram = "direnv";
  };
}
