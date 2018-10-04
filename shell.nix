{ pkgs ? import <nixpkgs> {} }:
with pkgs;
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
    elvish
    fish
    tcsh
    zsh
  ];
}
