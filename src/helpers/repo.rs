use std::path::{PathBuf, Path};
use anyhow::{anyhow, Context, Result};
use std::fs;
use path_absolutize::*;

pub fn get_relative_path(root_dir: &PathBuf, file_path: &PathBuf) -> Result<PathBuf> {
    let abs_file_string = file_path.absolutize().unwrap().to_str().unwrap().to_string();
    let abs_file_path = PathBuf::from(abs_file_string);

    let abs_root_dir = root_dir.canonicalize()?;

    match abs_file_path.strip_prefix(abs_root_dir) {
        Ok(path) => return Ok(path.to_path_buf()),
        Err(e) => {
            return Err(anyhow!("could not get relative path for {} and {}: {e}", 
            root_dir.display(), file_path.display()))
        }
    }
}

fn is_git_repo(dir: &PathBuf) -> bool {
    dir.join(".git").is_dir()
}

pub fn is_directory_empty(directory: &Path) -> Result<bool> {
    let mut entries = fs::read_dir(directory)
    .with_context(|| format!("could not check if directory: {} is empty", directory.display()))?;
    Ok(entries.next().is_none())
}

pub fn get_nearest_repo_dir(dir: &PathBuf) -> Result<PathBuf> {
    let mut directory = dir.canonicalize().with_context(|| format!("could not find directory {}", dir.display()))?;

    if is_git_repo(&dir) {return Ok(directory)}

    while directory != PathBuf::from("/") {
        directory = match directory.parent() {
            Some(_) => {
                if is_git_repo(&directory.to_path_buf()) {return Ok(directory.to_path_buf())}
                else {directory.parent().unwrap().to_path_buf()}
            }
            None => directory,
        };
    }
    return Err(anyhow!("no nearby git repo"));
}