package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type item struct {
	value      string
	exp        int64
	lastAccess int64
	fetch      bool
}

type TTLMap struct {
	m map[string]*item
	l sync.Mutex
}

func New(ln int) (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*item, ln)}
	go func() {
		for now := range time.Tick(time.Minute) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix() > v.exp {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k, v string) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item{value: v, exp: math.MaxInt64, fetch: false}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k string) (v *string) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		if time.Now().Unix() > it.exp {
			delete(m.m, k)
			m.l.Unlock()
			return nil
		}

		v = &it.value
		it.lastAccess = time.Now().Unix()
		it.fetch = true
	} else {
		m.l.Unlock()
		return nil
	}

	m.l.Unlock()
	return
}

func (m *TTLMap) SetExpire(k string, exp int64) bool {
	m.l.Lock()
	var found bool
	if it, ok := m.m[k]; ok {
		if time.Now().Unix() > it.exp {
			delete(m.m, k)
			m.l.Unlock()
			return false
		}

		it.exp = exp
		found = true
	}
	m.l.Unlock()
	return found
}

func Write(c net.Conn, msg string) {
	c.Write([]byte(msg))
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
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

	m := New(10)
	reader := bufio.NewReader(c)

	for true {
		command, _ := reader.ReadString('\n')
		args := strings.Split(command[:len(command)-1], " ")

		cmd := args[0]
		if strings.EqualFold(cmd, "q") || strings.EqualFold(cmd, "quit") {
			os.Exit(0)
		} else if strings.EqualFold(cmd, "SET") {
			if len(args) < 2 {
				Write(c, "Too few arguments\n")
				continue
			} else if len(args) > 3 {
				Write(c, "Too many arguments\n")
				continue
			}

			keyName := args[1]
			var value string
			if len(args) == 2 {
				value = ""
			} else {
				value = args[2]
			}

			m.Put(keyName, value)
			Write(c, "OK\n")
		} else if strings.EqualFold(cmd, "GET") {
			if len(args) != 2 {
				Write(c, "Incorrect number of arguments\n")
				continue
			}

			keyName := args[1]
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

			keyName := args[1]
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
			if len(args) != 2 {
				Write(c, "Incorrect number of arguments\n")
				continue
			}

			keyName := args[1]
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
				continue
			}
		}
	}
}
