{ pkgs ? import ./nix { } }:
with pkgs;

buildGoModule rec {
  name = "direnv-${version}";
  version = lib.fileContents ./version.txt;
  subPackages = [ "." ];

  vendorSha256 = "1kmgl3wkbyirqkf7rwlcss8ml7incxclsb36nskwx1w6b1ilphrz";

  src = lib.cleanSource ./.;

  # we have no bash at the moment for windows
  makeFlags = stdenv.lib.optional (!stdenv.hostPlatform.isWindows) [
    "BASH_PATH=${bash}/bin/bash"
  ];

  installPhase = ''
    make install PREFIX=$out
  '';

  meta = with stdenv.lib; {
    description = "A shell extension that manages your environment";
    homepage = https://direnv.net;
    license = licenses.mit;
    maintainers = with maintainers; [ zimbatm ];
  };
}
