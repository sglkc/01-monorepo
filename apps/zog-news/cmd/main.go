package main

import (
	"fmt"
	"os"
	"zog-news/cmd/commands"
	"zog-news/config"
)

func init() {
	config.LoadEnv()
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Expected a command.")
        os.Exit(1)
    }

    command := os.Args[1]
    args := os.Args[2:]

    err := commands.Execute(command, args)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
