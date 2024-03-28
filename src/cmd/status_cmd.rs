use std::path::PathBuf;

use crate::internal::config::config;
use crate::internal::file::hash;
use crate::internal::git::repo;
use crate::internal::meta::{file, parse};

#[derive(Debug)]
pub struct JsonFileResult {
    pub path: PathBuf,
    pub status: String,
    pub file_size: u64,
    pub file_hash: String,
    pub time_stamp: String,
    pub saved_by: String,
    pub message: String
}

pub fn run_status_cmd(files: &Vec<String>) -> Result<Vec<JsonFileResult>, std::io::Error> {
  
    let mut json_logger: Vec<JsonFileResult> = Vec::new();

    let mut meta_paths: Vec<PathBuf> = Vec::new();

    // if no arguments are provided, get the status of all files in the current git repository
    if files.is_empty() {
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

        // get config
        match config::read(&git_repo) {
            Ok(config) => config,
            Err(_) => return Err(std::io::Error::other("devious is not initialized")),
        };

        // get meta files
        meta_paths = parse::get_all_meta_files(git_repo);
    } // if doing all files
    else {
        meta_paths = parse::parse_globs(files);
    } // else specific files

    if meta_paths.is_empty() {return Ok(json_logger)}

    json_logger  = meta_paths.into_iter().map(|path| {
        // get relative path
        let rel_path = repo::get_relative_path(&PathBuf::from("."), &path).expect("couldn't get relative path");
        
        // get file info
        let metadata = file::load(&path).expect("couldn't get metadata");
        
        // get whether file was hashable and file hash
        let file_hash_result = hash::get_file_hash(&path);

        let mut file_hash = String::from("placeholder");
        if file_hash_result.is_ok() {file_hash = file_hash_result.unwrap()}
        // else, file_hash_result was an error, so stick with the nonsense default value so they don't match

        // asign status: not-present by default
        let mut status = String::from("out-of-date");
        if !path.exists() {status = String::from("not-present")}
        else if file_hash == metadata.file_hash {
            status = String::from("up-to-date")
        }
        // else, the file exists, but the hash isn't up to date, so still with default: out-of-date

        // assemble info into JsonFileResult
        JsonFileResult{
            path: rel_path,
            status: status,
            file_size: metadata.file_size,
            file_hash: metadata.file_hash,
            time_stamp: metadata.time_stamp,
            saved_by: metadata.saved_by,
            message: metadata.message
        }
    }).collect::<Vec<JsonFileResult>>();

Ok(json_logger)
} // run_status_cmd