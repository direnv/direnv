Installation
============

Requirements
------------

direnv is a go program and thus compiles to most POSIX systems and Windows.

direnv is built for all go build targets although Windows is probably
less likely to work. One of the following compatible shell is also necessary: 
bash, fish or zsh.


Installing direnv is a two-step process. First install the binary into your path
and then activate the shell hook. Various methods are availble for both steps so
choose accordingly to your setup.

Installing the binary
---------------------

OSX: Install from Homebrew
^^^^^^^^^^^^^^^^^^^^^^^^^^

If you're using `Homebrew <http://brew.sh>`_ then installing direnv is just a
single command away::

    $ brew install direnv

Otherwise you can always take the binary from below.

Arch Linux: Install from AUR
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

direnv is available as an AUR. See https://aur.archlinux.org/packages/direnv/

If you're using the `yaourt pacman frontend <http://archlinux.fr/yaourt-en>`_ then
installing is just a single command away::

    $ yaourt -S direnv

NixOS: Install from nixpkgs
^^^^^^^^^^^^^^^^^^^^^^^^^^^

direnv is available in the `nixpkgs repository <http://nixos.org/nixpkgs/>`_. 
To install type::

    $ nix-env -i direnv

Other: Binary builds
^^^^^^^^^^^^^^^^^^^^

Get the binary for your OS below and put it in your path.

.. tip:: don't forget to make the file executable. eg: ``chmod +x direnv``

=======  =====  ==============
OS       Arch   Download
=======  =====  ==============
Darwin   386    Darwin_386_
Darwin   amd64  Darwin_amd64_
FreeBSD  386    FreeBSD_386_
FreeBSD  amd64  FreeBSD_amd64_
Linux    386    Linux_386_
Linux    amd64  Linux_amd64_
Linux    arm    Linux_arm_
Windows  386    Windows_386_
Windows  amd64  Windows_amd64_
=======  =====  ==============

Other: Build from source
^^^^^^^^^^^^^^^^^^^^^^^^

direnv depends on `Go <http://golang.org>`_ to compile properly. Once installed
building direnv is quite easy::

    $ git clone https://github.com/zimbatm/direnv.git
    $ cd direnv
    $ make install

By default direnv will be installed in /usr/local. It's possible to change the
destination by setting the DESTDIR environment varialbe. eg: 
``make install DESTDIR=/opt/direnv``

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




