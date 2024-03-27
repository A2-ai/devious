use std::path::PathBuf;

pub fn parse_globs(globs: &Vec<String>) -> Vec<PathBuf> {
    let mut meta_files: Vec<PathBuf> = Vec::new();

    for glob in globs {
        // remove meta file extension
        let file_string = glob.replace(".dvsmeta", "");
        let file_path = PathBuf::from(file_string);

        // skip is already queued
        if meta_files.contains(&file_path) {continue}

        // skip directories
        
        if file_path.is_dir() {
            // json skipping directory
            continue
        }

        meta_files.push(file_path);
    }
    
    return meta_files;
}