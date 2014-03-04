.. _direnv-command:

**************
direnv command
**************

Synopsis
========

direnv
    <command> [...args]

Description
===========

direnv is a shell extension that changes the environment variables depending on
the current directory. By adding an ".envrc" bash script to any folder, it's
exported variables are then loaded in the current shell when entering the
directory, and unloaded when exitting.

The ``direnv`` executable is a command dispatcher. Below are all the
sub-commands available.

Public commands
===============

.. _direnv_allow:

direnv allow [path_to_rc]
-------------------------

Grants direnv to load the given .envrc

.. _direnv_deny:

direnv deny [path_to_rc]
------------------------

Revokes the auhorization of a given .envrc


.. _direnv_edit:

direnv edit [path_to_rc]
-------------------------

Opens PATH_TO_RC or the current .envrc into an $EDITOR and allow the file to be loaded afterwards.

.. _direnv_exec:

direnv exec [dir] <command> [...args]
-------------------------------------

Executes a command after loading the first .envrc found in [dir]. Unless specified, the directory is the parent directory of the command.

.. _direnv_help:

direnv help [show_private]
--------------------------

Shows the available commands of the executable. If anything is passed as an argument, the below private commands are also displayed.

.. _direnv_hook:

direnv hook <shell>
-------------------

When evaluated inside the target shell, setups direnv as the shell extension.

Current shell supported: fish, bash, zsh

.. _direnv_reload:

direnv reload
-------------

Triggers an env reload

.. _direnv_status:

direnv status
-------------

Prints some debug status informations

Private commands
================

.. note:: these commands are used internally and there is no stability guarantee
          in regards to their interfaces.


.. _direnv_apply_dump:

direnv apply_dump <file>
------------------------

Accepts a filename containing `direnv dump` output and generates a series of bash export statements to apply the given env.

.. _direnv_dotenv:

direnv dotenv [shell] [path_to_dotenv]
--------------------------------------

Transforms a .env file to evaluatable `export KEY=PAIR` statements.

.. _direnv_dump:

direnv dump
-----------

Used to export the inner bash state at the end of execution.

.. _direnv_expand_path:

direnv expand_path <path> [rel_to]
----------------------------------

Transforms a PATH to an absolute path to REL_TO or $PWD.

.. _direnv_export:

direnv export <shell>
---------------------

Loads an .envrc and prints the diff in terms of exports.

.. note:: this is what is executed and evaluated on every prompt.

.. _direnv_stdlib:

direnv stdlib
-------------

Displays the stdlib available in the .envrc execution context.

.. _direnv_version:

direnv version
--------------

Prints the current direnv version.

Files
=====

:file:`.envrc`
  A script loaded in a bash subshell. When placed inside a directory direnv will
  load it if the current directory is the same as the script or in a sub=directory.

:file:`~/.config/direnv/direnvrc`
  Your personal standard library. This file is loaded in the context of an 
  ".envrc"

Environment
===========

Some environment variables can be set globally to change the behavior of
direnv.

DIRENV_CONFIG
-------------

Sets the configuration directory where direnv will put it's allow files.
If not set it will default to :file:`$XDG_CONFIG_HOME/direnv` or 
:file:`~/.config/direnv`.

DIRENV_BASH
-----------

If set, this will be the bash executable used to evaluate the .envrc files.
Otherwise bash is looked up in the PATH.

DIRENV_LOG_FORMAT
-----------------

If set, this will be the format used to output direnv's log messages.
By default it is set to "direnv: %s".

.. note:: don't forget to include the "%s" inside the string.

See also
========

:ref:`direnv-stdlib(1) <direnv-stdlib>`

Reporting bugs
==============

Report bugs to direnv's issue tracker at
<https://github.com/zimbatm/direnv/issues>

