package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shubh-man007/Tizu/internal/database"
	"github.com/shubh-man007/Tizu/internal/repository"
	"github.com/shubh-man007/Tizu/internal/tui"
)

func main() {
	db, err := database.InitDB("tizu.db")
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
		id := atoi(os.Args[2])
		taskList, _ := tiz.ReadTasks()
		for _, t := range taskList {
			if t.ID == id {
				tiz.ToggleTask(id, !t.Status)
			}
		}

	case "delete":
		id := atoi(os.Args[2])
		tiz.DeleteTask(id)
		fmt.Println("Task deleted")

	case "edit":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tizum edit <id> <new_text>")
			return
		}
		id := atoi(os.Args[2])
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

func atoi(s string) int {
	var n int
	fmt.Sscan(s, &n)
	return n
}
