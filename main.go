package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

type Task struct {
	ID        int
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

var db *sql.DB

func main() {
	app := &cli.App{
		Name:  "Task Timer",
		Usage: "A CLI tool to manage tasks and their timers",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a new task",
				Action:  addTask,
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all tasks",
				Action:  listTasks,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "status",
						Usage: "Filter tasks by status",
					},
					&cli.StringFlag{
						Name:  "date",
						Usage: "Filter tasks by date (YYYY-MM-DD)",
					},
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start a task",
				Action:  startTask,
			},
			{
				Name:    "stop",
				Aliases: []string{"t"},
				Usage:   "Stop a task timer",
				Action:  stopTask,
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "d"},
				Usage:   "Delete a task",
				Action:  deleteTask,
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "Complete a task",
				Action:  completeTask,
			},
			{
				Name:    "save",
				Aliases: []string{"sv"},
				Usage:   "Save tasks to SQLite database",
				Action:  saveTasks,
			},
		},
	}

	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        start_time DATETIME,
        end_time DATETIME,
        status TEXT
    );
    `
	_, err := db.Exec(query)
	return err
}

func addTask(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("task name is required")
	}
	_, err := db.Exec("INSERT INTO tasks (name, status) VALUES (?, ?)", name, "pending")
	if err != nil {
		return err
	}
	fmt.Println("Task added:", name)
	return nil
}

func listTasks(c *cli.Context) error {
	statusFilter := c.String("status")
	dateFilter := c.String("date")

	query := "SELECT id, name, start_time, end_time, status FROM tasks WHERE 1=1"
	var args []interface{}

	if statusFilter != "" {
		query += " AND status = ?"
		args = append(args, statusFilter)
	}

	if dateFilter != "" {
		query += " AND DATE(start_time) = ?"
		args = append(args, dateFilter)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, status string
		var startTime, endTime sql.NullTime
		if err := rows.Scan(&id, &name, &startTime, &endTime, &status); err != nil {
			return err
		}

		var duration string
		if startTime.Valid && endTime.Valid {
			duration = formatDuration(endTime.Time.Sub(startTime.Time))
		} else if startTime.Valid {
			duration = "in progress"
		} else {
			duration = "not started"
		}

		fmt.Printf("%d: %s [%s]\n", id, name, status)
		fmt.Printf("   Duration: %s Start: %s End: %s\n", duration, formatTime(startTime.Time), formatTime(endTime.Time))
	}
	return nil
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

func startTask(c *cli.Context) error {
	id := c.Args().First()
	_, err := db.Exec("UPDATE tasks SET start_time = ?, status = ? WHERE id = ?", time.Now(), "in-progress", id)
	if err != nil {
		return err
	}
	fmt.Println("Task started:", id)
	return nil
}

func stopTask(c *cli.Context) error {
	id := c.Args().First()
	_, err := db.Exec("UPDATE tasks SET end_time = ?, status = ? WHERE id = ?", time.Now(), "stopped", id)
	if err != nil {
		return err
	}
	fmt.Println("Task stopped:", id)
	return nil
}

func deleteTask(c *cli.Context) error {
	id := c.Args().First()
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	fmt.Println("Task deleted:", id)
	return nil
}

func completeTask(c *cli.Context) error {
	id := c.Args().First()
	_, err := db.Exec("UPDATE tasks SET end_time = ?, status = ? WHERE id = ?", time.Now(), "completed", id)
	if err != nil {
		return err
	}
	fmt.Println("Task completed:", id)
	return nil
}

func saveTasks(c *cli.Context) error {
	// Tasks are already saved in the SQLite database
	fmt.Println("Tasks saved to SQLite database")
	return nil
}
