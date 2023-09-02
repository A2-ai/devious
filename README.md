# 🌀 Devious
A file linker that allows you to version large or sensitivie files under Git.

Instead of tracking the file itself, which can result in bloated and sluggish to sync repositories, Devious has Git track a metadata file containing a reference to a particular version of a file. This allows you to version large files under Git without tracking the file directly in your repository.

## 📦 Installation
### Linux (User)
This installs the latest version of Devious to `~/.local/bin` and adds `~/.local/bin` to your PATH in `~/.profile`.
```
curl -o- https://raw.githubusercontent.com/A2-ai/devious/main/scripts/install_user.sh | bash
```

## 📚 Usage
### Provide a storage location
Start by navigating to a Git repository for which you'd like to version large files. You can then initialize Devious by telling it where you want to store tracked files for the current repository. **This directory should be accessible in a shared location to all future users of the repository.**
```
dvs init <storage-path>
```

### Adding files
Once Devious is intialized, you can start adding files.
```
dvs add <glob> <another-glob>
```
`dvs add` accepts one or more globs, each representing a file or set of files to be tracked. Ignores files outside of current git repository. For example, `dvs add *.png subdir/*.csv` will add all PNG files in the current directory and all CSV files in the `subdir` directory.

### Updating files
If you want to update a file after changing it, you can simply run `dvs add <glob>` again. Devious will automatically update the file's reference and add the new version to the storage location.

### Getting files
You can get file(s) by running
```
dvs get <glob> [another-glob]
```
`dvs get` works the same way as `dvs add`, using globs, but instead of adding files to the storage location, it will retrieve them and place them in the current directory.


### Listing tracked files
You can list all tracked files and their statuses by running
```
dvs status
```
or status for a specific file with
```
dvs status <path>
```

## 🧰 Building & Developing

### Prerequisites
- Go

### Updating dependencies
```
go mod download
```
