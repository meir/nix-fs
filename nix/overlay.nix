inputs:
{ pkgs, lib, config, ... }:
let
  new-state = pkgs.writeText "nix-fs.json" (builtins.toJSON {
    version = 1;
    time = "2000-01-01T00:00:00Z";
    locations = lib.mapAttrsToList (name: file: {
      source = if (file.source != null) then
        file.source
      else if (file.text != null) then
        pkgs.writeTextFile {
          name = name + "-content";
          text = file.text;
        }
      else
        throw "Either 'source' or 'text' must be provided for file '${name}'";

      destination = config.nix-fs.home + "/" + name;
    }) config.nix-fs.files;
  });
in
{
  options.nix-fs = with lib; {
    package = mkOption {
      default = pkgs.nix-fs;
      type = types.package;
    };
    home = mkOption {
      type = types.str;
      default = null;
      description = "The home directory where files will be managed.";
    };
    files = mkOption {
      type = types.attrsOf (types.submodule {
        options = {
          source = mkOption {
            type = types.nullOr types.path;
            default = null;
            description = "Source path of the file or directory to manage.";
          };

          text = mkOption {
            type = types.nullOr types.str;
            default = null;
            description = "Content to write to the destination file. If set, 'source' is ignored.";
          };
        };
      });
      default = { };
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
