use std::path::PathBuf;
// use clap::{arg, Command, ArgAction, Arg};
use std::io;
use std::fs;
use std::os::unix::fs::PermissionsExt;

mod internal;
mod cmd;
use crate::internal::config::config;
use crate::cmd::init_cmd;
use crate::internal::file::hash;
use crate::internal::storage::copy;



fn main() -> io::Result<()> {
    // test init
    let storage_dir = PathBuf::from(r"src/test_directory");                                                      
    let mode: u32 = 0o764;
    let group = String::from("datascience");

    fs::set_permissions(&storage_dir.canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();
    init_cmd::get_init_runner(&storage_dir, &mode, &group)?;
    println!("initialized devious");

    // test hash
    let hash_path = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src/test_directory/test.txt");
    let hash_output = hash::hash_file_with_blake3(&hash_path)?;
    assert_eq!(hash_output, "71fe44583a6268b56139599c293aeb854e5c5a9908eca00105d81ad5e22b7bb6");
    println!("hash is {hash_output}");

    // test copy
    let src = PathBuf::from("src/test_directory/test.txt");
    let dest = PathBuf::from("src/test_directory/test_copy.txt");
    copy::copy(&src, &dest)?;
    assert!(dest.exists());


    Ok(())
 }

    
    


