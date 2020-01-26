{ pkgs ? import ./nix {} }:
with pkgs;
mkShell {
  buildInputs = [
    # Build
    gitAndTools.git-extras # for git-changelog
    gnumake
    go
    go-md2man
    gox

    # Shells
    bashInteractive
    elvish
    fish
    tcsh
    zsh

    # Test dependencies
    golangci-lint
    ruby
    shellcheck
    shfmt
  ];
}
