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
	cd
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
	case cd:
		return "cd"
	default:
		return "unknown"
	}	
}

var builtins = map[string]bool {
	echo.String(): true,
	exit.String(): true,
	type_.String(): true,
	pwd.String(): true,
	cd.String(): true,
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
			case cd.String():
				handleCd(args)
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

func handlePwd(){
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting working directory: %s\n", err)
		return
	}
	fmt.Printf("%v\n", wd)
}

func handleCd(targetPath[] string) {
	if len(targetPath) == 0 || targetPath[0] == "" || targetPath[0] == "~" {
		err := os.Chdir(os.Getenv("HOME"))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	wd, _ := os.Getwd()
	if targetPath[0] == "-" {
		os.Chdir(os.Getenv("OLDPWD"))
		return
	}
	err := os.Chdir(targetPath[0])
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", targetPath[0])
	}
	os.Setenv("OLDPWD", wd)
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
