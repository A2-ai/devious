use std::path::PathBuf;
use walkdir::WalkDir;

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

pub fn get_all_meta_files(dir: PathBuf) -> Vec<PathBuf> {
    //let mut meta_files: Vec<String> = Vec::new();
    WalkDir::new(&dir)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.path().extension().map_or(false, |ext| ext == "dvsmeta"))
        .map(|e| {
            let string = e.into_path().display().to_string().replace(".dvsmeta", "");
            PathBuf::from(string)
        })
        .collect()
}