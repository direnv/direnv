{ pkgs ? import <nixpkgs> {} }:
let
  direnv = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  #imputsFrom = [ direnv ];
  buildInputs = with pkgs; [
    cargo
    carnix
    rustc
  ];
}
