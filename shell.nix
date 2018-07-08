{ pkgs ? import <nixpkgs> {} }:
with pkgs;
let
  elvish012 = elvish.overrideDerivation (self: {
    name = "elvis-0.12-rc2";

    src = fetchFromGitHub {
      owner = "elves";
      repo = "elvish";
      rev = "0.12-rc2";
      sha256 = "1nxijphdbqbv032s9mw7miypy92xdyzkql1bcqh0s1d0ghi66hvm";
    };

    excludedPackages = [ "github.com/elves/elvish/website" ];
  });
in
mkShell {
  buildInputs = [
    gnumake
    go
    go-md2man
    gox
    shellcheck
    shfmt
    which

    # Shells
    bashInteractive
    elvish012
    fish
    tcsh
    zsh
  ];
}
