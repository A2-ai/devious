use std::path::PathBuf;

use crate::internal::config::config;
use crate::internal::git::repo;
use crate::internal::meta::parse;
use crate::internal::storage::get;

pub fn run_get_cmd(globs: &Vec<String>) -> Result<(), std::io::Error> {
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

    // parse each glob
    let queued_paths = parse::parse_globs(globs);

    // Get the queued files
    for path in &queued_paths {
        match get::get(&path, &conf.storage_dir, &git_dir) {
            Ok(_) => {}
            Err(e) => {
                // json
                return Err(e);
            }
        };
    }

    // warn if no files were queued
    if queued_paths.is_empty() {
        // json: warning
    }

    Ok(())
}