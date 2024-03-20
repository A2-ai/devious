use yaml_rust::{YamlLoader, YamlEmitter};
use std::{os, path};
use std::env;
use std::fs;
use std::path::PathBuf;


pub struct Config {
    storage_directory: String
}

impl Default for Config {
    fn default() -> Self {
        Config {storage_directory: String::from("yaml\":storage-dir")}
    }
  }

//static CONFIG_FILE_NAME: PathBuf = PathBuf::from(r"dvs.yaml");

pub fn read(root_dir: PathBuf) {
    let config_file_contents = fs::read_to_string(root_dir.join(PathBuf::from(r"dvs.yaml"))).unwrap();
    
    // match config_file_contents {
    //     Ok(contents) => {
            
    //     }
} // read






