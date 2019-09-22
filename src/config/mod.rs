#![deny(warnings)]

use toml;
//use std::env::{current_exe};
use std::fs::File;
use std::io::Read;
use std::option::Option;
use std::path::PathBuf;

const BASH_PATH: Option<&'static str> = option_env!("BASH_PATH");

#[derive(Debug, Deserialize)]
pub struct Config {
    global: ConfigGlobal,
    logging: ConfigLogging,
}

#[derive(Debug, Deserialize)]
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

#[derive(Debug, Deserialize)]
pub struct ConfigLogging {
    /// The log format for error
    error_format: String,

    /// The log format for information
    info_format: String,
}

pub fn load(path: PathBuf) -> std::io::Result<Config> {
    //ConfigGlobal.self_path = current_exe()
    let mut file = File::open(&path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    let mut config: Config = toml::from_str(&contents)?;

    if config.global.bash_path.to_string_lossy() == "" {
        match BASH_PATH {
            Some(bash_path) => config.global.bash_path = PathBuf::from(bash_path),
            None => (),
        }
    }

    return Ok(config);
}
