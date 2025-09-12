package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3" // Required for SQLite driver
)

// --- Database and PID Setup ---
const dbFile = "pomodoro.db"
const pidFile = "pomo-cli.pid"

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_name TEXT NOT NULL,
		duration_minutes INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func saveTask(db *sql.DB, taskName string, duration int, startTime, endTime time.Time) {
	insertSQL := `INSERT INTO tasks (task_name, duration_minutes, start_time, end_time) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, taskName, duration, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task '%s' saved to the database.\n", taskName)
}

// --- Timer Logic ---
func runTimer(taskName string, duration int) {
	startTime := time.Now()

	// Wait for the duration or until a termination signal is received
	timer := time.NewTimer(time.Duration(duration) * time.Minute)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-timer.C:
		// Timer finished
	case <-sigChan:
		// Termination signal received
		fmt.Printf("\nTimer for '%s' was stopped.\n", taskName)
		return
	}

	endTime := time.Now()
	
	fmt.Printf("Good job! You managed to complete the session for '%s'!\n", taskName)
	fmt.Println("\a") // Play a terminal beep sound
	
	db := initDB()
	defer db.Close()
	saveTask(db, taskName, duration, startTime, endTime)
}

// --- Stop Timer Logic ---
func stopTimer() {
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No active Pomodoro timer found in the background.")
			return
		}
		log.Fatalf("Error reading PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		log.Fatalf("Invalid PID in file: %v", err)
	}

	// Remove the PID file immediately to prevent race conditions
	os.Remove(pidFile)

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Error finding process: %v", err)
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		log.Fatalf("Error signaling process: %v", err)
	}

	fmt.Printf("Successfully stopped Pomodoro timer (PID: %d).\n", pid)
}

// --- View Tasks Logic ---
func viewTasks() {
	db := initDB()
	defer db.Close()

	rows, err := db.Query("SELECT task_name, duration_minutes, start_time FROM tasks ORDER BY start_time DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\n--- Pomodoro Task History ---")
	foundTasks := false
	for rows.Next() {
		var taskName string
		var duration int
		var startTimeStr string
		err := rows.Scan(&taskName, &duration, &startTimeStr)
		if err != nil {
			log.Fatal(err)
		}
		
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- Task: '%s' | Duration: %d mins | Started: %s\n", taskName, duration, startTime.Format("2006-01-02 15:04:05"))
		foundTasks = true
	}

	if !foundTasks {
		fmt.Println("No tasks found in the database.")
	}
}

// --- Main Function ---
func main() {
	// Define command-line flags
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	task := startCmd.String("task", "", "The name of the task to work on.")
	timer := startCmd.Int("time", 25, "The duration of the Pomodoro session in minutes. Default is 25.")
	background := startCmd.Bool("background", false, "Run the timer silently in the background.")
	
	// Internal flag to prevent recursive background spawning
	isBackgroundChild := startCmd.Bool("is-background-child", false, "")

	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: 'start', 'view', or 'stop'")
		os.Exit(1)
	}

	// Parse the subcommands
	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])

		if *task == "" {
			fmt.Println("--task flag is required for the 'start' command.")
			startCmd.PrintDefaults()
			os.Exit(1)
		}

		if *background && !*isBackgroundChild {
			// This is the parent process. We will fork a new process and exit.
			args := []string{"start", "--task", *task, "--time", strconv.Itoa(*timer), "--is-background-child"}
			cmd := exec.Command(os.Args[0], args...)
			
			// We redirect output to the current terminal so messages appear
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				log.Fatalf("Failed to start background process: %v", err)
			}
			
			// Save the child process ID and exit immediately
			pid := cmd.Process.Pid
			err = os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
			if err != nil {
				log.Printf("Warning: Could not save PID file: %v\n", err)
			}
			
			fmt.Printf("Pomodoro timer for '%s' is running in the background. PID: %d\n", *task, pid)
			fmt.Println("Use 'pomo-cli stop' to end it.")
			return
		}

		// This is the interactive or background child process
		runTimer(*task, *timer)

	case "view":
		viewCmd.Parse(os.Args[2:])
		viewTasks()
	case "stop":
		stopCmd.Parse(os.Args[2:])
		stopTimer()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
