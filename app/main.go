package main

import (
	"fmt"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		var input string
		fmt.Print("$ ")
		fmt.Scan(&input)
		
		input = strings.TrimSpace(input)
		
		if input == "exit" {
			break
		} else if strings.HasPrefix(input, "echo ") {
			fmt.Println(input[5:])
		} else {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
