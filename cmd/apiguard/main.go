package main

import (
	"fmt"
	"os"

	"github.com/yokkesan/apiguard/internal/watcher"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("usage: apiguard <command>")
		return
	}

	switch args[0] {
	case "watch":
		path := "."

		if len(args) > 1 {
			path = args[1]
		}

		err := watcher.Watch(path)
		if err != nil {
			fmt.Println(err)
		}

	case "scan":
		fmt.Println("scan started")

	default:
		fmt.Printf("unknown command: %s\n", args[0])
	}
}
