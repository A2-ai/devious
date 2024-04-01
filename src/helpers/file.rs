use std::path::PathBuf;
use std::fs::File;
use serde::{Deserialize, Serialize};
use serde_json::Result;
use std::fs;

#[derive(Serialize, Deserialize)]
pub struct Metadata {
    pub file_hash: String,
    pub file_size: u64,
    pub time_stamp: String,
    pub message: String,
    pub saved_by: String
}

// static FILE_EXTENSION: String = String::from(".dvsmeta");

pub fn save(metadata: &Metadata, path: &PathBuf) -> Result<()> {
    // compose path file/to/file.ext.dvsmeta
    let metadata_file_path = PathBuf::from(path.display().to_string() + ".dvsmeta");

    // create file
    let _ = File::create(&metadata_file_path);
    // write to json
    let contents = serde_json::to_string_pretty(&metadata).unwrap();
    let _ = fs::write(&metadata_file_path, contents);
    Ok(())
}

pub fn load(path: &PathBuf) -> Result<Metadata> {
    let metafile_path = PathBuf::from(path.display().to_string() + ".dvsmeta");
    let contents = fs::read_to_string(metafile_path).unwrap();
    let metadata: Metadata = serde_json::from_str(&contents).unwrap();

    return Ok(metadata);
}