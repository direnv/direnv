{ sources ? import ./sources.nix }:
let
  overlay = _: pkgs: {
    sources = sources;
    naersk = pkgs.callPackage sources.naersk {};
  };

  pkgs = import sources.nixpkgs {
    config = {};
    overlays = [ overlay ];
  };
in
  pkgs

