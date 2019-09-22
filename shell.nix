{ sources ? import ./nix/sources.nix
, pkgs ? import ./nix { inherit sources; }
}:
let
  direnv = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  #imputsFrom = [ direnv ];
  buildInputs = [
    pkgs.rust.packages.stable.cargo
    pkgs.rust.packages.stable.clippy
    pkgs.rust.packages.stable.rls
    pkgs.rust.packages.stable.rustc
    pkgs.rust.packages.stable.rustfmt
  ];
}
