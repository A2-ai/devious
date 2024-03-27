use extendr_api::prelude::*;
use extendr_api::robj::{Robj, IntoRobj};
use std::path::PathBuf;

mod internal;
use crate::internal::git::repo;
use crate::internal::storage::init;

// initiate storage directory
// @export
#[extendr]
pub fn run_init_cmd_R(storage_dir: &PathBuf, mode: &u32, group: &String) -> Robj {
    
    // Get git root
    let git_repo = match repo::get_nearest_repo_dir(&PathBuf::from(".")) {
        Ok(git_repo) => {
            // json
            git_repo
        }
        Err(e) => {
            // json
            return Robj::from(format!("Error getting git repo: {}", e));
        }
    };

    // Initialize
    match init::init(&git_repo, &storage_dir, &mode, &group) {
        Ok(_) => {
            // json
        }
        Err(e) => {
            // json
            return Robj::from(format!("Error initializing: {}", e));
        }
    };
    
    let val: Result<()> = Ok(());
    return val.into_robj();
}


// Macro to generate exports.
// This ensures exported functions are registered with R.
// See corresponding C code in `entrypoint.c`.
extendr_module! {
    mod devious;
    fn run_init_cmd_R;
}