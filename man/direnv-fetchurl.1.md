DIRENV-FETCHURL 1 "2019" direnv "User Manuals"
==============================================

NAME
----

direnv fetchurl - Fetch a URL to disk

SYNOPSIS
--------

*direnv fetchurl* <url> [<integrity-hash>]

DESCRIPTION
-----------

This command downloads the given URL into a fixed disk location, based on the
content of the retrieved file.

This has been introduced to avoid a dependency on `curl` or `wget`, while also
promoting a more secure way to fetch data from the Internet. Use this instead
of `curl <url> | sh`.

This command has two modes of operation:

1. Just pass the URL to discover the integrity hash.
2. Pass the URL *and* the integrity hash to get back the on-disk location.

Since the on-disk location is based on the hash, it also acts as a cache. One
implication of this design is that URLs must be stable and always return the
same content.

Downloaded files are marked as read-only and executable so it can also be used
to fetch and execute static binaries.

OPTIONS
-------

<url>
    A HTTP URL that returns content on HTTP GET requests. 301 and other
    redirects are followed.

<integrity-hash>
    When passed, the integrity of the retrieved content will be validated
    against the given hash. The hash encoding is based on the SRI W3C
    specification (see https://www.w3.org/TR/SRI/ ).

OUTPUT
------

*direnv fetchurl* outputs different content based on the argument.

If the `integrity-hash` is being passed, it will output the path to the
on-disk location, if the retrieval was successful.

If only the `url` is being passed, it will output the hash as well as some
human-readable instruction. If stdout is not a tty, only the hash will be
displayed.

EXAMPLE
-------

    $ ./direnv fetchurl https://releases.nixos.org/nix/nix-2.3.7/install
    Found hash: sha256-7Gxl5GzI9juc/U30Igh/pY+p6+gj5Waohfwql3jHIds=

    Invoke fetchurl again with the hash as an argument to get the disk location:

      direnv fetchurl "https://releases.nixos.org/nix/nix-2.3.7/install" "sha256-7Gxl5GzI9juc/U30Igh/pY+p6+gj5Waohfwql3jHIds="
      #=> /home/zimbatm/.cache/direnv/cas/sha256-7Gxl5GzI9juc_U30Igh_pY+p6+gj5Waohfwql3jHIds=

ENVIRONMENT VARIABLES
---------------------

**XDG_CACHE_HOME**
    This variable is used to select the on-disk location of the fetched URLs
    as `$XDG_CACHE_HOME/direnv/cas`. If **XDG_CACHE_HOME** is unset or empty,
    defaults to `$HOME/.cache`.

COPYRIGHT
---------

MIT licence - Copyright (C) 2019 @zimbatm and contributors

SEE ALSO
--------

direnv-stdlib(1)
