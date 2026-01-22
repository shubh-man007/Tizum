package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shubh-man007/Tizu/internal/database"
	"github.com/shubh-man007/Tizu/internal/repository"
	"github.com/shubh-man007/Tizu/internal/tui"
)

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

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	tiz := repository.NewTiz(db)

	if len(os.Args) < 2 {
		fmt.Println("Usage: tizum <add|list|delete|toggle|edit|tui>")
		return
	}

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
		for _, t := range tasks {
			status := " "
			if t.Status {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.TaskName)
		}

	case "toggle":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tizum toggle <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Println("Error parsing: ", err)
		}
		taskList, _ := tiz.ReadTasks()
		for _, t := range taskList {
			if t.ID == id {
				tiz.ToggleTask(id, !t.Status)
			}
		}

	case "delete":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Println("Error parsing: ", err)
		}
		tiz.DeleteTask(id)
		fmt.Println("Task deleted")

	case "edit":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tizum edit <id> <new_text>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Println("Error parsing: ", err)
		}
		newText := os.Args[3]
		tiz.EditTask(id, newText)
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
