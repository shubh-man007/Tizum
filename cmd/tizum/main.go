package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shubh-man007/Tizu/internal/database"
	"github.com/shubh-man007/Tizu/internal/repository"
	"github.com/shubh-man007/Tizu/internal/tui"
)

var (
	version = "dev"
)

func appendBanner() {
	banner := `

████████╗██╗███████╗██╗   ██╗███╗   ███╗
╚══██╔══╝██║╚══███╔╝██║   ██║████╗ ████║
   ██║   ██║  ███╔╝ ██║   ██║██╔████╔██║
   ██║   ██║ ███╔╝  ██║   ██║██║╚██╔╝██║
   ██║   ██║███████╗╚██████╔╝██║ ╚═╝ ██║
   ╚═╝   ╚═╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝
   `

	fmt.Println(banner)
}

func checkSQLite(dbPath string) error {
	db, err := database.InitDB(dbPath)
	if err != nil {
		return err
	}
	defer db.DB.Close()

	return db.DB.Ping()
}

func runDoctor(dbPath string) {
	binPath, err := os.Executable()
	if err != nil {
		binPath = "<unknown>"
	}

	fmt.Printf("DB path:     %s\n", dbPath)
	err = checkSQLite(dbPath)
	if err != nil {
		fmt.Printf("DB:          Err: %s\n", err)
	} else {
		fmt.Printf("DB:          %d\n", 1)
	}
	fmt.Printf("OS:          %s\n", runtime.GOOS)
	fmt.Printf("Binary path: %s\n", binPath)
	fmt.Printf("tizum version: %s\n", version)
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not find home directory: %v", err)
	}

	tizumDir := filepath.Join(homeDir, ".tizum")
	if err := os.MkdirAll(tizumDir, 0755); err != nil {
		log.Fatalf("Could not create config directory: %v", err)
	}

	dbPath := filepath.Join(tizumDir, "tizu.db")

	if len(os.Args) < 2 {
		appendBanner()
		fmt.Println("Usage: tizum <add|list|delete|toggle|edit|tui|doctor>")
		return
	}

	if os.Args[1] == "doctor" {
		appendBanner()
		runDoctor(dbPath)
		return
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	tiz := repository.NewTiz(db)

	switch os.Args[1] {

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tizum add <task>")
			return
		}
		err := tiz.CreateTask(os.Args[2])
		if err != nil {
			log.Println("Error:", err)
		}
		fmt.Println("Task added!")

	case "list":
		tasks, _ := tiz.ReadTasks()
		for i, t := range tasks {
			status := " "
			if t.Status {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, i+1, t.TaskName)
		}

	case "toggle":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tizum toggle <id>")
			return
		}
		pos, err := strconv.Atoi(strings.TrimSpace(os.Args[2]))
		if err != nil {
			log.Println("Error parsing: ", err)
			return
		}
		taskList, _ := tiz.ReadTasks()
		if pos < 1 || pos > len(taskList) {
			log.Printf("No task at position %d (valid: 1–%d)", pos, len(taskList))
			return
		}
		t := taskList[pos-1]
		tiz.ToggleTask(t.ID, !t.Status)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tizum delete <id1> [, <id2>, ...]")
			return
		}
		var positions []int
		for _, arg := range os.Args[2:] {
			for _, part := range strings.Split(arg, ",") {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				pos, err := strconv.Atoi(part)
				if err != nil {
					log.Printf("Error parsing %q: %v", part, err)
					continue
				}
				positions = append(positions, pos)
			}
		}
		tasks, _ := tiz.ReadTasks()
		deleted := 0
		for _, pos := range positions {
			if pos < 1 || pos > len(tasks) {
				log.Printf("No task at position %d (valid: 1–%d)", pos, len(tasks))
				continue
			}
			tiz.DeleteTask(tasks[pos-1].ID)
			deleted++
		}
		if deleted > 0 {
			if deleted == 1 {
				fmt.Println("Task deleted")
			} else {
				fmt.Printf("%d tasks deleted\n", deleted)
			}
		}

	case "edit":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tizum edit <id> <new_text>")
			return
		}
		pos, err := strconv.Atoi(strings.TrimSpace(os.Args[2]))
		if err != nil {
			log.Println("Error parsing: ", err)
			return
		}
		taskList, _ := tiz.ReadTasks()
		if pos < 1 || pos > len(taskList) {
			log.Printf("No task at position %d (valid: 1–%d)", pos, len(taskList))
			return
		}
		newText := strings.Join(os.Args[3:], " ")
		tiz.EditTask(taskList[pos-1].ID, newText)
		fmt.Println("Task updated")

	case "tui":
		m := tui.NewModel(tiz)
		if _, err := tui.Run(m); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
