{ sources ? import ./nix/sources.nix
, pkgs ? import ./nix { inherit sources; }
}:
let
  direnv = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  #imputsFrom = [ direnv ];
  buildInputs = [
    pkgs.rustPlatform.rust.cargo
    pkgs.rustPlatform.rust.rustc
    pkgs.rustfmt
  ];
}
