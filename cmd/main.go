package main

//"forumAA/web"

import (
	"fmt"
	"forumAA/web"
	"os"
)

func main() {
	args := os.Args[1:]
	port := ":4888"

	if len(args) == 1 {
		port = args[0]
	} else if len(args) > 1 {
		fmt.Println("Number of arguments must be one.")
		os.Exit(1)
	}

	web.Run(port)
}
