use std::path::PathBuf;
use clap::{arg, Command, ArgAction, Arg};
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
    // let intro = Command::new("🌀 Devious")
    //     .version("0.0.1")
    //     .author("Andriy Massimilla, Jenna Johnson")
    //     .about("version large files under Git")
    //     .arg_required_else_help(false)
    //     .get_matches();

    // test init
    let storage_dir = PathBuf::from(r"src/test_directory");                                                      
    let mode: u32 = 0o664;
    let gid: u32 = 993;
    fs::set_permissions(&storage_dir.canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();
    init_cmd::get_init_runner(&storage_dir, &mode, &gid)?;
    fs::set_permissions(&storage_dir.canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();

    // test hash
    let hash_path = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src/test_directory/test.txt");
    let hash_output = hash::get_file_hash(&hash_path)?;
    println!("{hash_output}");

    // test copy
    let src = PathBuf::from("src/test_directory/test.txt");
    let dest = PathBuf::from("src/test_directory/test_copy.txt");
    copy::copy(&src, &dest)?;
    assert!(dest.exists());


    Ok(())
 }

    
    


