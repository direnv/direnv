---
title: direnv 2 is coming
layout: default
---

After a long time, version 2.0 of direnv if finally approaching. I'm looking
for beta-testers so if you want to try it out see the install section below
and leave your comments on the [mailing list](mailto:direnv@librelist.com) or
as [github issues](http://github.com/zimbatm/direnv/issues).

About
-----

direnv 2 is a full rewrite of the tool in [Go](http://golang.org). There are
some operations that we can't do in bash and until now direnv was depending on
ruby for these bits. Now that the tool is rewritten in Go we just need to
distribute a single binary to support any operations (plus Go is good at
cross-compiling). People with rotating disks should also see startup
improvements since ruby was previously loading a tons of files on each
invokation.

On top of this, direnv has also seen improvements in terms of shell escaping,
eval loop robustness and security. Now all .envrc are denied by default unless
you authorize to load them with the `direnv allow` command or after you have
opened it with `direnv edit` (which opens your $EDITOR).

Installation
------------

Once direnv will become stable it will be available to Homebrew users as
`brew install direnv` but for now you'll have to install it by hand (help for
other operating systems welcome !)

One option is to `git clone --branch go1 git://github.com:zimbatm/direnv.git`
and run `make install`. Or you can use the `go get github.com/zimbatm/direnv`
command.

If you don't have Go installed on your system you can also fetch a single
binary and install it in your PATH. I guess the windows binaries also need to
be renamed to add the .exe extension.

* [Darwin/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.darwin-386)
* [Darwin/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.darwin-amd64)
* [FreeBSD/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.freebsd-386)
* [FreeBSD/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.freebsd-amd64)
* [Linux/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.linux-386)
* [Linux/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.linux-amd64)
* [Linux/arm](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.linux-arm)
* [Windows/386](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.windows-386)
* [Windows/amd64](http://zimbatm.s3.amazonaws.com/direnv/direnv2.0.0-rc.1.windows-amd64)

Finally, insert the following line at the very end of either your .bashrc or
.zshrc: `eval "$(direnv hook $0)"` and restart your shell.

Happy hacking !
