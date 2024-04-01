use std::path::PathBuf;
use std::os::unix::fs::PermissionsExt;
use file_owner::{Group, PathExt};
use std::{fs, u32};
use anyhow::{anyhow, Context, Result};
use crate::helpers::hash;
use crate::helpers::copy;
use crate::helpers::file;
use crate::helpers::ignore;
use crate::helpers::config;
use crate::helpers::repo;

pub fn dvs_add(files: &Vec<String>, message: &String) -> Result<()> {
   // Get git root
   let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

    // load the config
    let conf = config::read(&git_dir)?;

    let mut queued_paths: Vec<PathBuf> = Vec::new();

    for file_in in files {
        let file_without_meta = file_in.replace(".dvsmeta", "");
        let file = PathBuf::from(file_without_meta);

        if queued_paths.contains(&file) {continue}

        // ensure file is inside of the git repo
        let abs_path = match file.canonicalize() {
            Ok(file) => file,
            Err(_) => { // swallowing error here because the command can still run
                println!("skipping {} - doesn't exist", file.display());
                continue;
            }
        };
        if abs_path.strip_prefix(&git_dir).unwrap() == abs_path {
            println!("skipping {} - outside of git repository", file.display());
            continue;
        }

        // skip directories
        if file.is_dir() {
            println!("skipping {} - is a directory", file.display());
            continue
        }

        // all checks passed, finally add file to queued_paths
        queued_paths.push(file);
    } // for
    
    
    // add each file in queued_paths to storage
    for file in &queued_paths {
        add(file, &conf, &message)?;
    }

    if queued_paths.is_empty() {
        // json warning: no files were queued
    }
   
    Ok(())
} // run_add_cmd

fn add(local_path: &PathBuf, conf: &config::Config, message: &String) -> Result<String> {
    // get file hash
    let file_hash = hash::hash_file_with_blake3(local_path).with_context(|| format!("could not hash file"))?;

    // get storage path
    let storage_dir_abs = conf.storage_dir.canonicalize().with_context(|| format!("could not find storage directory: {}", conf.storage_dir.display()))?;
    let dest_path = hash::get_storage_path(&storage_dir_abs, &file_hash);

    // check if group exists again
    let group = Group::from_name(&conf.group).with_context(|| format!("group not found: {}", conf.group))?;

    // Copy the file to the storage directory if it's not already there
    if !dest_path.exists() {
        // copy
        copy::copy(&local_path, &dest_path).with_context(|| format!("could not copy {} to storage directory: {}", local_path.display(), dest_path.display()))?;
    }
    else {
        println!("{} already exists in storage directory", local_path.display())
    }

    // set permissions
    set_permissions(&conf.permissions, &dest_path)?;

    // set group ownership
    dest_path.set_group(group).with_context(|| format!("unable to set group: {}", group))?;

    // get file size
    let file_size = get_file_size(&local_path)?;

    // get user name
    let user_name = get_user_name(&local_path)?;

    // create + write metadata file
    let metadata = file::Metadata{
        file_hash: file_hash.clone(),
        file_size,
        time_stamp: chrono::offset::Local::now().to_string(),
        message: message.clone(),
        saved_by: user_name
    };
    file::save(&metadata, &local_path).with_context(|| format!("could not save metadata into file"))?;

    // Add file to gitignore
    ignore::add_gitignore_entry(local_path).with_context(|| format!("could not add .gitignore entry"))?;
    
    return Ok(file_hash);
}

fn set_permissions(mode: &u32, dest_path: &PathBuf) -> Result<()> {
    dest_path.metadata().unwrap().permissions().set_mode(*mode);
    let _file_mode = dest_path.metadata().unwrap().permissions().mode();
    let new_permissions = fs::Permissions::from_mode(*mode);
    fs::set_permissions(&dest_path, new_permissions).with_context(|| format!("unable to set permissions: {}", mode))?;
    Ok(())
}

fn get_file_size(local_path: &PathBuf) -> Result<u64> {
    let local_path_data = local_path.metadata().with_context(|| format!("unable to get size of file: {}", local_path.display()))?;
    return Ok(local_path_data.len());
}

fn get_user_name(local_path: &PathBuf) -> Result<String> {
    let owner = local_path.owner().with_context(|| format!(""))?;
    match owner.name() {
        Ok(name) => return Ok(name.unwrap()),
        Err(e) => return Err(anyhow!("could not get name of file owner: {e}")),
    };
}