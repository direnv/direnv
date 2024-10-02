# Installation

The installation has two parts.

1. Install the package or binary, which is presented in this document
2. [hook into your shell](hook.md).

## From system packages

direnv is packaged for a variety of systems:

* [Fedora](https://src.fedoraproject.org/rpms/direnv)
* [Arch Linux](https://archlinux.org/packages/extra/x86_64/direnv/)
* [Debian](https://packages.debian.org/search?keywords=direnv&searchon=names&suite=all&section=all)
* [Gentoo Guru](https://wiki.gentoo.org/wiki/Project:GURU/Information_for_End_Users)
* [NetBSD pkgsrc-wip](http://www.pkgsrc.org/wip/)
* [NixOS](https://search.nixos.org/options?query=programs.direnv)
* [macOS Homebrew](https://formulae.brew.sh/formula/direnv#default)
* [openSUSE](https://build.opensuse.org/package/show/openSUSE%3AFactory/direnv)
* [MacPorts](https://ports.macports.org/port/direnv/)
* [Ubuntu](https://packages.ubuntu.com/search?keywords=direnv&searchon=names&suite=all&section=all)
* [GNU Guix](https://packages.guix.gnu.org/search/?query=direnv)
* [Windows](https://learn.microsoft.com/en-us/windows/package-manager/winget/)

See also:

[![Packaging status](https://repology.org/badge/vertical-allrepos/direnv.svg)](https://repology.org/metapackage/direnv)

## From binary builds

To install binary builds you can run this bash installer:

```sh
curl -sfL https://direnv.net/install.sh | bash
```

Binary builds for a variety of architectures are also available for
[each release](https://github.com/direnv/direnv/releases).

Fetch the binary, `chmod +x direnv` and put it somewhere in your `PATH`.

## Compile from source

See the [Development](development.md) page.

# Next step

[hook installation](hook.md)
