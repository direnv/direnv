{ pkgs ? import ./nix {} }:
with pkgs;
mkShell {
  buildInputs = [
    # Build
    gitAndTools.git
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

  shellHook = ''
    unset GOPATH GOROOT
    export GO111MODULE=on
    export SSL_CERT_FILE=${cacert}/etc/ssl/certs/ca-bundle.crt
  '';
}
