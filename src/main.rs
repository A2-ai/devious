use std::path::PathBuf;
use std::io::prelude::*;
use anyhow::Result;
use std::fs::{self, File};
use std::os::unix::fs::PermissionsExt;

mod helpers;
mod library;
use crate::library::init;
use crate::library::add;
use crate::library::get;
use crate::library::status;

fn main() -> Result<()> {
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

    // set permissions for files
    test1.metadata().unwrap().permissions().set_mode(0o777);
    test2.metadata().unwrap().permissions().set_mode(0o777);
    test3.metadata().unwrap().permissions().set_mode(0o777);

    // let _test1_mode = test1.metadata().unwrap().permissions().mode() & 0o777;
    // let _test2_mode = test2.metadata().unwrap().permissions().mode() & 0o777;
    // let _test3_mode = test3.metadata().unwrap().permissions().mode() & 0o777;
   
    // dvs init src/test_directory_storage --mode=0o764 --group=datascience
    let storage_dir = PathBuf::from(r"src/test_directory_storage");   
    let new_mode = 0o776; 
    let group = String::from("datascience");
    init::dvs_init(&storage_dir, &new_mode, &group)?;

    // dvs add src/test_directory/test1.txt src/test_directory/test2.txt src/test_directory/test3.txt "assembled data"
    let files: Vec<String> = vec![test1_path.clone(), test2_path.clone(), test3_path.clone()];
    let message = String::from("assembled data");
    add::dvs_add(&files, &message)?;

    // TODO: check permissions and group
    // let test1_storage = PathBuf::from("src/test_directory_storage/7f/08b8682ee8258389605201d65ed6a9104eed809c000d7975186bc4cd8a3efe");
    // let test2_storage = PathBuf::from("src/test_directory_storage/d3/9dfa7e18189f9d4cacabdaffc941191508ffac753e9eafa28155c154d76d5d");
    // let test3_storage = PathBuf::from("src/test_directory_storage/8e/287466df9f3e8b1a2bd177ba15efe111aae572ea8859bb24557cfb4418a5b4");
    // let test1_mode_new = test1_storage.metadata().unwrap().permissions().mode() & 0o777;
    // let test2_mode_new = test2_storage.metadata().unwrap().permissions().mode() & 0o777;
    // let test3_mode_new = test3_storage.metadata().unwrap().permissions().mode() & 0o777;
    // assert_eq!(0o777, test1_mode_new);
    // assert_eq!(0o777, test2_mode_new);
    // assert_eq!(0o777, test3_mode_new);

    // remove one of the files
    fs::remove_file(&test1_path)?;
    // change one of the files
    test2.write(b"\nadded a line")?;
    // keep test3.txt the same

    // dvs status
    let status = status::dvs_status(&Vec::new())?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    // dvs get src/test_directory/test1.txt 
    get::dvs_get(&vec![test1_path.clone()])?;

    // dvs status src/test_directory/test1.txt
    let status = status::dvs_status(&vec![test1_path.clone()])?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    // dvs add rc/test_directory/test2.txt "assembled data again"
    let message = String::from("assembled data again");
    add::dvs_add(&vec![test2_path.clone()], &message)?;

    // dvs status src/test_directory/test1.txt src/test_directory/test2.txt src/test_directory/test3.txt 
    let status = status::dvs_status(&vec![test1_path.clone(), test2_path.clone(), test3_path.clone()])?;
    let status_string = serde_json::to_string_pretty(&status).unwrap();
    println!("new status:\n{status_string}");

    fs::remove_file(&test1_path)?;
    fs::remove_file(&test2_path)?;
    fs::remove_file(&test3_path)?;
    Ok(())
 }