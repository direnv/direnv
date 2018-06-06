{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
  buildInputs = [
    bashInteractive
    elvish
    fish
    go
    tcsh
    zsh
  ];

  shellHook = ''
    unset GOPATH
  '';
}
