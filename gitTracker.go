package main

import (
	"flag"
	"fmt"
)

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "add folder to scan for Git repos")
	flag.StringVar(&email, "email", "example@email.com", "email to scan")
	flag.Parse()

	if folder != "" {
		Scan(folder)
	} else {
		fmt.Println("stats")
	}

}
