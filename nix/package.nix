{ buildGoApplication, lib, stdenv, bash, git }:
buildGoApplication {
  pname = "direnv";
  version = "0+git";
  subPackages = [ "." ];

  nativeBuildInputs = [ git ];

  src = ../.;

  modules = ./gomod2nix.toml;

  # we have no bash at the moment for windows
  BASH_PATH =
    lib.optionalString (!stdenv.hostPlatform.isWindows)
      "${bash}/bin/bash";

  # replace the build phase to use the GNUMakefile instead
  buildPhase = ''
    make BASH_PATH=$BASH_PATH
  '';

  installPhase = ''
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
