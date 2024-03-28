use std::os::unix::fs::PermissionsExt;
use std::path::PathBuf;
use chrono::Utc;
use file_owner::{Group, PathExt};
use std::fs::{self, Permissions};

use crate::internal::config::config::Config;
use crate::internal::file::hash;
use crate::internal::storage::copy;
use crate::internal::meta::file;
use crate::internal::git::ignore;


pub fn add(local_path: &PathBuf, conf: &Config, message: &String) -> Result<String, std::io::Error> {
    // get file hash
    let file_hash = match hash::hash_file_with_blake3(local_path) {
        Ok(file_hash) => {
            //json
            file_hash
        }
        Err(e) => {
            // json
            return Err(e);
        }
    };

    // get storage path
    let storage_dir_abs = match conf.storage_dir.canonicalize() {
        Ok(path) => path,
        Err(e) => return Err(e),
    };
    let dest_path = hash::get_storage_path(&storage_dir_abs, &file_hash);

    // Copy the file to the storage directory
	// if the destination already exists, skip copying
    if !dest_path.exists() {
        // copy
        match copy::copy(&local_path, &dest_path) {
            Ok(_) => {
                
                // json
             }
             Err(e) => {
                //json
                return Err(e)
             }
        };
        //json
    }
    else {
        // json
    }

    // set permissions
    let mode = conf.permissions;
    let permissions: Permissions = fs::Permissions::from_mode(mode);
    match fs::set_permissions(&dest_path, permissions) {
        Ok(_) => {
            // json: success
        },
        Err(e) => {
            // json: fail
            return Err(e)
        }
    };

    // set group ownership
    let group_name = conf.group.as_str();
    let group = match Group::from_name(group_name) {
        Ok(group) => {
            // json
            group
        }
        Err(_) => return Err(std::io::Error::other("group name is invalid"))
    };
    match dest_path.set_group(group) {
        Ok(_) => {},
        Err(_) => return Err(std::io::Error::other("group name is invalid"))
    };

    // get file size
    let local_path_data = match local_path.metadata() {
        Ok(data) => data,
        Err(e) => return Err(e),
    };
    let file_size = local_path_data.len();

    // get user name
    let owner = match local_path.owner() {
        Ok(owner) => owner,
        Err(_) => return Err(std::io::Error::other("file owner not found")),
    };
    let owner_name = match owner.name() {
        Ok(name) => name.unwrap(),
        Err(_) => {return Err(std::io::Error::other("file owner not found"))},
    };

    // create + write metadata file
    let metadata = file::Metadata{
        file_hash: file_hash.clone(),
        file_size,
        time_stamp: chrono::offset::Local::now().to_string(),
        message: message.clone(),
        saved_by: owner_name
    };
    match file::save(&metadata, &local_path) {
        Ok(_) => {},
        Err(_) => return Err(std::io::Error::other("metadata file not created"))
    };

    // Add file to gitignore
    match ignore::add_gitignore_entry(local_path) {
        Ok(_) => {},
        Err(_) => return Err(std::io::Error::other("gitignore entry could not be created"))
    };
    
    return Ok(file_hash);
}