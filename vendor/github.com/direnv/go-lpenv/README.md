# LookPathEnv

This is a fork of the stdlib os/exec.LookPath function.

It does two things differently:

* pass the environment variables by argument instead of getting it from the
  process global context.
* fix a potential security issue on Windows.

That's it.
