use std::path::PathBuf;
use anyhow::{Context, Result};
use serde::{Serialize, Deserialize};
use crate::helpers::config;
use crate::helpers::hash;
use crate::helpers::repo;
use crate::helpers::file;
use crate::helpers::parse;


#[derive(Serialize, Deserialize, PartialEq, Debug)]
pub struct JsonFileResult {
    pub path: PathBuf,
    pub status: String,
    pub file_size: u64,
    pub file_hash: String,
    pub time_stamp: String,
    pub saved_by: String,
    pub message: String
}

pub fn dvs_status(files: &Vec<String>) -> Result<Vec<JsonFileResult>> {
    // struct for each file's status and such
    let mut json_logger: Vec<JsonFileResult> = Vec::new();

    // vector of files
    let mut meta_paths: Vec<PathBuf> = Vec::new();

    // if no arguments are provided, get the status of all files in the current git repository
    if files.is_empty() {
        // Get git root
        let git_dir = repo::get_nearest_repo_dir(&PathBuf::from(".")).with_context(|| "could not find git repo root - make sure you're in an active git repository")?;

        // get config
        config::read(&git_dir)?;

        // get meta files
       meta_paths = [meta_paths, parse::get_all_meta_files(git_dir)].concat();
    } // if doing all files
    else {
        meta_paths = [meta_paths, parse::parse_globs(files)].concat();
    } // else specific files

    if meta_paths.is_empty() {return Ok(json_logger)}

    json_logger  = meta_paths.into_iter().map(|path| {
        // get relative path
        let rel_path = repo::get_relative_path(&PathBuf::from("."), &path).expect("couldn't get relative path");
        
        // get file info
        let metadata = file::load(&path).expect("couldn't get metadata");
        
        // asign status: not-present by default
        let mut status = String::from("out-of-date");
        // if the file path doesn't exist assign status to "not-present"
        if !path.exists() {status = String::from("not-present")}
        else {
            // get whether file was hashable and file hash
            match hash::get_file_hash(&path) {
                Some(file_hash) => {
                    if file_hash == metadata.file_hash {
                        status = String::from("up-to-date")
                    }
                }
                None => (),
            }; 
        }
       
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