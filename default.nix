with import <nixpkgs> {};
let
  rstrip = s:
    let
      inherit (builtins) substring;
      len = builtins.stringLength s;
      suffix = substring (len - 1) 1 s;
      ws = [ " " "\r" "\n" "\t" ];
    in
      if len > 0 && builtins.any (char: char == suffix) ws then
        rstrip (substring 0 (len - 1) s)
      else
        s;
  readVersion = f:
    rstrip (builtins.readFile f);
in
buildGoPackage rec {
  version = readVersion ./version.txt;
  name = "direnv-${version}";
  goPackagePath = "github.com/zimbatm/direnv";

  src = ./.;

  meta = with stdenv.lib; {
    homepage = http://direnv.net;
    description = "path-dependent environments";
    maintainers = with maintainers; [ zimbatm ];
    license = licenses.mit;
    platforms = go.meta.platforms;
  };
}

