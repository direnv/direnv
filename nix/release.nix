{ buildGoApplication, lib, stdenv, bash, git, goreleaser }:
buildGoApplication {
  pname = "direnv-release";
  version = "0+git";
  src = ../.;

  nativeBuildInputs = [ git goreleaser ];

  modules = ./gomod2nix.toml;

  # replace the build phase to use the GNUMakefile instead
  buildPhase = ''
    goreleaser release --clean --snapshot
  '';

  installPhase = ''
    mv dist/ $out
  '';
}
