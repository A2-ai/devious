use crate::internal::config::config;
use crate::internal::git::repo;
use crate::internal::utils::utils::normalize_path;
// use std::os::unix::fs::chown;
use std::path::PathBuf;
use std::fs::{self, create_dir};
use std::os::unix::fs::PermissionsExt;
use std::fs::Permissions;
use file_owner::{Group, PathExt};

pub fn init(root_dir: &PathBuf, storage_dir: &PathBuf, mode: &u32, group_name: &String) -> Result<(), std::io::Error> { // 
    // get storage directory as a normalized path
    let storage_dir_norm = normalize_path(storage_dir);

    // should I use .exists or .is_dir?
    if !storage_dir.exists() { // if storage directory doesn't exist
        // create storage dir
        match create_dir(&storage_dir_norm) {
            Ok(_) => {
                // json: success
            }
            Err(e) => {
                // json: fail
                return Err(e)
            }
        };
    } // if

    else { // else, storage directory exists
        // json: Storage directory already exists

        //  Warn if storage dir is empty
        match repo::is_directory_empty(storage_dir_norm.as_path()) {
            Ok(empty) => {
                if empty {
                    // json
                println!("directory is empty")
                }
                else {
                    println!("directory is not empty")
                }
            }
            Err(e) => {
                //json
                return Err(e)
            }
        }
    } // else

    // get absolute path
    let storage_dir_abs = match storage_dir_norm.canonicalize() {
        Ok(path) => path,
        Err(e) => return Err(e),
    };

    // set permissions
    let permissions: Permissions = fs::Permissions::from_mode(*mode);
    match fs::set_permissions(&storage_dir_abs, permissions) {
        Ok(_) => {
            // json: success
        },
        Err(e) => {
            // json: fail
            return Err(e)
        }
    };

    // set group ownership
    let group = match Group::from_name(group_name) {
        Ok(group) => {
            // json
            group
        }
        Err(_) => return Err(std::io::Error::other("group name is invalid"))
    };
    match storage_dir_abs.set_group(group) {
        Ok(_) => {},
        Err(_) => return Err(std::io::Error::other("group name is invalid"))
    };

    // write config
    match config::write(&config::Config{storage_dir: storage_dir_abs.clone(), permissions: mode.clone(), group: group_name.clone()}, &root_dir) {
        Ok(_) => {
            // json
        }
        Err(e) => {
            //json
            return Err(e);
        }
    };

    Ok(())
    // json: success
}