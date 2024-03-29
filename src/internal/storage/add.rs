use std::os::unix::fs::PermissionsExt;
use std::path::PathBuf;
use file_owner::{Group, PathExt};
use std::fs;
use anyhow::{anyhow, Context, Result};
use crate::internal::config::config::Config;
use crate::internal::file::hash;
use crate::internal::storage::copy;
use crate::internal::meta::file;
use crate::internal::git::ignore;


pub fn add(local_path: &PathBuf, conf: &Config, message: &String) -> Result<String> {
    // get file hash
    let file_hash = hash::hash_file_with_blake3(local_path).with_context(|| format!("could not hash file"))?;

    // get storage path
    let storage_dir_abs = conf.storage_dir.canonicalize().with_context(|| format!("could not find storage directory: {}", conf.storage_dir.display()))?;
    let dest_path = hash::get_storage_path(&storage_dir_abs, &file_hash);

    // check if group exists again
    let group = Group::from_name(&conf.group).with_context(|| format!("group not found: {}", conf.group))?;

    // Copy the file to the storage directory
	// if the destination already exists, skip copying
    if !dest_path.exists() {
        // copy
        copy::copy(&local_path, &dest_path).with_context(|| format!("could not copy {} to storage directory: {}", local_path.display(), dest_path.display()))?;
    }
    else {
        println!("{} already exists in storage directory", local_path.display())
    }

    // set permissions
    let mode = conf.permissions;
    fs::set_permissions(&dest_path, fs::Permissions::from_mode(mode)).with_context(|| format!("unable to set permissions: {}", mode))?;

    // set group ownership
    dest_path.set_group(group).with_context(|| format!("unable to set group: {}", group))?;

    // get file size
    let local_path_data = local_path.metadata().with_context(|| format!("unable to get size of file: {}", local_path.display()))?;
    let file_size = local_path_data.len();

    // get user name
    let owner = local_path.owner().with_context(|| format!(""))?;
    let owner_name = match owner.name() {
        Ok(name) => name.unwrap(),
        Err(e) => {return Err(anyhow!("could not get name of file owner: {e}"))},
    };

    // create + write metadata file
    let metadata = file::Metadata{
        file_hash: file_hash.clone(),
        file_size,
        time_stamp: chrono::offset::Local::now().to_string(),
        message: message.clone(),
        saved_by: owner_name
    };
    file::save(&metadata, &local_path).with_context(|| format!("could not save metadata into file"))?;

    // Add file to gitignore
    ignore::add_gitignore_entry(local_path).with_context(|| format!("could not add .gitignore entry"))?;
    
    return Ok(file_hash);
}