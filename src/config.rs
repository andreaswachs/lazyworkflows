use dirs::{self};
use std::{fs, path::PathBuf, io::{Write, Read}};
use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Repo {
    pub token: String,
    pub owner: String,
    pub repo: String,
}

#[derive(Deserialize, Debug)]
pub struct Config {
    pub repos: Vec<Repo>,
}

impl Config {
    /// Get the default content of the config file
    fn default_config() -> String {
        String::from("\
[[repos]]
token = 'INSERT_TOKEN_HERE'
owner = 'andreaswachs'
repo = 'lazyworkflows'
")
    }

    /// Constructs the path to the configuration directory
    fn config_dir_path() -> PathBuf {

        let mut config_dir = dirs::config_dir().expect("should get path to config dir");
        config_dir.push("lazyworkflows");
        config_dir
    }

    /// Builds on top of the path to the configuration directory with the config file name
    fn config_file_path() -> String {
        
        let mut config_dir = Config::config_dir_path();
        config_dir.push("config.toml");
        config_dir.to_str().expect("path should convert to string slice").to_string();
        config_dir.to_str().expect("should convert to string slice").to_string()
    }

    /// Creates the path to the configuration file, it trusts that you've checked that it doesn't exist already
    fn create_config_file() {

        let file = Config::config_file_path();
        let mut stream = std::fs::File::create(&file).expect("should create config file");
        stream.write_all(Config::default_config().as_bytes()).expect("should write to config file");
    }

    /// AIO function to ensure all the config files and directories are created
    fn ensure_created_config() {

        // Ensure directories exist
        let dir = Config::config_dir_path().to_str().expect("should convert to string slice").to_string();
        if !fs::metadata(&dir).is_ok() {
            std::fs::create_dir_all(Config::config_dir_path()).expect("should recursively create full config dir path");
        }

        // Ensure config file exists
        let file = Config::config_file_path();
        if !fs::metadata(&file).is_ok() {
            Config::create_config_file();
        }
    }

    /// Loads the config file from the config directory
    pub fn load() -> Config {
        // First we make sure that a config file is created
        Config::ensure_created_config();

        // Open and read config file contents
        let mut config_file = std::fs::File::open(Config::config_file_path()).expect("should open config file");
        let mut contents = String::new();
        config_file.read_to_string(&mut contents).expect("should read config file");
        
        // Serialize the contents to a config struct and return it
        toml::from_str(&contents).expect("should parse config file")
    }
}



#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn config_file_path_should_return_non_empty_string() {
        let path = Config::config_file_path();
        assert!(!path.is_empty());
    }
    
    #[test]

    fn config_dir_path_should_return_non_empty_path() {
        let path = Config::config_dir_path();
        assert!(!path.to_str().expect("should return non-empty string").is_empty());
    }
}
