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
	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide server node ports")
		return
	}

	ch := ConsistentHash(len(arguments) - 2)
	cmap := make(map[string]net.Conn)
	for i := 2; i < len(arguments); i++ {
		port := arguments[i]
		ch.Add(serverNode{port: port})
		CONNECT := "localhost:" + port
		c, err := net.Dial("tcp", CONNECT)
		if err != nil {
			fmt.Println(err)
			return
		}

		cmap[port] = c
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		command, _ := reader.ReadString('\n')
		args := strings.Split(command[:len(command)-1], " ")
		if strings.EqualFold(command, "q\n") || strings.EqualFold(command, "quit\n") {
			fmt.Println("TCP client exiting...")
			return
		}

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Not enough arguments")
			continue
		}

		keyName := args[1]
		node := ch.LocateKey([]byte(keyName)).String()
		c := cmap[node]

		fmt.Fprintf(c, command+"\n")

		message, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Server down")
			ch.Remove(node)
			continue
		}

		fmt.Print("->: " + message)
	}
}
