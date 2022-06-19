package main

import (
	"bufio"
	"fmt"
	"math"
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
		for now := range time.Tick(time.Second) {
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

func (m *TTLMap) Get(k string) (v string) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
		it.fetch = true
	}
	m.l.Unlock()
	return
}

func (m *TTLMap) SetExpire(k string, exp int64) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.exp = exp
	}
	m.l.Unlock()
	return
}

func main() {
	m := New(10)
	reader := bufio.NewReader(os.Stdin)

	for true {
		command, _ := reader.ReadString('\n')
		args := strings.Split(command[:len(command)-1], " ")
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Please type something")
			continue
		}

		cmd := args[0]

		if strings.EqualFold(cmd, "q") || strings.EqualFold(cmd, "quit") {
			os.Exit(0)
		} else if strings.EqualFold(cmd, "SET") {
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Too few arguments")
				continue
			} else if len(args) > 3 {
				fmt.Fprintln(os.Stderr, "Too many arguments")
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
			fmt.Println("OK")
		} else if strings.EqualFold(cmd, "GET") {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "Incorrect number of arguments")
				continue
			}

			keyName := args[1]
			fmt.Println(m.Get(keyName))
		} else if strings.EqualFold(cmd, "EXPIRE") {
			if len(args) != 3 {
				fmt.Fprintln(os.Stderr, "Incorrect number of arguments")
				continue
			}

			keyName := args[1]
			expTime, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid argument")
				continue
			}

			m.SetExpire(keyName, time.Now().Unix()+int64(expTime))
			fmt.Println("(integer) 1")
		} else if strings.EqualFold(cmd, "TTL") {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "Incorrect number of arguments")
				continue
			}

			keyName := args[1]
			if it, ok := m.m[keyName]; ok {
				fmt.Println(it.exp - time.Now().Unix())
			} else {
				fmt.Fprintln(os.Stderr, "Read error")
				continue
			}
		}
	}
}
