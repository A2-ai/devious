use std::path::PathBuf;
use crate::internal::git::repo;
use crate::internal::storage::add;
use crate::config;

pub fn run_add_cmd(files: &Vec<PathBuf>, message: &String) -> Result<(), std::io::Error> {
   // Get git root
   let git_dir = match repo::get_nearest_repo_dir(&PathBuf::from(".")) {
        Ok(git_repo) => {
            // json
            git_repo
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    // load the config
    let conf = match config::read(&git_dir) {
        Ok(config) => config,
        Err(_) => return Err(std::io::Error::other("config not readable")),
    };

    let mut queued_paths: Vec<PathBuf> = Vec::new();

    for file_in in files {
        let string = file_in.display().to_string();
        let string_without_meta = string.replace(".dvsmeta", "");
        let file = PathBuf::from(string_without_meta);

        if queued_paths.contains(&file) {continue}

        // ensure file is inside of the git repo
        let abs_path = match file.canonicalize() {
            Ok(file) => file,
            Err(_) => {
                // json warning: skipping invalid path
                continue;
            }
        };
        if abs_path.strip_prefix(&git_dir).unwrap() == abs_path {
            // json warning: skipped file outside of git repository
            continue;
        }

        // skip directories
        if file.is_dir() {
            // json warning: skipped directory
            continue
        }

        // all checks passed, finally add file to queued_paths
        queued_paths.push(file);
    } // for
    
    
    // add each file in queued_paths to storage
    for file in &queued_paths {
        match add::add(file, &conf, &git_dir, &message) {
            Ok(_) => {}
            Err(e) => return Err(e),
        };
    }

    if queued_paths.is_empty() {
        // json warning: no files were queued
    }
   
    Ok(())
} // run_add_cmd