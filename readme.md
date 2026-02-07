# Tizum

```
████████╗██╗███████╗██╗   ██╗███╗   ███╗
╚══██╔══╝██║╚══███╔╝██║   ██║████╗ ████║
   ██║   ██║  ███╔╝ ██║   ██║██╔████╔██║
   ██║   ██║ ███╔╝  ██║   ██║██║╚██╔╝██║
   ██║   ██║███████╗╚██████╔╝██║ ╚═╝ ██║
   ╚═╝   ╚═╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝
```

Tizum is a lightweight terminal-based task manager written in Go. It uses SQLite for storage and a Bubble Tea TUI for interactive use. 

![TUI](https://github.com/shubh-man007/Tizum/blob/main/assets/tizum.png)
---

## Commands

Usage:

```
tizum <command> [arguments]
```

| Command | Description |
|---------|-------------|
| `add`   | Add a task |
| `list`  | List all tasks with position and status |
| `toggle`| Toggle completion for the task at the given position |
| `delete`| Delete one or more tasks by position |
| `edit`  | Edit the task at the given position |
| `tui`   | Launch the interactive TUI |
| `doctor`| Print diagnostic info (DB path, OS, binary path, version) |

### add

Add a new task. The task text is the single argument (use quotes for multiple words).

```
tizum add "Buy groceries"
tizum add Shower
```

### list

List all tasks. Each line shows status (`[x]` done, `[ ]` pending), position number, and task text. Positions are 1-based and used by `toggle`, `delete`, and `edit`.

```
tizum list
```

Example output:

```
[ ] 1: Buy groceries
[x] 2: Learn Go
[ ] 3: Review PR
```

### toggle

Toggle the completion state of the task at the given position.

```
tizum toggle 2
```

### delete

Delete tasks by position. One or more positions, separated by spaces or commas. Invalid positions are reported but do not stop other deletions.

```
tizum delete 3
tizum delete 2 4 6
tizum delete 1, 3, 5
```

### edit

Change the text of the task at the given position. Everything after the position is the new text (use quotes for multiple words).

```
tizum edit 1 "Buy groceries and milk"
tizum edit 2 Updated task text
```

### tui

Start the interactive terminal UI. You can browse tasks, filter, and use keys to add, edit, toggle, and delete. 

```
tizum tui
```

### doctor

Print diagnostic information: DB path, database status, OS, binary path, and tizum version. 

```
tizum doctor
```

---

## Project structure

```
Tizu/
  cmd/
    tizum/
      main.go
  internal/
    database/
      database.go
    models/
      tasks.go
    repository/
      crud.go
    tui/
      delegate.go
      model.go
      styles.go
      update.go
      view.go
  assets/
    Banner.md
    Tizum.png
  go.mod
  go.sum
```

---

## Building

**Requirements:** Go 1.24+ (or compatible). SQLite is used via `go-sqlite3`, which requires CGO on all platforms. On Windows, a C compiler (e.g. TDM-GCC, MinGW) is required for building.

Clone the repository:

```
git clone https://github.com/shubh-man007/Tizum.git
cd Tizum
```

Build a binary:

```
go build -o tizum ./cmd/tizum
```

On Windows, produce `tizum.exe`:

```
go build -o tizum.exe .\cmd\tizum
```

Optional: set version and CGO status for `tizum doctor`:

```
go build -ldflags "-X main.version=1.0.0 -X main.cgoStatus=enabled" -o tizum ./cmd/tizum
```

---

## Installation

### Linux and macOS

1. Build the binary (see above) and move it into your PATH, for example:

   ```
   go build -o tizum ./cmd/tizum
   sudo mv tizum /usr/local/bin/
   sudo chmod +x /usr/local/bin/tizum
   ```

2. The database is created automatically on first use under `~/.tizum/tizu.db`. No extra setup is required.

3. Run `tizum list` or `tizum tui` to verify.

### Windows (global install)

1. Create a directory in your PATH (e.g. `C:\Users\<username>\bin`) and add it to the PATH environment variable if needed.

2. Enable CGO and install a C compiler if you have not already:
   - Set `CGO_ENABLED=1` (e.g. `setx CGO_ENABLED 1` and restart the terminal).
   - Install TDM-GCC or MinGW so that `gcc` is available.

3. Build the binary from the project root **in Windows** (PowerShell or cmd):

   ```
   go build -o tizum.exe .\cmd\tizum
   ```

4. Move `tizum.exe` into your PATH directory (e.g. `Move-Item -Force .\tizum.exe C:\Users\<username>\bin\`).

5. The database is created automatically on first use under `C:\Users\<username>\.tizum\tizu.db`.

6. Run `tizum list` or `tizum doctor` to verify.

### WSL (global install)

1. Build from the project root in WSL:

   ```
   go build -o tizum ./cmd/tizum
   sudo mv tizum /usr/local/bin/
   sudo chmod +x /usr/local/bin/tizum
   ```

2. Choose how to store the database:

   - **Separate database for WSL:**  
     `mkdir -p ~/.tizum`  
     The DB file will be created at `~/.tizum/tizu.db` on first run.

   - **Shared database with Windows:**  
     Use the same file from both environments:
     ```
     mkdir -p ~/.tizum
     ln -s /mnt/c/Users/<username>/.tizum/tizu.db ~/.tizum/tizu.db
     ```
     Replace `<username>` with your Windows username. Tasks added in Windows will then appear in WSL and vice versa.

3. Run `tizum list` and `tizum tui` to verify. Use `tizum doctor` to confirm OS, binary path, and version if you use both WSL and Windows builds.

---

## CGO on Windows

Tizum depends on `go-sqlite3`, which requires CGO. If you see:

```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work.
```

then the binary was built without CGO (for example from WSL, which cannot produce a Windows CGO binary). To fix:

1. Build the Windows executable **on Windows** (PowerShell or cmd), not in WSL.
2. Ensure CGO is enabled: `setx CGO_ENABLED 1` and restart the terminal.
3. Install a Windows C compiler (TDM-GCC or MinGW) and ensure `gcc` is on the PATH.
4. Rebuild: `go build -o tizum.exe .\cmd\tizum`.

---

## Database location

The database file is stored in the user’s home directory so it is available from any working directory:

| Platform        | Path |
|-----------------|------|
| Windows         | `C:\Users\<username>\.tizum\tizu.db` |
| Linux / macOS / WSL | `~/.tizum/tizu.db` |

The `.tizum` directory is created automatically on first run. For a shared setup between Windows and WSL, use the symlink approach described in the WSL installation section.
