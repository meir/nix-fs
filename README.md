# Nix-FS

A small Nix flake for filesystem/dotfiles management using Go.
This project is supposed to provide an alternative to HomeManager to manage your filesystem/dotfiles.
It can be used to include `.config` files in your Nix flake for example and link them to your actual `.config` folder.

## How to install/use

First add the flake to your inputs as followed:
```nix
    inputs = {
        # ...
        nix-fs.url = "github:meir/nix-fs";
    };
```

After that, add it as a module in your `nixosConfiguration`:
```nix
    nixpkgs.lib.nixosSystem {
        # ...
        modules = [
            # ...
            nix-fs.nixosModules.nix-fs
        ];
    }
```

Then you have to specify your home folder in your config as such:
```nix
    nix-fs.home = "/home/your-username";
```

And finally you can specify files like such:
```nix
    nix-fs.files = {
        ".config/kitty/kitty.conf".text = ''
            font_size 12
        '';

        ".config/eww".source = ./eww;
    };
```

This works both with files and folders.

## How it works

The Go program has a State file in JSON saved in `/etc/nix-fs.json`.
In this file is the current state saved of the program, aka all the files that are (supposed to be) linked by Nix-FS.
If the actual files deviate from the current state, Nix-FS will try to fix this.
> If a file already exists in place without it being linked, Nix-FS will not delete the file or give an error, it will just give a warning.

Once a new state is given, Nix-FS compares it to the current state and goes through a list of actions that have to be taken and eventually saves the new state in `/etc/nix-fs.json`.

