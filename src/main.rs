use std::path::PathBuf;

use config::{read, write};


mod config;

fn main() {
    let dir = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src");
    let config = config::Config{storage_dir: dir.clone()};
    
    write(&config, &dir);

    read(&dir);
}
