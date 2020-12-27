{ pkgs ? import ./nix {} }:
with pkgs;

buildGoModule rec {
  name = "direnv-${version}";
  version = lib.fileContents ./version.txt;
  subPackages = ["."];

  vendorSha256 = "084x7d7sshcsyim76d6pl6127nlqacgwwnm965srl9y5w5nqzba6";

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
