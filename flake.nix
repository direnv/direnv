{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

          gomod2nixBuilder = callPackage "${gomod2nix}/builder" {
            gomod2nix = gomod2nix';
          };
          gomod2nix' = callPackage "${gomod2nix}/default.nix" {
            inherit (gomod2nixBuilder) mkGoEnv buildGoApplication;
          };
        in
        {
          packages.default = callPackage ./. {
            inherit (gomod2nixBuilder) buildGoApplication;
          };
          devShells.default = callPackage ./shell.nix {
            inherit (gomod2nixBuilder) mkGoEnv;
            gomod2nix = gomod2nix';
          };
        }));
}
