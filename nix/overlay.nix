{ pkgs, lib, config, ... }:
let
  home-folder = lib.getEnv "HOME";
  new-state = pkgs.writeText "nix-fs.json" (pkgs.lib.toJSON {
    version = 1;
    time = pkgs.lib.currentTimeString;
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
  options.nix-fs.files = with lib; mkOption {
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

  config = {
    system.activationScripts = {
      nix-fs = {
        deps = [
          "etc"
          "nix"
          "users"
        ];
        text =
          ''
            ${pkgs.nix-fs}/bin/nix-fs --state-file ${new-state} --old-state-file /etc/nix-fs.json
          '';
      };
    };
  };
}
