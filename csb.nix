with import <nixpkgs> {};
 
stdenv.mkDerivation {
    name = "csb";
    buildInputs = [
        go-task
    ];
    shellHook = ''
        export PATH="$PWD/node_modules/.bin/:$PATH"
    '';
}
