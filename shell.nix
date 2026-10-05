{ stdenv
, pkgs
, mkGoEnv
, gomod2nix
, git
, git-extras
, gnumake
, go
, go-md2man
, gox
, bashInteractive
, elvish
, fish
, tcsh
, zsh
, powershell
, murex
, golangci-lint
, python3
, ruby
, shellcheck
, shfmt
, cacert
}:
stdenv.mkDerivation {
  name = "shell";
  nativeBuildInputs = with pkgs; [

    (mkGoEnv { pwd = ./.; go = go; })

    go

    # Build
    git
    git-extras # for git-changelog
    gnumake
    go-md2man
    gomod2nix

    # Shells
    bashInteractive
    elvish
    fish
    tcsh
    zsh
    powershell
    murex

    golangci-lint
    python3
    ruby
    shellcheck
    shfmt
  ];

  shellHook = ''
    unset GOPATH GOROOT
    # needed in pure shell
    export HOME=''${HOME:-$TMPDIR}

    export GO111MODULE=on
    export SSL_CERT_FILE=${cacert}/etc/ssl/certs/ca-bundle.crt
  '';
}
