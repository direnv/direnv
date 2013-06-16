---
title: direnv 2 released !
layout: default
---

Finally direnv 2.0.0 is available to everyone !

First, I would like to express a special thanks to
[Peter Waller](https://github.com/pwaller) for helping me bring this release
to light.

About
-----

direnv is a shell extension that allows you to change your environment
variables dynamically depending on your path. You can use it load specific
environment variables for your project or change your PATH to load different
versions of an interpreter. In that sense it's a bit like [RVM](https://rvm.io/)
or other version managers except that it's language-agnostic. It's just one
more tool that you can add to your unix toolbelt.

Changes
-------

direnv 2 is a full rewrite in [Go](http://golang.org) to remove runtime
dependencies and provide a single binary for deploys, execution should also be
faster as a nice side-effect. In some ways I liked the simplicity of the
earlier version but this rewrite also brings better error handling and shell
escaping. I hope it will make the whole experience much better ! Support is
still limited to the BASH and ZSH shells but extensions for FISH or other
shells should also be easier to write.

Most of the changes are backward-compatible and upgrading from direnv 1.x
should be a drop-in replacement although you probably need to restart your
shells.

The biggest change is the introduction of a security model. .envrc are not
automatically loaded unless you allow them. When you see an authorization
issue, review the .envrc and run `direnv allow` to authorize the loading.
Or better yet, use the new `direnv edit` command to open the file in your
$EDITOR and automatically allow it on exit.

Example:

```bash
$ cd path/to/project
.envrc is not allowed
$ direnv allow
direnv: reloading
direnv: loading ~/.direnvrc
direnv: loading .envrc
direnv export: +BUNDLE_BIN +CPATH +GEM_HOME +LD_LIBRARY_PATH +LIBRARY_PATH
+PKG_CONFIG_PATH ~MANPATH ~PATH
```

Installation
------------

There are multiple ways to install direnv, pick the one you prefer:

One option is to `git clone git://github.com:zimbatm/direnv.git`
and run `make install`. Or you can use the `go get github.com/zimbatm/direnv`
command. For Mac OSX users you can `brew install direnv` (just waiting on the
approval of the [pull request](https://github.com/mxcl/homebrew/pull/20540).

If you don't have Go installed on your system you can also fetch a single
binary and install it in your PATH. I guess the windows binaries also need to
be renamed to add the .exe extension and the other binaries need to be marked
as executable with `chmod +x path/to/direnv`.

Example:

    sudo wget -O /usr/local/bin/direnv http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.linux-amd64
    sudo chmod +x /usr/local/bin/direnv

* [Darwin/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.darwin-386)
* [Darwin/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.darwin-amd64)
* [FreeBSD/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.freebsd-386)
* [FreeBSD/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.freebsd-amd64)
* [Linux/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.linux-386)
* [Linux/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.linux-amd64)
* [Linux/arm](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.linux-arm)
* [Windows/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.windows-386)
* [Windows/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0.windows-amd64)

Finally, insert the following line at the very end of either your .bashrc or
.zshrc: `eval "$(direnv hook $0)"` and restart your shell.

Happy hacking !
