{ pkgs ? import ./nix {} }:
with pkgs;

buildGoPackage rec {
  name = "direnv-${version}";
  version = lib.fileContents ./version.txt;
  goPackagePath = "github.com/direnv/direnv";
  subPackages = ["."];

  src = lib.cleanSource ./.;

  postConfigure = ''
    cd $NIX_BUILD_TOP/go/src/$goPackagePath
  '';

  # we have no bash at the moment for windows
  makeFlags = stdenv.lib.optional (!stdenv.hostPlatform.isWindows) [
    "BASH_PATH=${bash}/bin/bash"
  ];

  installPhase = ''
    mkdir -p $out
    make install DESTDIR=$bin
    mkdir -p $bin/share/fish/vendor_conf.d
    echo "eval ($bin/bin/direnv hook fish)" > $bin/share/fish/vendor_conf.d/direnv.fish
  '';

  meta = with stdenv.lib; {
    description = "A shell extension that manages your environment";
    homepage = https://direnv.net;
    license = licenses.mit;
    maintainers = with maintainers; [ zimbatm ];
  };
}
