{ pkgs ? import ./nix {} }:
with pkgs;
mkShell {
  buildInputs = [
    gnumake
    go
    go-md2man
    golangci-lint
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
