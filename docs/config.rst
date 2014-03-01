Configuration
=============

Some environment variables can be set globally to change the behavior of
direnv.

DIRENV_CONFIG
-------------

Sets the configuration directory where direnv will put it's allow files.
If not set it will default to $XDG_CONFIG_HOME/direnv or ~/.config/direnv.

DIRENV_BASH
-----------

If set, this will be the bash executable used to evaluate the .envrc files.
Otherwise bash is looked up in the PATH.

DIRENV_LOG_FORMAT
-----------------

If set, this will be the format used to output direnv's log messages.
By default it is set to "direnv: %s".

.. note:: don't forget to include the %s inside the string.

