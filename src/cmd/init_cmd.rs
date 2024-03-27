use crate::internal::git::repo;
use crate::internal::storage::init;
use std::path::PathBuf;

pub fn run_init_cmd(storage_dir: &PathBuf, mode: &u32, group: &String) -> Result<(), std::io::Error> {
    
    // Get git root
    let git_repo = match repo::get_nearest_repo_dir(&PathBuf::from(".")) {
        Ok(git_repo) => {
            // json
            git_repo
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    // Initialize
    match init::init(&git_repo, &storage_dir, &mode, &group) {
        Ok(_) => {
            // json
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    Ok(())
}