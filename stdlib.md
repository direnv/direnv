---
layout: default
---

direnv-stdlib(1) -- the ".envrc" stdlib
=======================================

## SYNOPSIS

`direnv stdlib`

## DESCRIPTION

Outputs a bash script called the *stdlib*. The following commands are included in that script and loaded in the context of an ".envrc". Additionnaly to that, it also loads the file in "~/.direnvrc" if it exists.

## STDLIB

  * `has` `<command>`:
    Returns 0 if the `<command>` is available. Returns 1 otherwise. It can be a binary in the PATH or a shell function.

Example:

    if has curl; then
      echo "Yes we do"
    fi

* `expand_path` `<rel_path>` [`<relative_to>`]:
    Outputs the absolute path of `<rel_path>` relaitve to `<relative_to>` or the current directory.

Example:

    cd /usr/local/games
    expand_path ../foo
    # 1> /usr/local/foo

* `dotenv` [`<dotenv_path>`]: 
    Loads a ".env" file into the current environment

* `user_rel_path` `<abs_path>`:
    Outputs a path relative to the user's home if possible.

Example:

    user_rel_path $HOME/my/project
    # 1> ~/my/project

* `find_up` `<filename>`:
    Outputs the path of `<filename>` when searched from the current directory up to /. Returns 1 if the file has not been found.

Example:

    cd /usr/local/my
    mkdir -p project/foo
    touch bar
    cd project/foo
    find_up bar
    # 1> /usr/local/my/bar

* `source_env` `<file_or_dir_path>`:
    Loads another ".envrc" either by specifying it's path or filename.

* `source_up` [`<filename>`]:
    Loads another ".envrc" if found with the `find_up` command.

* `PATH_add` `<path>`:
    Prepends the expanded `<path>` to the PATH environment variable. It prevents a common mistake where PATH is replaced by only the new <path>.

Example:

    pwd
    # 1> /home/user/my/project
    PATH_add bin
    echo $PATH
    # 1> /home/user/my/project/bin:/usr/bin:/bin

* `path_add` `<varname>` `<path>`:
    Works like `PATH_add` except that it's for an arbitrary <varname>.

* `load_prefix` `<prefix_path>`: 
    Expands some common path variables for the given `<prefix_path>` prefix.

Variables set:

    CPATH
    LD_LIBRARY_PATH
    LIBRARY_PATH
    MANPATH
    PATH
    PKG_CONFIG_PATH

Example:

    ./configure --prefix=$HOME/rubies/ruby-1.9.3
    make && make install
    # Then in the .envrc
    load_prefix ~/rubies/ruby-1.9.3

* `layout` `<type>`:
    A semantic dispatch used to describe common project layouts.

* `layout ruby`:
    Sets the GEM_HOME environment variable to "$PWD/.direnv/ruby/RUBY_VERSION". This forces the installation of any gems into the project's sub-folder.
    If you're using bundler it will create wrapper programs that can be invoked directly instead of using the `bundle exec` prefix.

* `layout python`: 
    Creates and loads a virtualenv environment under "$PWD/.direnv/virtualenv". This forces the installation of any egg into the project's sub-folder.

* `layout node`:
    Adds "$PWD/node_modules/.bin" to the PATH environment variable.

* `layout go`: 
    Sets the GOPATH environment variable to the current directory.

* `use` `<program_name>` [`<version>`]:
    A semantic command dispatch intended for loading external dependencies into the environment.

Example:

    use_ruby() {
      echo "Ruby $1"
    }
    use ruby 1.9.3
    # 1> Ruby 1.9.3

* `use rbenv`: 
    Loads rbenv which add the ruby wrappers available on the PATH.

* `rvm` [...]: 
    Should work just like in the shell if you have rvm installed.

