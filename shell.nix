with import <nixpkgs> {};

mkShell {
  buildInputs = [
    go
    gopls
  ];
}
