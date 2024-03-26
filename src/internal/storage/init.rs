use crate::config::{self, Config};
use crate::internal::git::repo;
use crate::internal::utils::utils::normalize_path;
// use std::os::unix::fs::chown;
use std::path::PathBuf;
use std::fs::{self, create_dir};
use std::os::unix::fs::PermissionsExt;
use std::fs::Permissions;
use file_owner::{Group, PathExt};

pub fn init(root_dir: &PathBuf, storage_dir: &PathBuf, mode: &u32, group_name: &String) -> Result<(), std::io::Error> { // 
    // get storage directory as an absolute path (don't use canoniicalize which checks if the path exists)
    let storage_dir_abs = normalize_path(storage_dir);
    
    // check storage directory permissions
    let md = fs::metadata(&storage_dir_abs)?;
    let permissions = md.permissions();
    let readonly = permissions.readonly();
    if readonly { // create storage dir
        // create storage dir
        match create_dir(&storage_dir_abs) {
            Ok(_) => {
                // json: success
            }
            Err(e) => {
                // json: fail
                return Err(e)
            }
        };
    }
    else { // storage dir already exists
        // json: Storage directory already exists

        // Ensure storage dir is a directory
        if !storage_dir_abs.clone().is_dir() {
            // json: fail
            return Err(std::io::Error::other("storage dir is not a directory"));
        }

        //  Warn if storage dir is not readable
        match repo::is_directory_empty(storage_dir_abs.as_path()) {
            Ok(_) => {
                // json
                println!("directory is empty")
            }
            Err(e) => {
                //json
                return Err(e)
            }
        }
    } // else

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
    match config::write(&Config{storage_dir: storage_dir_abs.clone()}, &root_dir) {
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