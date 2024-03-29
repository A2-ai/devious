use std::path::PathBuf;
use anyhow::{Context, Result};
use crate::internal::config::config;
use crate::internal::git::repo;
use crate::internal::meta::parse;
use crate::internal::storage::get;

pub fn run_get_cmd(globs: &Vec<String>) -> Result<()> {
    // Get git root
   let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

    // load the config
    let conf = config::read(&git_dir).with_context(|| "dvs.yaml is not present in your directory - have you initialized devious?")?;

    // parse each glob
    let queued_paths = parse::parse_globs(globs);

    // Get the queued files
    for path in &queued_paths {
        get::get(&path, &conf.storage_dir).with_context(|| format!("could not retrieve {} from storage directory", path.display()))?;
    }

    if queued_paths.is_empty() {
       println!("warning: no files were queued")
    }

    Ok(())
}