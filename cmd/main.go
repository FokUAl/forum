package main

//"forumAA/web"

import (
	"fmt"
	"forumAA/web"
	"log"
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

	err := web.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
