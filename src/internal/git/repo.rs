use std::path::{PathBuf};

fn get_relative_path(root_dir: &PathBuf, file_path: &PathBuf) -> PathBuf {
    let abs_file_path = file_path.canonicalize().expect("no such file or directory");

    let abs_root_dir = root_dir.canonicalize().expect("no such file or directory");

    abs_file_path.strip_prefix(abs_root_dir).expect("paths not relative").to_path_buf()
}

fn is_git_repo(dir: &PathBuf) -> bool {
    dir.join(".git").is_dir()
}

pub fn get_nearest_repo_dir(dir: &PathBuf) -> Result<PathBuf, std::io::Error> {
    let mut directory = match dir.canonicalize() {
        Ok(directory) => directory,
        Err(e) => return Err(e)
    };

    if is_git_repo(&dir) {return Ok(directory)}

    while directory != PathBuf::from("/") {
        directory = match directory.parent() {
            Some(directory) => {
                if is_git_repo(&directory.to_path_buf()) {return Ok(directory.to_path_buf())}
                else {continue;}
            }
            None => directory,
        };
    }
    return Err(std::io::Error::other("no nearby git repo"));
}