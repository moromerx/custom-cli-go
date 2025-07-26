package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const sentinel = "__DONE__" // to know when the shell is done producing output

func startCLI() {

	// starting shell for faster output

	var shell *exec.Cmd

	// Pick the shell based on OS
	if runtime.GOOS == "windows" {
		// PowerShell: no banner, no profile, read commands from stdin (“-Command -”)
		shell = exec.Command("powershell", "-NoLogo", "-NoProfile", "-Command", "-")
	} else {
		// Bash (Linux/macOS): skip rc/profile files so no prompt/banner is printed
		shell = exec.Command("bash", "--noprofile", "--norc")
	}

	// Get input/output pipes
	stdin, _ := shell.StdinPipe()
	stdout, _ := shell.StdoutPipe()
	stderr, _ := shell.StderrPipe()

	// Start the shell process
	shell.Start()

	cleanup := func() {
		stdin.Close()
		shell.Wait()
	}

	out := bufio.NewReader(stdout)  // read synchronously
	errR := bufio.NewReader(stderr) // forward stderr manually

	printBanner()

	for {
		scanner := bufio.NewScanner(os.Stdin) // creates a new scanner that reads from standard input (the console).

		inputPrompt()

		scanner.Scan()          // waits for the user to enter a line of text and press Enter.
		input := scanner.Text() //  stores the line of input into the variable text.

		cleanedInput := cleanInput((input))
		availableCommands := getCommands(cleanup)

		if len(cleanedInput) == 0 {
			continue // If the user just presses enter and does not provide anything, print another >>, this is the same functionality as
			// a typical cli
		}

		commandName := cleanedInput[0]

		command, ok := availableCommands[commandName]

		if ok {
			command.callback()
			continue
		}

		// if we get pass this means we have a valid command for execution

		// If cd is used for changing directory or the user wants to clear the screen
		// seperate from other commands as these are normal execution commands
		if commandName == "cd" && len(cleanedInput) > 1 {
			err := os.Chdir(cleanedInput[1])
			if err != nil {
				fmt.Printf("cd: %v\n", err)
			}
			// Restart shell in new working directory
			stdin.Close()
			shell.Process.Kill()
			shell.Wait()

			// Restart new shell in updated dir
			if runtime.GOOS == "windows" {
				shell = exec.Command("powershell", "-NoLogo", "-NoProfile", "-Command", "-")
			} else {
				shell = exec.Command("bash", "--noprofile", "--norc")
			}

			shell.Dir, _ = os.Getwd() // set working directory

			stdin, _ = shell.StdinPipe()
			stdout, _ = shell.StdoutPipe()
			stderr, _ = shell.StderrPipe()
			out = bufio.NewReader(stdout)
			errR = bufio.NewReader(stderr)
			shell.Start()
			continue
		}

		if commandName == "clear" || commandName == "cls" {
			clearScreen()
			continue
		}

		fullCommand := strings.Join(cleanedInput, " ") + "; echo " + sentinel

		io.WriteString(stdin, fullCommand+"\n")

		for errR.Buffered() > 0 {
			line, _ := errR.ReadString('\n')
			fmt.Print(line)
		}

		for {
			line, _ := out.ReadString('\n')
			trimmed := strings.TrimSpace(line)
			if trimmed == sentinel {
				break
			}
			fmt.Print(line)
		}
		fmt.Println()
	}
}

// an interface for a cli command
type cliCommand struct {
	name        string       // the name of the command
	description string       // what the command does
	callback    func() error // the action the command will perform
}

// Getting commands
func getCommands(cleanup func()) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the help menu",
			callback:    callbackHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    callbackExit(cleanup),
		},
	}
}

// we are going to take input, clean it, and then return each word in a slice
func cleanInput(input string) []string {
	loweredInput := strings.ToLower(input) // convert it to lowercase
	words := strings.Fields(loweredInput)  // splitting it into words

	return words
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	printBanner()
}

func inputPrompt() {
	cwd, err := os.Getwd() // Getting the current directory the user is in

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s> ", cwd)
}

var banner string = `
                       __                                  .__  .__ 
  ____  __ __  _______/  |_  ____   _____             ____ |  | |__|
_/ ___\|  |  \/  ___/\   __\/  _ \ /     \   ______ _/ ___\|  | |  |
\  \___|  |  /\___ \  |  | (  <_> )  Y Y  \ /_____/ \  \___|  |_|  |
 \___  >____//____  > |__|  \____/|__|_|  /          \___  >____/__|
     \/           \/                    \/               \/       
`

func printBanner() {
	fmt.Println("\033[1;92m" + banner + "\033[0m")
}
