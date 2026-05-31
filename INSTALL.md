# Global CLI Deployment Guide

This document outlines how to transform a standalone Go binary into a globally accessible terminal command, just like `git` or `node`.

## 1. PATH Configuration: The Basics

Your operating system uses an environment variable called `PATH` to determine where to look when you type a command. If a binary is located inside any of the directories listed in your `PATH`, you can run it from anywhere.

### Windows (PowerShell/CMD)
1. **Find/Create a Folder**: A common standard is `C:\Users\<YourUsername>\.local\bin` or creating a specific app folder like `%USERPROFILE%\.warp-speed\bin`.
2. **Move the Binary**: Place `warp-speed.exe` into that folder.
3. **Update PATH (UI)**: 
   - Search for "Environment Variables" in the Start Menu.
   - Click "Edit the system environment variables" -> "Environment Variables".
   - Under "User variables", find `Path`, select it, and click "Edit".
   - Click "New" and paste the path to your bin folder. Click OK.
4. **Update PATH (PowerShell)**: 
   ```powershell
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Path\To\Your\Bin", [EnvironmentVariableTarget]::User)
   ```

### macOS / Linux (Zsh/Bash)
1. **Find/Create a Folder**: The standard user-level bin directory is `~/.local/bin`. System-wide is `/usr/local/bin` (requires `sudo`).
2. **Move the Binary**: `mv warp-speed ~/.local/bin/`
3. **Make Executable**: `chmod +x ~/.local/bin/warp-speed`
4. **Update PATH**:
   - For **Zsh** (Default on macOS): Add `export PATH="$PATH:$HOME/.local/bin"` to `~/.zshrc`.
   - For **Bash** (Linux): Add the same line to `~/.bashrc`.
5. **Reload Shell**: `source ~/.zshrc` or restart the terminal.

---

## 2. Best Practices: Symlinking vs. Copying

When installing a CLI tool, you can either **Copy** the compiled binary to the `PATH` directory, or create a **Symlink** (symbolic link or shortcut) from the `PATH` directory pointing to your local development folder.

### Copying (The standard for end-users)
* **How it works**: You take the compiled `warp-speed.exe` and duplicate it into `~/.local/bin/` or `%USERPROFILE%\.warp-speed\bin\`.
* **Pros**: It isolates the installed version from your source code. You can keep writing code, breaking things, and compiling new tests in your dev folder without affecting the global `warp-speed` command you rely on for daily use.
* **Cons**: If you *want* the global command to update, you must remember to copy it over again after building.

### Symlinking (The standard for active developers)
* **How it works**: You leave `warp-speed.exe` in your `D:\Project\Personal\CLI tool` directory. You create a special shortcut in your `PATH` directory that points to it.
* **Pros**: **Incredible Developer Experience (DX)**. Every time you run `go build -o warp-speed.exe`, your global command is instantly updated. No manual copying required. 
* **Cons**: If you delete or move your project folder, the global command breaks.

**Verdict**: As the developer of the tool, **Symlinking** is vastly superior for local development. For end-users installing your tool, **Copying** is the standard approach.

> **How to Symlink**:
> * **Windows (PowerShell run as Admin)**: `New-Item -ItemType SymbolicLink -Path "C:\Users\Admin\.local\bin\warp-speed.exe" -Target "D:\Project\Personal\CLI tool\warp-speed.exe"`
> * **macOS/Linux**: `ln -s "/path/to/your/repo/warp-speed" "$HOME/.local/bin/warp-speed"`

*(Note: The provided `install.sh` and `install.ps1` scripts use the Copying method, as it is the safest, most universal way to distribute a tool to end-users via a GitHub repo).*

---

## 3. The 'Pro' Install Scripts

I have created two professional installation scripts in your workspace:
- `install.ps1` for Windows users.
- `install.sh` for macOS/Linux users.

### How they work:
1. **OS Detection**: Handled implicitly by users choosing the correct script based on their shell (PowerShell vs Bash/Zsh).
2. **Directory Creation**: They check for and create a dedicated `bin` folder (`$HOME/.local/bin` for Unix, `$USERPROFILE\.warp-speed\bin` for Windows).
3. **Move/Copy**: They copy the compiled binary into the target directory.
4. **PATH Injection**: They intelligently append the new directory to the user's permanent environment variables or shell configuration files (`.zshrc`, `.bashrc`).

You can host these scripts on your GitHub repository. Users can download the binaries from your GitHub Releases page and run these scripts to install the CLI seamlessly!
