use std::path::PathBuf;
use std::time::SystemTime;
use file_owner::PathExt;

use crate::internal::file::hash;
use crate::internal::storage::copy;
use crate::internal::meta::file;
use crate::internal::git::ignore;

pub fn add(local_path: &PathBuf, storage_dir: &PathBuf, git_dir: &PathBuf, message: String) -> Result<String, std::io::Error> {
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
    let dest_path = hash::get_storage_path(&storage_dir, &file_hash);

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

    // get file size
    let local_path_data = match local_path.metadata() {
        Ok(data) => data,
        Err(e) => return Err(e),
    };
    let file_size = local_path_data.len();

    // get user
    let owner = match local_path.owner() {
        Ok(owner) => owner,
        Err(_) => return Err(std::io::Error::other("file owner not found")),
    };
    // get user name
    let owner_name = match owner.name() {
        Ok(name) => name.unwrap(),
        Err(_) => {return Err(std::io::Error::other("file owner not found"))},
    };

    // get group 
    let group = match local_path.group() {
        Ok(group) => group,
        Err(_) => return Err(std::io::Error::other("file group not found")),
    };
    // get group name
    let group_name = match group.name() {
        Ok(name) => name.unwrap(),
        Err(_) => {return Err(std::io::Error::other("file group not found"))},
    };

    // create + write metadata file
    let metadata = file::Metadata{
        file_hash: file_hash.clone(),
        file_size: file_size,
        time_stamp: SystemTime::now(),
        message: message,
        group: group_name,
        saved_by: owner_name
    };

    match file::save(&metadata, &local_path) {
        Ok(_) => {},
        Err(_) => return Err(std::io::Error::other("metadata file not created"))
    };

    // Add file to gitignore
    ignore::add_gitignore_entry(git_dir, local_path).expect("gitignore entry unable to be added");
    
    return Ok(file_hash);
}