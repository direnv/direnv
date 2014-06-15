# Manage your rubies with direnv and ruby-install

direnv is just a shell extension that manages your environment variables
depending on the folder you live in. In this article we will explore how it
can be used in combination with
[ruby-install](https://github.com/postmodern/ruby-install) to manage and
select the version of ruby that you want to use in a project.

## The setup

First install direnv. This is the quick version on OSX + Bash:

```bash
brew install direnv
echo 'eval $(direnv hook bash)' >> .bashrc
exec $0
```

Then use [ruby-install](https://github.com/postmodern/ruby-install) to
install a couple of ruby versions. We're also creating a couple of aliases
for convenience.

```
brew install ruby-install
ruby-install ruby 1.9
ruby-install ruby 2.0
cd ~/.rubies
ln -s 1.9.3-p448 1.9.3
ln -s 1.9.3-p448 1.9
ln -s 2.0.0-p247 2.0.0
ln -s 2.0.0-p247 2.0
```

The end goal is that each project will have an ".envrc" file that contains
a descriptive syntax like `use ruby 1.9.3` to selects the right ruby version
for the project.

For that regard we are going to use a couple of commands available in the
[direnv stdlib](/stdlib.html) and expand it a bit in the ~/.direnvrc file.

Add this to the ~/.direnvrc file (you have to create it if it doesn't exist):

```bash
# Usage: use ruby <version>
#
# Loads the specified ruby version into the environent
#
use_ruby() {
  local ruby_dir=$HOME/.rubies/$1
  load_prefix $ruby_dir
  layout ruby
}
```

That's it. Now in any project you can run `direnv edit .` and add
`use ruby 1.9.3` or `use ruby 2.0` in the file like you want and direnv will
select the right ruby version when you enter the project's folder.

## A bit of explanation

The last part probably needs a bit more explanation. We make use of a couple
of commands that are part of the [stdlib](/stdlib.html) which is availabe in
the execution context of an envrc.

`use` is a command dispatch that's just there to build the
`use something something` dsl so that `use ruby <version>` will translate into
`use_ruby <version>`.

`load_prefix` will add a couple of things into the environment, notably add
`<prefix>/bin` into the PATH. This is what makes the specified ruby available.

And finally `layout ruby` who like `use` translates into the `layout_ruby`
function call. It's used to decribe common project layouts. In the stdlib, the
ruby layout will configure rubygems (with the `GEM_HOME` environment variable)
to install all the gems into the .direnv/ruby/RUBY_VERSION folder under the
project root. This is a bit similar to rvm's gemsets except that they live
inside your project's folder. It also configures bundler to install wrapper
shims into the .direnv/bin folder which allows you to invoke the commands
directly instead of prefixing your ruby programs with `bundle exec` all the
time.

## Conclusion

As you see this approach is not restricted to ruby. You could have various
versions of python installed under ~/.pythons and a `use_python` defined in
your ~/.direnvrc. Or perl, php, ...  This is the good thing about direnv, it's
not restricted to a single language.

Actually, wouldn't it be great to have all your project's dependencies
available when you enter the project folder ? Not only your ruby version but
also the exact redis or mysql or ... version that you want to use, without
having to start a VM. I think that's definitely possible using something like
the [Nix package manager](http://nixos.org/nix/), something that still needs
to be explored in a future post.

