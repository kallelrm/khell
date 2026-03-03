package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"slices"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	commands := []string{"exit", "echo", "type"}
	reader := bufio.NewReader(os.Stdin)
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading command: ", err)
			os.Exit(1)
		}
		
		command = strings.TrimSpace(command)
		
		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		} else if strings.HasPrefix(command, "type ") {
			arg := command[5:]
			if slices.Contains(commands, arg) {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				fmt.Printf("%s: command not found\n", arg)
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
