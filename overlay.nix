{ pkgs, lib, config, ... }:
{
  options.nix-fs.files = with lib; mkOption {
    type = types.listOf types.submodule {
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

  };
}
