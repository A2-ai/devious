# 🌀 Devious
A file linker that enables you to work with large files while keeping them under version control.

## 📦 Installation
### Linux (User)
This installs the latest version of Devious to `~/.local/bin` and adds `~/.local/bin` to your PATH.
```
curl -o- https://raw.githubusercontent.com/A2-ai/devious/main/scripts/install_user.sh | bash
```

## 🪄 Usage
The CLI is invoked with `dvs` which stands for data versioning system.

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