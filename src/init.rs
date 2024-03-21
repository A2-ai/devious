use crate::config::{self, Config};
use std::os::unix::fs::chown;
use std::path::PathBuf;
use std::fs::{self, create_dir, read_dir, metadata, File};
// use std::os::unix::fs::PermissionsExt;
// use std::os::unix::fs::MetadataExt;
//use std::os::unix::fs;

//const STORAGE_DIR_PERMISSIONS: st_mode = 0o644;



fn init(root_dir: PathBuf, storage_dir: PathBuf) -> Result<(), std::io::Error> { // 
    // get storage directory as an absolute path
    let storage_dir_abs: PathBuf = match storage_dir.canonicalize() {
        Ok(storage_dir_abs) => {
            // json: success
            storage_dir_abs
        },
        Err(e) => {
            // json: fail
            return Err(e)
        }
    };
    
    // check storage directory permissions
    //chown(dir, uid, gid)
    let md = fs::metadata(&storage_dir_abs)?;
    let permissions = md.permissions();
    let readonly = permissions.readonly();
    if readonly { // create storage dir
        // create storage dir
        match create_dir(&storage_dir_abs) {
            Ok(_) => {
                // json: success
            }
            Err(e) => {
                // json: fail
                return Err(e)
            }
        };

        // set permissions
        match fs::set_permissions(&storage_dir_abs, permissions) {
            Ok(_) => {
                // json: success
            },
            Err(e) => {
                // json: fail
                return Err(e)
            }
        };

        // json: success
    }
    else { // storage dir already exists
        // json: Storage directory already exists

        // Ensure storage dir is a directory
        if !storage_dir_abs.clone().is_dir() {
            // json: fail
            return Err(std::io::Error::other("storage dir is not a directory"));
        }

        // Warn if storage dir is not readable
        let mut dir = match read_dir(&storage_dir_abs) {
            Ok(dir) => {
                // json
                dir
            }
            Err(e) => {
                //json
                return Err(e);
            }
        };

        // Warn if storage dir is not empty
        if !dir.next().is_none() {
            // json warn
        }
    }

    // write config
    match config::write(&Config{storage_dir: storage_dir_abs.clone()}, &root_dir) {
        Ok(_) => {
            // json
        }
        Err(e) => {
            //json
            return Err(e);
        }
    };

    Ok(())
    // json: success
    
}