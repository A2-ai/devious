use blake3::{self, hash};
use std::path::PathBuf;
use std::fs;

pub fn get_file_hash(path: &PathBuf) -> Result<String, std::io::Error> {
    // TODO: get cache if possible
    
    let file_contents = match fs::read_to_string(path) {
        Ok(file_contents) => {
            // json
            file_contents
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    let bytes = file_contents.as_bytes();

    let hash = hash(&bytes);

    // TODO: cache bytes

    return Ok(hash.to_string());
}

pub fn get_storage_path(storage_dir: &PathBuf, file_hash: &String) -> PathBuf {
    let first_hash_segment: &str = &file_hash[..2];
    let second_hash_segment: &str = &file_hash[2..];
    return storage_dir.join(first_hash_segment).join(second_hash_segment);

}