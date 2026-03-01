# Skill Manager (skillmgr)

A lightweight CLI tool to manage and symlink your local skill repositories across projects efficiently

## Features

- **Repository Management**: Register local directories as "skill repositories".
- **Skill Linking**: Create symlinks from repositories to any target directory (e.g., project workspaces).
- **Style-Based Linking**: Use customizable styles to organize symlinks (default: `.opencode/skills/`).
- **Safe Moving**: Move local directories into a repository while automatically leaving a symlink behind.
- **Git Awareness**: Automatically detects and leaves `.git` directories in their original location during `move` (customizable).
- **Conflict Detection**: Refuses to overwrite existing files or symlinks unless forced.
- **Dry Run Support**: Preview filesystem changes with `--dry-run` on all modifying commands.
- **Status Reporting**: Check the health and source of symlinks in your workspace.
- **Configuration Management**: Set default styles and other preferences.

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed.

```bash
# Clone the repository
git clone <repository-url>
cd my_skill_manager

# Build the binary
go build -o skillmgr .

# (Optional) Move to your PATH
mv skillmgr /usr/local/bin/
```

## Quick Start

### 1. Register a Repository
Tell `skillmgr` where your skills are stored:
```bash
./skillmgr repo add my-skills ~/Documents/my-skill-vault
```

### 2. Move a Local Skill into the Repository
Move a directory to the vault and keep it accessible via symlink:
```bash
./skillmgr move ./my-new-component my-skills:ui/button
```
*Note: If `./my-new-component` is a git repo, the `.git` folder stays in your local workspace.*

### 3. Link an Existing Skill to a Project
```bash
./skillmgr link my-skills ui/button --target ./src/components/
```
*Note: By default, links are created in `./.opencode/skills/` relative to the current directory. You can customize the style.*

### 4. Customize Default Style
Set your preferred style for organizing links:
```bash
./skillmgr config set-default-style myproject
```
Now links will be created in `./.myproject/skills/`.

### 5. Check Symlink Status
```bash
./skillmgr status
```

## Commands

- `repo add <name> <path>`: Register a new skill repository.
- `repo list`: Show all registered repositories.
- `repo remove <name>`: Unregister a repository.
- `link <repo> <skill> [--target <dir>] [--style <style>]`: Create a symlink to a skill. Uses default style if not specified.
- `move <src> <repo>[:<dest>]`: Move a directory to a repo and link it back.
- `status`: Check the validity of symlinks in the current directory.
- `config set-default-style <style>`: Set the default style for links (default: opencode).
- `config get-default-style`: Show the current default style.

## Configuration

Settings are stored in `~/.skillmgr/config.json`. The directory and file are automatically created on first run.

- `default_style`: The default style for organizing symlinks (default: "opencode"). Links are created in `./.{style}/skills/`.

## Testing Your Installation

To verify `skillmgr` is working correctly on your machine:

1. **Build & Help**: `go build -o skillmgr . && ./skillmgr --help`
2. **Dry Run Test**: 
   ```bash
   mkdir -p ./test-src
   ./skillmgr move ./test-src my-repo --dry-run
   ```
3. **Status Check**: Run `./skillmgr status` in a directory containing symlinks to see them mapped to your repos.
4. **Config Test**: `./skillmgr config get-default-style` should show "opencode".

## License

[MIT](LICENSE)
