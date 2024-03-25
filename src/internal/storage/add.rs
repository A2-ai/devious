use std::path::PathBuf;

use crate::internal::file::hash;

pub fn add(local_path: &PathBuf, storage_dir: &PathBuf, git_dir: &PathBuf, message: String) -> Result<String, std::io::Error> {
    // get file hash
    let file_hash = match hash::get_file_hash(local_path){
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
    let storage_path = hash::get_storage_path(&storage_dir, &file_hash);

}