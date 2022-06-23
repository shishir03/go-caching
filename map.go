package main

import (
	"math"
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
