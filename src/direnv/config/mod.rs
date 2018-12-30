use std::path::{PathBuf};
use std::option::{Option};

const BASH_PATH: Option<&'static str> = option_env!("BASH_PATH");

#[derive(Debug)]
pub struct Config {
    global: ConfigGlobal,
    logging: ConfigLogging,
}

#[derive(Debug)]
pub struct ConfigGlobal {
    /// The absolute location of direnv.
    self_path: PathBuf,

    /// The location of bash
    bash_path: PathBuf,

    /// The location of the direnv config directory
    config_dir: PathBuf,

    /// The location of the direnv state directory
    state_dir: PathBuf,
}

#[derive(Debug)]
pub struct ConfigLogging {
    /// The log format for error
    error_format: &'static str,

    /// The log format for information
    info_format: &'static str,
}
