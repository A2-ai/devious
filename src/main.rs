use std::path::PathBuf;
// use clap::{arg, Command, ArgAction, Arg};
use std::io;
use std::fs;
use std::os::unix::fs::PermissionsExt;

mod internal;
mod cmd;
use crate::cmd::init_cmd;
use crate::cmd::add_cmd;
use crate::cmd::get_cmd;
use crate::cmd::status_cmd;

fn main() -> io::Result<()> {
    // dvs init src/test_directory_storage --mode=0o764 --group=datascience
   let storage_dir = PathBuf::from(r"src/test_directory_storage");   
   // if mode specified, set mode
   // else, mode is set to default                                                 
   let mode: u32 = 0o764;
   // if group specified, set group
   // else, group is set to default    
   let group = String::from("datascience");
   fs::set_permissions(&storage_dir.canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();
   init_cmd::run_init_cmd(&storage_dir, &mode, &group)?;
   println!("initialized devious");

   // dvs add src/test_directory/test1.txt src/test_directory/test2.txt "derived DA files"
   let files: Vec<String> = vec![String::from("src/test_directory/test1.txt"), String::from("src/test_directory/test2.txt")];
   let message = String::from("derived DA files");
   //add_cmd::run_add_cmd(&files, &message)?;

   //fs::remove_file(PathBuf::from("src/test_directory/test1.txt"))?;
//  fs::remove_file(PathBuf::from("src/test_directory/test2.txt"))?;

   // dvs get src/test_directory/test1.txt src/test_directory/test2.txt
   //get_cmd::run_get_cmd(&files)?;

   let empty_vec: Vec<String> = Vec::new();
   let vec = status_cmd::run_status_cmd(&empty_vec)?;
   println!("{:?}", vec);
    Ok(())
 }

    
    


 // test hash
    // let hash_path = PathBuf::from("/cluster-data/user-homes/jenna/Projects/devious/src/test_directory/test.txt");
    // let hash_output = hash::hash_file_with_blake3(&hash_path)?;
    // assert_eq!(hash_output, "71fe44583a6268b56139599c293aeb854e5c5a9908eca00105d81ad5e22b7bb6");
    // println!("hash is {hash_output}");

    // test copy
    // let src = PathBuf::from("src/test_directory/test.txt");
    // let dest = PathBuf::from("src/test_directory/test_copy.txt");
    // copy::copy(&src, &dest)?;
    // assert!(dest.exists());