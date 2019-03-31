{ sources ? import ./sources.nix }:
import sources.nixpkgs {
  config = {};
  overlays = [];
}
