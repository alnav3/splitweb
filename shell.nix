{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs_22
    pkgs.go
    pkgs.air
    pkgs.tailwindcss
    pkgs.templ
  ];
  shellHook = ''
    if ! npm list tailwindcss | grep -q tailwindcss && ! npm list @catppuccin/tailwindcss | grep -q tailwindcss; then
      echo "Installing tailwindcss and @catppuccin/tailwindcss..."
      npm install -D tailwindcss
      npm install -D @catppuccin/tailwindcss
    else
      echo "tailwindcss and @catppuccin/tailwindcss are already installed."
    fi
  '';
}
