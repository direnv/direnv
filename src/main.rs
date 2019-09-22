use direnv::stdlib;
use std::io::{self, Write};

const USAGE: &str = "
Per-directory shell environment variables

Usage:
    direnv <command> [<args>...]
    direnv [options]
Options:
    -h, --help          Display this message
    -V, --version       Print version info and exit
    --list              List installed commands
    -v, --verbose ...   Use verbose output (-vv very verbose/build.rs output)
    -q, --quiet         No output printed to stdout
    --color WHEN        Coloring: auto, always, never
Some common direnv commands are (see all commands with --list):
    init        Compile the current project
    reload      Analyze the current project and report errors, but don't build object files
    hook        Remove the target directory
See 'direnv help <command>' for more information on a specific command.
";

fn main() {
    if let Err(err) = try_main() {
        eprintln!("{}", err);
        std::process::exit(1);
    }
}

fn try_main() -> Result<(), Box<dyn std::error::Error>> {
    io::stdout().write_all(stdlib::STDLIB.as_bytes())?;

    io::stdout().write_all(USAGE.as_bytes())?;

    io::stdout().flush()?;
    Ok(())
}
