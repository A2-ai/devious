use std::path::PathBuf;
use anyhow::{Context, Result};
use crate::internal::file::hash;
use crate::internal::meta::file;
use crate::internal::storage::copy;

// gets a file from storage
pub fn get(local_path: &PathBuf, storage_dir: &PathBuf) -> Result<()> {
    // get metadata
    let metadata = file::load(&local_path).with_context(|| format!("could not get metadata file for {}", local_path.display()))?;

    // get storage data
    let storage_path = hash::get_storage_path(storage_dir, &metadata.file_hash);

    // get hashes to compare
    let local_hash = hash::get_file_hash(&local_path);
    let metadata_hash = metadata.file_hash;

    // check if up-to-date file is already present locally
    if !local_path.exists() || local_hash == String::from("") || metadata_hash == String::from("") || local_hash != metadata_hash {
        copy::copy(&storage_path, &local_path).with_context(|| format!("could not copy {} from storage directory: {}", local_path.display(), storage_dir.display()))?;
    }
    else {
        println!("{} already up to date", local_path.display())
    }

    Ok(())
} // get