inputs:
{ pkgs, lib, config, ... }:
let
  home-folder = lib.getEnv "HOME";
  new-state = pkgs.writeText "nix-fs.json" (pkgs.lib.toJSON {
    version = 1;
    time = lib.currentTimeString;
    locations = lib.mapAttrsToList (name: file: {
      source = if file.source != null then
        file.source
      else if file.text != null then
        pkgs.writeTextFile {
          name = name + "-content";
          text = file.text;
        }
      else
        throw "Either 'source' or 'text' must be provided for file '${name}'";

      destination = home-folder + "/" + name;
    }) config.nix-fs.files;
  });
in
{
  options.nix-fs = with lib; {
    package = mkOption {
      default = pkgs.nix-fs;
      type = types.package;
    };
    files = mkOption {
      type = types.attrsOf types.submodule {
        options = {
          source = mkOption {
            type = types.nullOr types.path;
            description = "Source path of the file or directory to manage.";
          };

          text = mkOption {
            type = types.nullOr types.str;
            description = "Content to write to the destination file. If set, 'source' is ignored.";
          };
        };
      };
      default = [ ];
      description = "List of files to manage with nix-fs.";
    };
  };

  config = {
    nixpkgs.overlays = [
      (final: prev: {
        nix-fs = inputs.self.packages.${final.system}.nix-fs;
      })
    ];

    system.activationScripts = {
      nix-fs = {
        deps = [
          "etc"
          "users"
        ];
        text =
          ''
            ${config.nix-fs.package}/bin/nix-fs --state-file ${new-state} --old-state-file /etc/nix-fs.json
          '';
      };
    };
  };
}
