# Development

Setup a go environment https://golang.org/doc/install

> go >= 1.24 is required

Clone the project:

    $ git clone git@github.com:direnv/direnv.git

Build by just typing make:

    $ cd direnv
    $ make

Test the projects:

    $ make test

To install to /usr/local:

    $ make install

Or to a different location like `~/.local`:

    $ make install PREFIX=~/.local

## Updating gomod2nix.toml

Execute `./script/update-gomod2nix`; if you don't have nix locally, can
do so via a docker container like so:

    $ docker run -it --platform linux/amd64 -v $(pwd):/workdir nixos/nix /bin/sh -c "cd workdir && ./script/update-gomod2nix"
