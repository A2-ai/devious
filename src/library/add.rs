use std::path::PathBuf;
use std::os::unix::fs::PermissionsExt;
use serde::Serialize;
use file_owner::{Group, PathExt};
use std::{fs, u32};
use anyhow::{Context, Result};
use crate::helpers::hash;
use crate::helpers::copy;
use crate::helpers::file;
use crate::helpers::ignore;
use crate::helpers::config;
use crate::helpers::repo;

#[derive(Clone, PartialEq, Serialize)]
enum Outcome {
    Success,
    AlreadyPresent,
    Error
}

#[derive(Clone, PartialEq, Serialize)]
pub struct AddedFile {
    path: PathBuf,
    hash: Option<String>,
    outcome: Outcome,
    error: Option<String>,
    size: Option<u64>,
}

pub fn dvs_add(files: &Vec<String>, message: &String) -> Result<Vec<AddedFile>> {
    // Get git root
    let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

    // load the config
    let conf = config::read(&git_dir).with_context(|| "could not load configuration file - no dvs.yaml in directory - be sure to initiate devious")?;

    let mut queued_paths: Vec<PathBuf> = Vec::new();

    for file_in in files {
        let file_without_meta = file_in.replace(".dvsmeta", "");
        let file = PathBuf::from(file_without_meta);

        if queued_paths.contains(&file) {continue}

        queued_paths.push(file);
    } // for
    
    // add each file in queued_paths to storage
    let added_files = queued_paths.into_iter().map(|file| {
        add(&file, &git_dir, &conf, &message)
    }).collect::<Vec<AddedFile>>();

    return Ok(added_files)
} // run_add_cmd

fn add(local_path: &PathBuf, git_dir: &PathBuf, conf: &config::Config, message: &String) -> AddedFile {
    // set error to None by default
    let mut error: Option<String> = None;

    if error.is_none() {error = get_preliminary_errors(local_path, git_dir, conf, message)}

    // get file hash
    let file_hash = hash::get_file_hash(&local_path);
    if file_hash.is_none() && error.is_none() {
        error = Some(String::from("hash not found"));
    }

    // get file size
    let file_size = get_file_size(&local_path);
    if file_size.is_none() && error.is_none() {
        error = Some(String::from("file size not found"));
    }

    // get user name
    let user_name = get_user_name(&local_path);
    if user_name.is_none() && error.is_none() {
        error = Some(String::from("file owner not found"));
    }

    // check group
    let group: Option<Group> = match Group::from_name(&conf.group) {
        Ok(group) => Some(group),
        Err(_) => {
            if error.is_none() {error = Some(String::from("group not found"))}
            None
        }
    };

    // now see if file can be added
    let storage_dir_abs: Option<PathBuf> = match conf.storage_dir.canonicalize() {
        Ok(path) => Some(path),
        Err(_) => {
            if error.is_none() {error = Some(String::from("storage directory not found"))}
            None
        }
    };


    

    if error.is_some() {
        return AddedFile{
            path: local_path.clone(), 
            hash: file_hash,
            outcome: Outcome::Error,
            error: error,
            size: file_size
        };
    }

    // can safely unwrap storage_dir_abs and file_hash 
    let storage_dir_abs_value = storage_dir_abs.unwrap();
    let file_hash_value = file_hash.clone().unwrap();
    
    // get storage path
    let dest_path = hash::get_storage_path(&storage_dir_abs_value, &file_hash_value);

    // Copy the file to the storage directory if it's not already there
    let mut outcome: Outcome = Outcome::Success;
    if !dest_path.exists() {
        // copy
        copy_file_to_storage_directory(local_path, &dest_path, &conf.permissions, &group.unwrap());
    }
    else {
        outcome = Outcome::AlreadyPresent;
    }

    // create metadata
    let metadata = file::Metadata{
        file_hash: file_hash_value,
        file_size: file_size.unwrap(),
        time_stamp: chrono::offset::Local::now().to_string(),
        message: message.clone(),
        saved_by: user_name.unwrap()
    };

    // write metadata file
    match file::save(&metadata, &local_path) {
        Ok(_) => {},
        Err(_) => if error.is_none() {error = Some(String::from("could not save metadata file"))}
    };

    // Add file to gitignore
    match ignore::add_gitignore_entry(local_path) {
        Ok(_) => {},
        Err(_) => {
            if error.is_none() {error = Some(String::from("could not add .gitignore entry"))}
        }
    };
    
    if error.is_some() {outcome = Outcome::Error}

    return AddedFile {
        path: local_path.clone(),
        hash: file_hash.clone(),
        outcome,
        error,
        size: file_size
    }
}

fn get_preliminary_errors(local_path: &PathBuf, git_dir: &PathBuf, conf: &config::Config, message: &String) -> Option<String> {
    // check if file exists
    match local_path.canonicalize() {
        Ok(local_path) => { // file exists
            // if file is outside of git repository
            if local_path.strip_prefix(&git_dir).unwrap() == local_path {
                return Some(String::from("file not in git repository"));
            }
        }
        Err(_) => { 
            return Some(String::from("file not found"))
        }
    };

    if local_path.is_dir() {
        return Some(String::from("path is a directory"))
    }

    
    None
}


fn get_file_size(local_path: &PathBuf) -> Option<u64> {
    match local_path.metadata() {
        Ok(data) => return Some(data.len()),
        Err(_) => return None,
    };
}


fn get_user_name(local_path: &PathBuf) -> Option<String> {
    let owner = match local_path.owner().with_context(|| format!("")) {
        Ok(owner) => owner,
        Err(_) => return None,
    };
    match owner.name() {
        Ok(name) => return Some(name.unwrap()),
        Err(_) => return None,
    };
}


fn copy_file_to_storage_directory(local_path: &PathBuf, dest_path: &PathBuf, mode: &u32, group: &Group) -> Option<String> {
    let mut error = None;
    match copy::copy(&local_path, &dest_path) {
        Ok(_) => {
            // set permissions
            match set_permissions(&mode, &dest_path) {
                Ok(_) => {},
                Err(_) => {
                    // set error
                    if error.is_none() {error = Some(String::from("could not set file permissions"))}
                    // delete copied file
                    fs::remove_file(&dest_path)
                    .expect(format!("could not set permissions after copying {} to {}: error deleting copied file. Delete {} manually.", local_path.display(), dest_path.display(), dest_path.display()).as_str());
                }
            };

            // set group ownership
            match dest_path.set_group(group.clone()) {
                Ok(_) => {},
                Err(_) => {
                    // set error
                    if error.is_none() {error = Some(String::from("could not set file group ownership"))}
                    // delete copied file
                    fs::remove_file(&dest_path)
                    .expect(format!("could not set group ownership after copying {} to {}: error deleting copied file. Delete {} manually.", local_path.display(), dest_path.display(), dest_path.display()).as_str());

                }
            };
        } // Ok, could copy
        Err(_) => {
            if error.is_none() {error = Some(String::from("could not copy file to storage directory"))}
        }
    };
    return error
}

fn set_permissions(mode: &u32, dest_path: &PathBuf) -> Result<()> {
    dest_path.metadata().unwrap().permissions().set_mode(*mode);
    let _file_mode = dest_path.metadata().unwrap().permissions().mode();
    let new_permissions = fs::Permissions::from_mode(*mode);
    fs::set_permissions(&dest_path, new_permissions).with_context(|| format!("unable to set permissions: {}", mode))?;
    Ok(())
}