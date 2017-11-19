use std::path::{PathBuf};
use std::option::{Option};

const BASH_PATH: Option<&'static str> = option_env!("BASH_PATH");

#[derive(Debug)]
pub struct Config {
    /// The location of the users's 'home' directory. OS-dependent.
    home_path: PathBuf,
}
