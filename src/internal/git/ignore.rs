use std::{fs::{File, OpenOptions}, path::PathBuf};
use crate::internal::git::repo;
use std::io::prelude::*;

pub fn add_gitignore_entry(git_dir: &PathBuf, path: &PathBuf) -> Result<(), std::io::Error> {
    // get relative path
    let ignore_entry_temp = match repo::get_relative_path(git_dir, path) {
        Ok(entry) => entry,
        Err(e) => return Err(e)
    };
    // Add leading slash
    let path = ignore_entry_temp.display().to_string();
    let ignore_entry = format!("/{path}");

    // open the gitignore file, creating one if it doesn't exist
    let ignore_file = git_dir.join(".gitignore");
    if !ignore_file.exists() {
       File::create(&ignore_file)?;
    }
    let contents = std::fs::read_to_string(&ignore_file).unwrap();
    if !contents.contains(&ignore_entry) {
        let mut file = OpenOptions::new()
        .write(true)
        .append(true)
        .open(ignore_file)
        .unwrap();

        if let Err(e) = writeln!(file, "\n\n# Devious entry\n{ignore_entry}" ) {
            eprintln!("Couldn't write to file: {}", e);
        }

    } // add ignore entry
    Ok(())
}