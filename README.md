Pomodoro CLI App
Description
This is a simple command-line interface (CLI) application to help you practice the Pomodoro productivity technique. It allows you to start a timed work session for a specific task and saves a record of your work. The app can run in the background and will provide a final status message upon completion.

Requirements & Installation
This application is written in Go, which provides a simple and portable executable.

Install Go: You must have Go installed on your system. You can download and install it from the official website. Below are simple steps for different operating systems.

For Windows:

Download the .msi installer.

Double-click the installer and follow the prompts. The installer will automatically set up the necessary environment variables.

For macOS:

Download the .pkg installer.

Double-click the installer and follow the prompts. The installer will place the Go binaries in /usr/local/go and set up your $PATH.

For Linux:

Download the .tar.gz archive.

Extract the archive to /usr/local with a command like sudo tar -C /usr/local -xzf go<version>.<os>-<arch>.tar.gz.

Add the Go binary path to your $PATH environment variable by adding export PATH=$PATH:/usr/local/go/bin to your ~/.profile or ~/.bashrc file and then restarting your terminal.

Initialize the project: In the directory with pomodoro.go, run these commands once to set up the project and download the database driver.

go mod init pomodoro
go mod tidy

Usage
Run the script with the executable created by go build, followed by one of the available commands.

Build the Executable:
You only need to do this once, or whenever you change the code.

go build pomodoro.go

This will create an executable named pomodoro in the same directory. To match your naming preference, you can rename it.

mv pomodoro pomo-cli

Installation (Making it a System Command)
To use pomo-cli from any directory, you need to place the executable in a directory that is part of your system's $PATH.

Move the pomo-cli executable to a common binaries directory, like /usr/local/bin/.

sudo mv pomo-cli /usr/local/bin/

You may be prompted for your password to complete this command.

Restart your terminal for the changes to take effect. You can now use pomo-cli from anywhere.

1. Start a New Pomodoro Session
Use the start command to begin a new timer.

--task <task_name> (required): The name of the task you are working on.

--time <minutes> (optional): The duration of the session in minutes. Defaults to 25.

--background (optional): Use this flag to run the timer silently in the background.

Interactive Mode (with output):

pomo-cli start --task "Write project documentation" --time 25

Background Mode (the simple way):
This is the ideal way to use the app for long sessions. The application will immediately run the timer in the background, and the "Good job!" message will appear in your terminal when it's done.

pomo-cli start --task "Focus on debugging" --time 50 --background

Your terminal prompt will return immediately, allowing you to continue working. The timer will run, and the final message will appear when it's finished.

2. View Task History
Use the view command to see all your completed Pomodoro sessions saved in the database.

pomo-cli view

Database
The application automatically creates a file named pomodoro.db in the same directory where you run your command to store all your completed tasks. This database file is a simple SQLite database and does not require any additional setup.