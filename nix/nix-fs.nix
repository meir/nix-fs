{ pkgs, lib }:
pkgs.buildGoApplication rec {
  pname = "nix-fs";
  version = lib.readFile ./VERSION;

  pwd = ./..;
  src = ./..;

  modules = ./gomod2nix.toml;
}
