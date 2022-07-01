package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide port number")
	}

	s := serverNode{port: arguments[1]}
	s.run()
}
