package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"os/exec"
)

type builtin int

const (
	echo builtin = iota
	exit
	type_
	pwd
)

func (b builtin) String() string {
	switch b{
	case echo:
		return "echo"
	case exit:
		return "exit"
	case type_:
		return "type"
	case pwd:
		return "pwd"
	default:
		return "unknown"
	}	
}

var builtins = map[string]bool {
	echo.String(): true,
	exit.String(): true,
	type_.String(): true,
}


func main() {
	for {
		fmt.Print("$ ")
		input, err := readFromStdin()
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading command: ", err)
			os.Exit(1)
		}
		
		cmd := input[0]	
		args := input[1:]
		PATH := os.Getenv("PATH")

		if len(cmd) == 0 {
			continue
		}

		switch cmd {
			case echo.String():
				handleEcho(args)
			case exit.String():
				return
			case type_.String():
				if len(args) == 0 {
					continue
				}
			  handleType(args[0], PATH)
			case pwd.String():
				handlePwd()
			default:
				handleProgram(PATH, cmd, args)
		}
	}
}

func readFromStdin() ([]string, error) {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return []string{}, fmt.Errorf("Error reading from stdin: %v", err)
	}

	return strings.Split(input[:len(input)-1], " "), nil
}

func handleEcho(args[] string){ 
	fmt.Println(strings.Join(args, " "))
}

func handleType(cmd, pathEnv string) {
	if len(cmd) == 0 {
		return
	}
	if ok := builtins[cmd]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	}

	pathElements := filepath.SplitList(pathEnv)

	for _, dir := range (pathElements) {
		fullpath := filepath.Join(dir, cmd)
		fileInfo, err := os.Stat(fullpath)
		if err != nil {
			continue
		}
		if !fileInfo.IsDir() && fileInfo.Mode()&0111 != 0 {
			fmt.Printf("%s is %s\n", cmd, fullpath)
			return
		} 
	}

	fmt.Printf("%s: not found\n", cmd)
}

func handlePwd() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting working directory: %s", err)
		return 
	}
	fmt.Printf("%v\n", pwd)
}

func handleProgram (pathEnv, cmd string, args []string) int {
	pathElements := filepath.SplitList(pathEnv)

	for _, dir := range (pathElements) {
		fullpath := filepath.Join(dir, cmd)
		fileInfo, err := os.Stat(fullpath)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() && fileInfo.Mode()&0111 != 0 {
			command := exec.Command(fullpath, args...)
			command.Args[0] = cmd

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Stdin = os.Stdin

			err := command.Run()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					return exitErr.ExitCode()
				}
				return 127
			}

			return 0
		} 
	}

	fmt.Printf("%s: command not found\n", cmd)

	return 0
}
