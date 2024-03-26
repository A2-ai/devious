use std::path::PathBuf;
use std::fs::{create_dir_all, File};
use std::os::unix::fs::PermissionsExt;
use std::fs;
use std::fs::Permissions;

pub fn copy(src_path: &PathBuf, dest_path: &PathBuf) -> Result<(), std::io::Error> {
    // Ignore .. and . paths
    if *src_path == PathBuf::from(r"..") || *src_path == PathBuf::from(r".") {
        return Err(std::io::Error::other("copy failed"));
    }

    // Open source file
    let src_file = match File::open(src_path) {
        Ok(file) => {
            // json
            file
        }
        Err(e) => {
            // json
            return Err(e)
        }
    };

    // Get file size
    let src_file_data = match src_file.metadata() {
        Ok(data) => data,
        Err(e) => return Err(e),
    };
    let _src_file_size = src_file_data.len();

    // ensure destination exists
    match create_dir_all(dest_path.parent().unwrap()) {
        Ok(_) => {}
        Err(e) => return Err(e),
    }

    // create destination file
    match File::create(dest_path) {
        Ok(file) => file,
        Err(e) => return Err(e),
    };

    // Set destination file permissions
    let mode: u32 = 0o664;
    let permissions: Permissions = fs::Permissions::from_mode(mode);
    fs::set_permissions(dest_path, permissions).expect("file unable to be copied, permissions denied");

    // Copy the file
    fs::copy(src_path, dest_path).expect("file unable to be copied");

    Ok(())
}