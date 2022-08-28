#![deny(warnings)]

mod config;

fn main() {
    // Load the config file
    let _config = config::Config::load();
}
