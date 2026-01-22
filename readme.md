# Tizum

Tizum is a lightweight terminal-based task manager written in Go. It provides both a traditional CLI interface and an optional TUI powered by the Bubble Tea framework. Tasks are stored locally using SQLite, making the tool portable and reliable without requiring any external database servers.

Tizum is designed for quick task entry, clean viewing, and efficient management directly from the command line.

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
tizum delete --id=3
```

### Toggle Task Completion

```
tizum toggle --id=3
```

### Edit a Task

```
tizum edit --id=3 --task="Updated task text"
```

### Launch the Bubble Tea TUI

```
tizum tui
```

---

## Bubble Tea

Tizum uses the Bubble Tea framework for its TUI mode. Bubble Tea is a functional and declarative TUI framework for Go that allows building clean and responsive terminal applications. The TUI displays tasks, allows navigation using the keyboard, and interacts with the same SQLite backend used by the CLI.

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

Test:

```
tizum list
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