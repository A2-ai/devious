use crate::internal::config::config;
use crate::internal::git::repo;
use std::path::PathBuf;
use std::fs::create_dir;
use path_absolutize::Absolutize;
use file_owner::Group;
use anyhow::{Context, Result};
use anyhow::anyhow;


pub fn init(root_dir: &PathBuf, storage_dir: &PathBuf, mode: &u32, group_name: &String) -> Result<()> { 
    // get absolute path, but don't check if it exists yet
    let storage_dir_abs = PathBuf::from(storage_dir.absolutize().unwrap());

    
    // check if directory exists
    if !storage_dir_abs.exists() { // if storage directory doesn't exist
        println!("storage directory doesn't exist. Creating storage directory...");
        // create storage dir
        create_dir(&storage_dir_abs).with_context(|| format!("Failed to create storage directory: {}", storage_dir.display()))?;
    } // if

    else { // else, storage directory exists
        println!("storage directory already exists");

        //  Warn if storage dir is not empty
        match repo::is_directory_empty(&storage_dir_abs) {
            Ok(empty) => {
                if !empty {
                    println!("warning: storage directory is not empty")
                }
            }
            Err(e) => {
                //json
                return Err(anyhow!("unable to check if directory is empty: {}", e))
            }
        }
    } // else

    // check if group_name refers to an actual group
    Group::from_name(group_name).with_context(|| "group not found: {group_name}")?;
   
    // write config
    config::write(
        &config::Config{storage_dir: storage_dir_abs.clone(), 
            permissions: mode.clone(), 
            group: group_name.clone()}, 
            &root_dir)
            .with_context(|| "unable to write configuration file")?;
       
    Ok(())
    // json: success
}