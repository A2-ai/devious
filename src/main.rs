use std::path::PathBuf;
// use clap::{arg, Command, ArgAction, Arg};
use std::io::{self, prelude::*};
use std::fs::{self, File};
use std::os::unix::fs::PermissionsExt;

mod internal;
mod cmd;
use crate::cmd::init_cmd;
use crate::cmd::add_cmd;
use crate::cmd::get_cmd;
use crate::cmd::status_cmd;

fn main() -> io::Result<()> {
    // write files
    let test1_path = String::from("src/test_directory/test1.txt");
    let mut test1 = File::create(&test1_path)?;
    test1.write_all(b"Hello, test1!")?;

    let test2_path = String::from("src/test_directory/test2.txt");
    let mut test2 = File::create(&test2_path)?;
    test2.write_all(b"Hello, test2!")?;

    let test3_path = String::from("src/test_directory/test3.txt");
    let mut test3 = File::create(&test3_path)?;
    test3.write_all(b"Hello, test3!")?;

    // write permissions for files
    fs::set_permissions(&PathBuf::from(&test1_path).canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();
    fs::set_permissions(&PathBuf::from(&test2_path).canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();
    fs::set_permissions(&PathBuf::from(&test3_path).canonicalize().unwrap(), fs::Permissions::from_mode(0o777)).unwrap();

    // dvs init src/test_directory_storage --mode=0o764 --group=datascience
    let storage_dir = PathBuf::from(r"src/test_directory_storage");   
    let mode: u32 = 0o764;
    let group = String::from("datascience");
    init_cmd::run_init_cmd(&storage_dir, &mode, &group)?;

    // dvs add src/test_directory/test1.txt src/test_directory/test2.txt src/test_directory/test3.txt "assembled data"
    let files: Vec<String> = vec![test1_path.clone(), test2_path.clone(), test3_path.clone()];
    let message = String::from("assembled data");
    add_cmd::run_add_cmd(&files, &message)?;

    // remove one of the files
    fs::remove_file(&test1_path)?;
    // change one of the files
    test2.write(b"added a line")?;
    // keep test3.txt the same

    // dvs status
    let status = status_cmd::run_status_cmd(&Vec::new())?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    // dvs get src/test_directory/test1.txt 
    get_cmd::run_get_cmd(&vec![test1_path.clone()])?;

    // dvs status src/test_directory/test1.txt
    let status = status_cmd::run_status_cmd(&vec![test1_path.clone()])?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    // dvs add rc/test_directory/test2.txt "assembled data again"
    let message = String::from("assembled data again");
    add_cmd::run_add_cmd(&vec![test2_path.clone()], &message)?;

    // dvs status src/test_directory/test1.txt src/test_directory/test2.txt src/test_directory/test3.txt 
    let status = status_cmd::run_status_cmd(&vec![test1_path.clone(), test2_path.clone(), test3_path.clone()])?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    Ok(())
 }