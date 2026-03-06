package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	// "os/exec"
	// "slices"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

type builtin int

const (
	echo builtin = iota
	exit
	type_
)

func (b builtin) String() string {
	switch b{
	case echo:
		return "echo"
	case exit:
		return "exit"
	case type_:
		return "type"
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
	forLoop: for {
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
				break forLoop
			case type_.String():
				if len(args) == 0 {
					continue
				}
			  handleType(args[0], PATH)
			default:
				fmt.Printf("%s: command not found\n", cmd)

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
		// files, _ := os.ReadDir(dir)
		// fmt.Println(files)
		// if slices.Contains([]string(files), cmd) {
		// 	fullPath := filepath.Join(dir, cmd)
		// 	permissions, _ := os.Stat(fullPath)
		//
		// 	if permissions.Mode()&0111 != 0 {
		// 		return
		// 	}
		// }
		// // if slices.Contains(cmd, files) {
		// // 	print(true)
		// // 	// program, err := os.Lstat(fmt.Sprintf("%s%s"))
		// // }
	}

	fmt.Printf("%s: not found\n", cmd)
}
