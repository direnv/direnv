pub mod config;
pub mod stdlib;
pub mod env;

pub struct VersionInfo {
    pub major: String,
    pub minor: String,
    pub patch: String,
    pub pre_release: Option<String>,
}

pub fn version() -> VersionInfo {
    macro_rules! env_str {
        ($name:expr) => { env!($name).to_string() }
    }
    macro_rules! option_env_str {
        ($name:expr) => { option_env!($name).map(|s| s.to_string()) }
    }
    VersionInfo {
        major: env_str!("CARGO_PKG_VERSION_MAJOR"),
        minor: env_str!("CARGO_PKG_VERSION_MINOR"),
        patch: env_str!("CARGO_PKG_VERSION_PATCH"),
        pre_release: option_env_str!("CARGO_PKG_VERSION_PRE"),
    }
}
