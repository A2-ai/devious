use crate::internal::git::repo;
use crate::internal::storage::init;
use std::path::PathBuf;

pub fn get_init_runner(storage_dir: &PathBuf, mode: &u32, gid: &u32) -> Result<(), std::io::Error> {
    
    // Get git root
    let git_repo = match repo::get_nearest_repo_dir(&storage_dir) {
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
    let init_run = match init::init(&git_repo, &storage_dir, &mode, &gid) {
        Ok(init_run) => {
            // json
            init_run
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    Ok(())
}