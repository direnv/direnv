{ pkgs ? import <nixpkgs> {} }:
with pkgs;
buildGoPackage rec {
  version = lib.fileContents ./version.txt;
  name = "direnv-${version}";
  goPackagePath = "github.com/zimbatm/direnv";

  src = lib.cleanSource ./.;

  postConfigure = "cd $NIX_BUILD_TOP/go/src/$goPackagePath";

  buildPhase = "make BASH_PATH=${bash}/bin/bash";

  installPhase = ''
    mkdir -p $out
    make install DESTDIR=$bin
    mkdir -p $bin/share/fish/vendor_conf.d
    echo "eval ($bin/bin/direnv hook fish)" > $bin/share/fish/vendor_conf.d/direnv.fish
  '';

  meta = with stdenv.lib; {
    homepage = http://direnv.net;
    description = "A shell extension that manages your environment";
    maintainers = with maintainers; [ zimbatm ];
    license = licenses.mit;
    platforms = go.meta.platforms;
  };
}
