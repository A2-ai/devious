use std::path::PathBuf;
use crate::internal::git::repo;
use crate::internal::storage::add;
use crate::internal::config::config;
use anyhow::{Context, Result};

pub fn run_add_cmd(files: &Vec<String>, message: &String) -> Result<()> {
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
        add::add(file, &conf, &message)?;
    }

    if queued_paths.is_empty() {
        // json warning: no files were queued
    }
   
    Ok(())
} // run_add_cmd