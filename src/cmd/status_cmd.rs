use std::time::SystemTime;
use std::path::PathBuf;

use crate::internal::config::config;
use crate::internal::file::hash;
use crate::internal::git::repo;
use crate::internal::meta::{file, parse};

pub struct JsonFileResult {
    pub path: String,
    pub status: String,
    pub file_size: u64,
    pub file_hash: String,
    pub time_stamp: SystemTime,
    pub saved_by: String,
    pub message: String
}

pub fn run_status_cmd(files: &Vec<PathBuf>) {
    




} // run_status_cmd