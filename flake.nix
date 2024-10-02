{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.systems.url = "github:nix-systems/default";

  outputs =
    {
      self,
      nixpkgs,
      gomod2nix,
      systems,
    }:
    let
      eachSystem =
        f:
        nixpkgs.lib.genAttrs (import systems) (
          system:
          f rec {
            callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
            gomod2nixPkgs = gomod2nix.legacyPackages.${system};
            inherit system;
            pkgs = nixpkgs.legacyPackages.${system};
          }
        );
    in
    {
      packages = eachSystem (
        { callPackage, gomod2nixPkgs, ... }:
        {
          default = callPackage ./. { inherit (gomod2nixPkgs) buildGoApplication; };
        }
      );

      devShells = eachSystem (
        { callPackage, gomod2nixPkgs, ... }:
        {
          default = callPackage ./shell.nix { inherit (gomod2nixPkgs) mkGoEnv gomod2nix; };
        }
      );

      checks = eachSystem (
        { system, ... }:
        {
          check-package = self.packages.${system}.default;
        }
      );
    };
}
