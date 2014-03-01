The command-line
================

The ``direnv`` executable is a command dispatcher. Below are all the
sub-commands available.

``direnv <command> [...args]``

Public commands
---------------

.. _direnv_allow:

``direnv allow [path_to_rc]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Grants direnv to load the given .envrc

.. _direnv_deny:

``direnv deny [path_to_rc]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Revokes the auhorization of a given .envrc

.. _direnv_edit:

``direnv edit [path_to_rc]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Opens PATH_TO_RC or the current .envrc into an $EDITOR and allow the file to be loaded afterwards.

.. _direnv_exec:

``direnv exec [dir] <command> [...args]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Executes a command after loading the first .envrc found in [dir]. Unless specified, the directory is the parent directory of the command.

.. _direnv_help:

``direnv help [show_private]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Shows the available commands of the executable. If anything is passed as an argument, the below private commands are also displayed.


.. _direnv_hook:

``direnv hook <shell>``
^^^^^^^^^^^^^^^^^^^^^^^

When evaluated inside the target shell, setups direnv as the shell extension.

Current shell supported: fish, bash, zsh

.. _direnv_reload:

``direnv reload``
^^^^^^^^^^^^^^^^^

Triggers an env reload

.. _direnv_status:

``direnv status``
^^^^^^^^^^^^^^^^^

Prints some debug status informations

Private commands
----------------

.. note:: these commands are used internally and there is no stability guarantee
          in regards to their interfaces.


.. _direnv_apply_dump:

``direnv apply_dump <file>``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Accepts a filename containing `direnv dump` output and generates a series of bash export statements to apply the given env.

.. _direnv_dotenv:

``direnv dotenv [shell] [path_to_dotenv]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Transforms a .env file to evaluatable `export KEY=PAIR` statements.

.. _direnv_dump:

``direnv dump``
^^^^^^^^^^^^^^^

Used to export the inner bash state at the end of execution.

.. _direnv_expand_path:

``direnv expand_path <path> [rel_to]``
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Transforms a PATH to an absolute path to REL_TO or $PWD.

.. _direnv_export:

``direnv export <shell>``
^^^^^^^^^^^^^^^^^^^^^^^^^

Loads an .envrc and prints the diff in terms of exports.

.. note:: this is what is executed and evaluated on every prompt.

.. _direnv_stdlib:

``direnv stdlib``
^^^^^^^^^^^^^^^^^

Displays the stdlib available in the .envrc execution context.

.. _direnv_version:

``direnv version``
^^^^^^^^^^^^^^^^^^

Prints the current direnv version.

