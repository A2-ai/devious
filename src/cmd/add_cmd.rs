use std::path::PathBuf;
use crate::internal::git::repo;
use crate::config;

pub fn run_add_cmd(files: Vec<PathBuf>, message: String) -> Result<(), std::io::Error> {
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
    let config = config::read(&git_dir).expect("could not open yaml");

    let mut queued_paths: Vec<PathBuf> = Vec::new();

    files.into_iter().map(|file_in| {
        // remove meta file extension
        let string = file_in.display().to_string();
        let string_without_meta = string.replace(".dvsmeta", "");
        let file = PathBuf::from(string_without_meta);

        if queued_paths.contains(&file) {continue}


    });


    Ok(())
} // run_add_cmd