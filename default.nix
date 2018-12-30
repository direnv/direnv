{ pkgs ? import <nixpkgs> {} }:
let
  inherit (pkgs) lib rustPlatform;
in
rustPlatform.buildRustPackage {
  name = "direnv";
  src = lib.cleanSource ./.;
  cargoSha256 = "0sjjj9z1dhilhpc8pq4154czrb79z9cm044jvn75kxcjv6v5l2m5";
}
