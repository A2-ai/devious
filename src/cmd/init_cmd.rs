use crate::internal::git::repo;
use crate::internal::storage::init;
use std::path::PathBuf;

fn get_init_runner(storage_dir: &PathBuf, mode: u32, gid: u32) -> Result<(), std::io::Error> {
    
    let git_repo = match repo::get_nearest_repo_dir(storage_dir) {
        Ok(git_repo) => {
            // json
            git_repo
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    let init_run = init::init(&git_repo, &storage_dir, &mode, &gid);
    


}