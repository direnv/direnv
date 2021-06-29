{ pkgs ? import ./nix { } }:
with pkgs;

buildGoModule rec {
  name = "direnv-${version}";
  version = lib.fileContents ./version.txt;
  subPackages = [ "." ];

  vendorSha256 = "0ig0vijmfqjbi6i4fv0b8vaglwnczq0kpk2ifx0q3byv8k4v2j2b";

  src = lib.cleanSource ./.;

  # we have no bash at the moment for windows
  makeFlags = lib.optional (!stdenv.hostPlatform.isWindows) [
    "BASH_PATH=${bash}/bin/bash"
  ];

  installPhase = ''
    make install PREFIX=$out
  '';

  meta = {
    description = "A shell extension that manages your environment";
    homepage = https://direnv.net;
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [ zimbatm ];
  };
}
