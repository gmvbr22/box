package main

import (
	"fmt"
	"os"

	"github.com/gmvbr/box/pkg/lexical"
)
func main() {
    args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		file := args[0]
		if _, err := os.Stat(file); err == nil {
			lexical.ParseFile(file)
		} else {
			fmt.Printf("[err 0] File not found %s\n", file)
		}
	} else {
		fmt.Println("Required file argument")
	}
}