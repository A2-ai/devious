use std::path::PathBuf;
use std::fs::{create_dir_all, File};
use std::os::unix::fs::PermissionsExt;
use std::io;
use std::fs;
use std::fs::Permissions;

pub fn copy(src_path: &PathBuf, dest_path: &PathBuf) -> Result<(), std::io::Error> {
    if *src_path == PathBuf::from(r"..") || *src_path == PathBuf::from(r".") {
        return Err(std::io::Error::other("copy failed"));
    }

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

    let src_file_data = match src_file.metadata() {
        Ok(data) => data,
        Err(e) => return Err(e),
    };

    let src_file_size = src_file_data.len();

    // ensure destination exists
    match create_dir_all(dest_path) {
        Ok(_) => {}
        Err(e) => return Err(e),
    }

    // create destination file
    let dest_file = match File::create(dest_path) {
        Ok(file) => file,
        Err(e) => return Err(e),
    };

    let mode: u32 = 0o664;
    let permissions: Permissions = fs::Permissions::from_mode(mode);
        match fs::set_permissions(dest_path.join(dest_file), permissions) {
            Ok(_) => {
                // json: success
            },
            Err(e) => {
                // json: fail
                return Err(e)
            }
        };

    Ok(())
}