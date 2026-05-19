# Don’t check in `.envrc` files

`.envrc` files are intended to be a developer’s own personal means of customizing their environment for a project.
Therefore, as much as the idea of checking in a `.envrc` file may look good, do not be tempted.
A checked-in `.envrc` file makes for a figurative hoop to jump through before contributing to a project.
However, a checked in `.envrc.example` is a kindness!

More so, `.envrc` should be added to `.gitignore` as a slight safety measure,
to reduce the risk of `.envrc` with secrets being accidentally committed.

If you see a checked in `.envrc` file in a public open source project,
consider contributing a change that moves it to `.envrc.example`
and adds a `.envrc` line to `.gitignore`.

Some projects that have repented (add yours!):

- https://github.com/NixOS/nixpkgs/pull/325793
- https://codeberg.org/Aviac/codeberg-cli/pulls/240
- https://github.com/firezone/firezone/pull/6496
- https://github.com/riverqueue/river/pull/1230
- https://github.com/mattermost/mattermost/pull/36567

Some `.envrc` check-in rejections:

- help us find some!
