use std::path::{PathBuf, Path};
use std::fs;

pub fn get_relative_path(root_dir: &PathBuf, file_path: &PathBuf) -> Result<PathBuf, std::io::Error> {
    let abs_file_path = match file_path.canonicalize() {
        Ok(path) => path,
        Err(e) => return Err(e),
    };

    let abs_root_dir = match root_dir.canonicalize() {
        Ok(path) => path,
        Err(e) => return Err(e),
    };

    match abs_file_path.strip_prefix(abs_root_dir) {
        Ok(path) => return Ok(path.to_path_buf()),
        Err(_) => return Err(std::io::Error::other("paths not relative")),
    }
}

fn is_git_repo(dir: &PathBuf) -> bool {
    dir.join(".git").is_dir()
}

pub fn is_directory_empty(directory: &Path) -> std::io::Result<bool> {
    let mut entries = fs::read_dir(directory)?;
    let first_entry = entries.next();
    Ok(first_entry.is_none())
}

pub fn get_nearest_repo_dir(dir: &PathBuf) -> Result<PathBuf, std::io::Error> {
    let mut directory = match dir.canonicalize() {
        Ok(directory) => directory,
        Err(e) => return Err(e)
    };

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
    return Err(std::io::Error::other("no nearby git repo"));
}