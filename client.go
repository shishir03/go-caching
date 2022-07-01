package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type serverNode struct {
	port string
}

func (s serverNode) String() string {
	return s.port
}

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Fprintln(os.Stderr, "Please provide host:port.")
		return
	}

	ch := ConsistentHash(len(arguments) - 2)
	for i := 2; i < len(arguments); i++ {
		ch.Add(serverNode{port: arguments[i]})
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		cmd := text[:len(text)-1]

		command, _ := reader.ReadString('\n')
		args := strings.Split(command, " ")
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Not enough arguments")
			continue
		}

		keyName := args[1]
		node := ch.LocateKey([]byte(keyName)).String()
		CONNECT := "localhost:" + node
		c, err := net.Dial("tcp", CONNECT)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(c, text)

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.EqualFold(cmd, "q") || strings.EqualFold(cmd, "quit") {
			fmt.Println("TCP client exiting...")
			return
		}

		c.Close()
	}
}
