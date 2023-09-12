{ pkgs ? import ./nix { }, vendorHash ? "sha256-UsdVGKIoiI3nJgbIdSg+BIDInoUODFjfyvoTdxg2a8Q=" }:
let
  inherit (pkgs)
    bash
    buildGoModule
    lib
    stdenv
    ;
in
buildGoModule rec {
  name = "direnv-${version}";
  version = lib.fileContents ./version.txt;
  subPackages = [ "." ];

  inherit vendorHash;

  src = builtins.fetchGit ./.;

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
