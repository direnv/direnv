Installing direnv
=================

Installing direnv is a two-step process. First install the binary into your path
and then activate the shell hook. Various methods are availble for both steps so
choose accordingly to your setup.

Installing the binary
---------------------

On OSX
^^^^^^

If you're using `Homebrew <http://brew.sh>`_ then installing direnv is just a
single command away::

    $ brew install direnv

Otherwise you can always take the binary from below.

On Arch Linux
^^^^^^^^^^^^^

direnv is available as an AUR. See https://aur.archlinux.org/packages/direnv/

If you're using the `yaourt pacman frontend <http://archlinux.fr/yaourt-en>`_ then
installing is just a single command away::

    $ yaourt -S direnv

On NixOS
^^^^^^^^

direnv is available in the `nixpkgs repository <http://nixos.org/nixpkgs/>`_. 
To install type::

    $ nix-env -i direnv


Build from source
-----------------

direnv depends on `Go <http://golang.org>`_ to compile properly. Once installed
building direnv is quite easy::

    $ git clone https://github.com/zimbatm/direnv.git
    $ cd direnv
    $ make install

By default direnv will be installed in /usr/local. It's possible to change the
destination by setting the DESTDIR environment varialbe. eg: 
``make install DESTDIR=/opt/direnv``

Binary builds
-------------

Get the binary for your OS below and put it in your path.

.. tip:: don't forget to make the file executable. eg: ``chmod +x direnv``

* `Darwin/386 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.darwin-386>`_
* `Darwin/amd64 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.darwin-amd64>`_
* `FreeBSD/386 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.freebsd-386>`_
* `FreeBSD/amd64 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.freebsd-amd64>`_
* `Linux/386 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.linux-386>`_
* `Linux/amd64 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.linux-amd64>`_
* `Linux/arm <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.linux-arm>`_
* `Windows/386 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.windows-386>`_
* `Windows/amd64 <http://zimbatm.s3.amazonaws.com/direnv/direnv2.2.1.windows-amd64>`_


Shell configuration
-------------------

Once the hook is installed, don't forget to restart your shell.

Bash
^^^^

::

    $ echo 'eval "direnv hook bash"' >> ~/.bashrc

Zsh
^^^

::

    $ echo 'eval "direnv hook zsh"' >> ~/.zshrc

Fish shell
^^^^^^^^^^

::

    $ echo 'eval (direnv hook fish)' >> ~/.config/fish/config.fish




