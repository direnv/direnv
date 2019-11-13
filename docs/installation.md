# Installation

The installation has two parts.

1. Install the package or binary, which is presented in this document
2. [hook into your shell](hook.md).

## From system packages

direnv is packaged for a variety of systems:

* [Fedora](https://apps.fedoraproject.org/packages/direnv)
* [Arch AUR](https://aur.archlinux.org/packages/direnv/)
* [Debian](https://packages.debian.org/search?keywords=direnv&searchon=names&suite=all&section=all)
* [Gentoo go-overlay](https://github.com/Dr-Terrible/go-overlay)
* [NetBSD pkgsrc-wip](http://www.pkgsrc.org/wip/)
* [NixOS](https://nixos.org/nixos/packages.html#direnv)
* [OSX Homebrew](http://brew.sh/)
* [openSUSE](https://build.opensuse.org/package/show/openSUSE%3AFactory/direnv)
* [MacPorts](https://www.macports.org/)
* [Ubuntu](https://packages.ubuntu.com/search?keywords=direnv&searchon=names&suite=all&section=all)
* [GNU Guix](https://www.gnu.org/software/guix/)
* [Snap](https://snapcraft.io/direnv)

See also:

[![Packaging status](https://repology.org/badge/vertical-allrepos/direnv.svg)](https://repology.org/metapackage/direnv)

## From binary builds

To install binary builds you can run this bash installer:

```
curl -sfL https://direnv.net/install.sh | bash
```

Binary builds for a variety of architectures are also available for
[each release](https://github.com/direnv/direnv/releases).

Fetch the binary, `chmod +x direnv` and put it somewhere in your PATH.

## Compile from source

See the [Development](development.md) page.

# Next step

[hook installation](hook.md)
