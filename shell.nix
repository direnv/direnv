{ pkgs ? import <nixpkgs> {} }:
let
  direnv = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  #imputsFrom = [ direnv ];
  buildInputs = with pkgs; [
    carnix
    rust.cargo
    rust.rustc
  ];
}
