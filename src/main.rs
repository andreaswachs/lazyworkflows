// #![deny(warnings)]

mod config;
mod workflows;

#[tokio::main]
async fn main() -> Result<(), isize> {
    println!("main()");
    // Load the config file
    let config = config::Config::load();

    // Load the workflows
    workflows::load(&config).await;

    println!("hello world from main");
    Ok(())
}
