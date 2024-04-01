use std::path::PathBuf;
use anyhow::{Context, Result};
use crate::helpers::hash;
use crate::helpers::copy;
use crate::helpers::file;
use crate::helpers::config;
use crate::helpers::repo;
use crate::helpers::parse;

pub fn dvs_get(globs: &Vec<String>) -> Result<()> {
    // Get git root
   let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

    // load the config
    let conf = config::read(&git_dir).with_context(|| "dvs.yaml is not present in your directory - have you initialized devious?")?;

    // parse each glob
    let queued_paths = parse::parse_globs(globs);

    // Get the queued files
    for path in &queued_paths {
        get(&path, &conf.storage_dir).with_context(|| format!("could not retrieve {} from storage directory", path.display()))?;
    }

    if queued_paths.is_empty() {
       println!("warning: no files were queued")
    }

    Ok(())
}

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
    if !local_path.exists() || local_hash.is_none() || metadata_hash == String::from("") || local_hash.unwrap() != metadata_hash {
        copy::copy(&storage_path, &local_path).with_context(|| format!("could not copy {} from storage directory: {}", local_path.display(), storage_dir.display()))?;
    }
    else {
        println!("{} already up to date", local_path.display())
    }

    Ok(())
} // get