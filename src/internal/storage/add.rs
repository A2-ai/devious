use std::path::PathBuf;

use crate::internal::file::hash;
use crate::internal::storage::copy;

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
    let local_path_size = local_path_data.len();

    // get user


    // Create + write metadata file


    // Add file to gitignore


    return Ok(file_hash);

}