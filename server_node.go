package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Write(c net.Conn, msg string) {
	c.Write([]byte(msg))
}

var m *TTLMap
var PORT string

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide port number")
		return
	}

	PORT = ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	m = New(10)
	reader := bufio.NewReader(c)

	for {
		command, _ := reader.ReadString('\n')
		if len(command) == 0 {
			continue
		}

		args := strings.Split(command[:len(command)-1], " ")
		cmd := args[0]
		if len(args) == 0 {
			Write(c, "Please type something")
			continue
		}

		if len(args) < 2 {
			Write(c, "Not enough arguments\n")
			continue
		}

		keyName := args[1]
		if strings.EqualFold(cmd, "SET") {
			if len(args) > 3 {
				Write(c, "Too many arguments\n")
				continue
			}

			var value string
			if len(args) == 2 {
				value = ""
			} else {
				value = args[2]
			}

			m.Put(keyName, value)
			Write(c, "OK\n")
		} else if strings.EqualFold(cmd, "GET") {
			v := m.Get(keyName)
			if v != nil {
				Write(c, *v+"\n")
			} else {
				Write(c, "<nil>\n")
			}
		} else if strings.EqualFold(cmd, "EXPIRE") {
			if len(args) != 3 {
				Write(c, "Incorrect number of arguments\n")
				continue
			}

			expTime, err := strconv.Atoi(args[2])
			if err != nil {
				Write(c, "Invalid argument\n")
				continue
			}

			code := 1
			ok := m.SetExpire(keyName, time.Now().Unix()+int64(expTime))
			if !ok {
				code = 0
			}

			Write(c, "(integer) "+strconv.Itoa(code)+"\n")
		} else if strings.EqualFold(cmd, "TTL") {
			if it, ok := m.m[keyName]; ok {
				if it.exp == math.MaxInt64 {
					Write(c, "-1\n")
				} else {
					ttl := it.exp - time.Now().Unix()
					Write(c, strconv.Itoa(int(ttl))+"\n")
					if ttl < 0 {
						delete(m.m, keyName)
					}
				}
			} else {
				Write(c, "Key does not exist\n")
			}
		} else if strings.EqualFold(cmd, "DELETE") {
			if _, ok := m.m[keyName]; ok {
				delete(m.m, keyName)
				Write(c, "OK\n")
			} else {
				Write(c, "Key not found\n")
			}
		} else if strings.EqualFold(cmd, "ADD") {
			if len(args) > 3 {
				Write(c, "Too many arguments\n")
				continue
			}

			var value string
			if len(args) == 2 {
				value = ""
			} else {
				value = args[2]
			}

			if _, ok := m.m[keyName]; ok {
				Write(c, "Item already exists\n")
				continue
			}

			m.Put(keyName, value)
			Write(c, "OK\n")
		}
	}
}
