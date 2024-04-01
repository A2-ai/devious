use std::{fs::{File, OpenOptions}, path::PathBuf};
use crate::helpers::repo;
use std::io::prelude::*;
use anyhow::{Context, Result};

pub fn add_gitignore_entry(path: &PathBuf) -> Result<()> {
    let dir = path.parent().unwrap().to_path_buf();
    // get relative path
    let ignore_entry_temp = repo::get_relative_path(&dir, path).with_context(|| format!("could not create .gitignore entry for {}", path.display()))?;
    // Add leading slash
    let path = ignore_entry_temp.display().to_string();
    let ignore_entry = format!("/{path}");

    // open the gitignore file, creating one if it doesn't exist
    let ignore_file = dir.join(".gitignore");
    if !ignore_file.exists() {
       File::create(&ignore_file).with_context(|| format!("could not create gitignore file: {}", ignore_file.display()))?;
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