{ sources ? import ./nix/sources.nix
, pkgs ? import ./nix { inherit sources; }
}:
let
  lib = pkgs.lib;

  # another attempt to make filterSource nicer to use
  allowSource = { allow, src }:
    let
      out = builtins.filterSource filter src;
      filter = path: _fileType:
        lib.any (checkElem path) allow;
      checkElem = path: elem:
        lib.hasPrefix (toString elem) (toString path);
    in
      out;

  src = allowSource {
    allow = [
      ./Cargo.lock
      ./Cargo.toml
      ./src
    ];
    src = ./.;
  };
in
pkgs.naersk.buildPackage src {
  cratePaths = [ "." ];
}
