{ pkgs, lib }:
pkgs.buildGoApplication rec {
  pname = "nix-fs";
  version = "1.0";

  pwd = ./..;
  src = ./..;

  modules = ../gomod2nix.toml;
}
