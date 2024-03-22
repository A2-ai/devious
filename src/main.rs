use std::path::PathBuf;
use std::env;
use clap::{arg, Command, ArgAction, Arg};
use std::io;

mod internal;
use crate::internal::storage::init;
use crate::internal::config::config;
use crate::internal::git::repo;
mod cmd;
use crate::cmd::init_cmd;

fn main() -> io::Result<()> {
    let intro = Command::new("🌀 Devious")
        .version("0.0.1")
        .author("Andriy Massimilla, Jenna Johnson")
        .about("version large files under Git")
        .arg_required_else_help(false)
        .get_matches();
        
       




    let current = env::current_dir().unwrap().display().to_string();
    println!("{current}");
    let nearest_repo = repo::get_nearest_repo_dir(&PathBuf::from(r"src")).unwrap().display().to_string();
    println!("NEAREST PATH: {nearest_repo}");

    let dir = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src");
    let config = config::Config{storage_dir: dir.clone()};
    
    let write_result = config::write(&config, &dir);
    match write_result {
        Ok(_) => {}
        Err(_) => {}
    }

    let read_result = config::read(&dir);
    match read_result {
        Ok(_) => {}
        Err(_) => {}
        
    }

    
    


    
}
