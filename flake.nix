{
  description = "a small nix flake for filesystem/dotfiles management in go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs =
    { nixpkgs, gomod2nix, pre-commit-hooks, ... }@inputs:
    let
      inherit (nixpkgs) lib;
      eachSystem = function: lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
      ]
      (system:
        function (import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        })
      );
    in
    rec {
      packages = eachSystem (pkgs:
        {
          nix-fs = pkgs.callPackage ./nix/nix-fs.nix { };
        }
      );

      nixosModules.nix-fs = args: import ./nix/overlay.nix (args // {
        nix-fs = packages.${args.system}.nix-fs;
      });

      devShells = eachSystem (pkgs: {
        default = pkgs.callPackage ./nix/shell.nix { inherit pre-commit-hooks pkgs; };
      });
    };
}
