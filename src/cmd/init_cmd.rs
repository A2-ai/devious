use crate::internal::git::repo;
use crate::internal::storage::init;
use std::path::PathBuf;
use anyhow::{Context, Result};

pub fn run_init_cmd(storage_dir: &PathBuf, mode: &u32, group: &String) -> Result<()> {
    // Get git root
   let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

    // Initialize
    init::init(&git_dir, &storage_dir, &mode, &group)?;

    Ok(())
}