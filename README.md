# 🍅 Pomo CLI

A simple command-line Pomodoro timer. 


## Description

This is a CLI tool to help you manage your time using the Pomodoro Technique. You can start timers for specific tasks, view your task history, and stop running timers.

## Installation

To use this tool, you need to have [Go](https://golang.org/doc/install) installed on your system.

1.  **Clone or download the repository:**
    ```bash
    git clone https://github.com/arushdesp/pomo-cli
    cd pomo-cli
    ```

2.  **Build the executable:**
    ```bash
    go build -o pomo-cli pomo-cli.go
    ```
    This will create an executable file named `pomo-cli` for your current OS.

    To specifically create a Windows executable (`.exe`) from a non-Windows OS, use the following command:
    ```bash
    GOOS=windows GOARCH=amd64 go build -o pomo-cli.exe
    ```

### Adding to PATH for Easy Usage

To use `pomo-cli` from any directory, you should move the executable to a directory that is in your system's PATH.

#### macOS and Linux

1.  **Move the executable:**
    ```bash
    sudo mv pomo-cli /usr/local/bin/
    ```
    This moves the executable to a common directory for user-installed executables.

2.  **Make it executable:**
    ```bash
    sudo chmod +x /usr/local/bin/pomo-cli
    ```

Now you can run the tool by simply typing `pomo-cli` in your terminal.

### Windows TO CHECK after making the .exe ( TO-DO )

## Usage

The following commands are available:

### `start`

Starts a new Pomodoro timer.

**Flags:**

*   `--task`: (Required) The name of the task you are working on.
*   `--time`: The duration of the Pomodoro session in minutes. Defaults to 25.
*   `--background`: Run the timer silently in the background.

**Example:**

```bash
pomo-cli start --task "Write documentation" --time 30
```

To run in the background:
```bash
pomo-cli start --task "Code review" --background
```

### `view`

Displays the history of completed Pomodoro sessions.

**Example:**

```bash
pomo-cli view
```

### `stop`

Stops a Pomodoro timer that is running in the background.

**Example:**

```bash
pomo-cli stop
```

### `stats`
```bash
pomo-cli stats
```

Have a suggestion, question, or found a bug?  
👉 [Open an Issue](https://github.com/arushdesp/pomo-cli/issues)


## License

This project is licensed under the terms of the LICENSE file.
