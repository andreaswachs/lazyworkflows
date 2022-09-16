// #![deny(warnings)]

mod config;
mod workflows;

#[tokio::main]
async fn main() -> Result<(), isize> {
    // Load the config file
    let config = config::manager::Config::load();

    // Load the workflows
    workflows::load(&config).await;

    Ok(())
}
