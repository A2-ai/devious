use std::path::PathBuf;
use config::{read, write};

mod config;
mod init;
mod repo;

fn main() {
    let dir = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src");
    let config = config::Config{storage_dir: dir.clone()};
    
    let write_result = write(&config, &dir);
    match write_result {
        Ok(_) => {}
        Err(_) => {}
    }

    let read_result = read(&dir);
    match read_result {
        Ok(_) => {}
        Err(_) => {}
    }

    
    


    
}
