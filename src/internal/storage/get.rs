use std::path::PathBuf;

use crate::internal::file::hash;
use crate::internal::meta::file;
use crate::internal::storage::copy;

// gets a file from storage
pub fn get(local_path: &PathBuf, storage_dir: &PathBuf, git_dir: &PathBuf) -> Result<(), std::io::Error> {
    // get metadata
    let metadata = match file::load(&local_path) {
        Ok(metadata) => metadata,
        Err(_) => return Err(std::io::Error::other("failed to load metadata")),
    };

    let storage_path = hash::get_storage_path(storage_dir, &metadata.file_hash);

    // check if file is already present locally
    let local_exists = local_path.exists();
    let local_hash = hash::get_file_hash(&local_path)?;
    let metadata_hash = metadata.file_hash;
    let empty_string = String::from("");

    if !local_exists || local_hash == empty_string || local_hash == empty_string || local_hash != metadata_hash {
        match copy::copy(&storage_path, &local_path) {
            Ok(_) => {}
            Err(e) => {
                // json: failed to copy the file
                return Err(e)
            }
        };
    }
    else {
        // json: file is already up to date
    }

    Ok(())
} // get