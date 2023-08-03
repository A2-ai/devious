# 👺 Devious
Devious is a large file linker that allows you to store large files outside of your git repository. This allows you to version control your large files without having to worry about the size of your repository.

## 🪄 Usage
The CLI is invoked with `dvs` which stands for data versioning system.

### Prerequisites
Since we aren't distributing binaries yet, you'll have to build Lumos before you can use it. See the [Building](#Building) section for more information on building Devious.

### Initialize Devious
To begin using Devious, you must first specify a storage directory for a given repository. This directory will contain all of your large files. To initialize a storage directory, run
```
dvs init <storage-path>
```

### File operations
You can get a list of possible commands using
```
dvs help
```

## 🧰 Building & Developing

### Prerequisites
- Go
- Just

### Building
```
just build
```

### Developing
Update dependencies after cloning or pulling using
```
just update
```