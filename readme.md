# Tizum

```
████████╗██╗███████╗██╗   ██╗███╗   ███╗
╚══██╔══╝██║╚══███╔╝██║   ██║████╗ ████║
   ██║   ██║  ███╔╝ ██║   ██║██╔████╔██║
   ██║   ██║ ███╔╝  ██║   ██║██║╚██╔╝██║
   ██║   ██║███████╗╚██████╔╝██║ ╚═╝ ██║
   ╚═╝   ╚═╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝           
 ```
Tizum is a lightweight terminal-based task manager written in Go. 

---

## Installation

Clone the repository:

```
git clone https://github.com/shubh-man007/Tizum.git
cd Tizu
```

---

## Commands

All commands follow the pattern:

```
tizum <command> [flags]
```

Available commands:

### Add a Task

```
tizum add --task="Buy groceries"
```

### List Tasks

```
tizum list
```

### Delete a Task

```
tizum delete <id>
```

### Toggle Task Completion

```
tizum toggle <id>
```

### Edit a Task

```
tizum edit --id=3 --task="Updated task text"
```

### Launch the Bubble Tea TUI

```
tizum tui
```

![Usage](https://github.com/shubh-man007/Tizum/blob/main/assets/tizum.png)


---

## Project Structure

```
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
    model.go
    update.go
    view.go

tizu.db
go.mod
go.sum
```

---

## Running the Application

### Linux / macOS

Simply run:

```
go run ./cmd/tizum add --task="Example task"
go run ./cmd/tizum list
go run ./cmd/tizum tui
```

To build a binary:

```
go build -o tizum ./cmd/tizum
```

Move it into your PATH if you want global access.

---

## Global Installation on Windows

To run `tizum` globally from any directory, follow these steps.

### 1. Create a bin directory (PowerShell)

```
mkdir C:\Users\<username>\bin
```

Add it to your PATH:

```
setx PATH "$Env:PATH;C:\Users\<username>\bin"
```

Restart PowerShell.

### 2. Build the Windows binary with CGO enabled

SQLite requires CGO. Windows binaries must be built inside Windows with a Windows C compiler installed (e.g., TDM-GCC).

Enable CGO:

```
setx CGO_ENABLED 1
```

Restart PowerShell.

Build the binary:

```
go build -o tizum.exe .\cmd\tizum
```

### 3. Move the binary to your PATH folder

```
mv -Force .\tizum.exe C:\Users\<username>\bin\
```

### 4. Copy the database to the user home directory

```
mkdir C:\Users\<username>\.tizum
cp .\tizu.db C:\Users\<username>\.tizum\
```

Test:

```
tizum list
```

---

## Global Installation on WSL

To run `tizum` globally from any directory in WSL, follow these steps.

### 1. Build the binary

From your project root:

```
go build -o tizum ./cmd/tizum
```

### 2. Move the binary to a global PATH location

```
sudo mv tizum /usr/local/bin/
sudo chmod +x /usr/local/bin/tizum
```

### 3. Set up the database

You have two options for database management:

#### Option A: Separate database for WSL

```
mkdir -p ~/.tizum
cp tizu.db ~/.tizum/
```

#### Option B: Shared database with Windows (Recommended)

Create a symlink to use the same database across Windows and WSL:

```
mkdir -p ~/.tizum
ln -s /mnt/c/Users/<username>/.tizum/tizu.db ~/.tizum/tizu.db
```

With a shared database, tasks added in Windows will appear in WSL and vice versa.

### 4. Test the installation

```
cd ~
tizum list
tizum tui
```

### Optional: Quick rebuild alias

Add this alias to your `~/.bashrc` or `~/.zshrc` for easier updates:

```
alias install-tizum='go build -o tizum ./cmd/tizum && sudo mv tizum /usr/local/bin/'
```

Then rebuild with:

```
install-tizum
```

---

## Handling the CGO Issue on Windows

If you see this error:

```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work.
```

It means the executable was built without CGO enabled or built from within WSL.
To fix this:

* Build the binary inside Windows PowerShell.
* Ensure CGO is enabled with `setx CGO_ENABLED 1`.
* Install a Windows C compiler (TDM-GCC or MinGW).
* Rebuild the binary.

WSL cannot produce a Windows CGO-enabled binary, so all Windows builds must happen natively.

---

## Database Storage Location

Tizum stores its database in the user's home directory to ensure consistent access from any working directory. The database locations are:

* **Windows**: `C:\Users\<username>\.tizum\tizu.db`
* **Linux/macOS/WSL**: `~/.tizum/tizu.db`

This approach follows standard conventions for CLI tools and separates executable files from data files. If you're using both Windows and WSL, you can share the same database by creating a symlink in WSL that points to the Windows database location, allowing seamless task management across both environments.
