package main

import (
	"fmt"
	"sync"
)

type MutexDB struct {
	repository map[string]int
	gate       sync.Mutex
}

func NewMutexDB() *MutexDB {
	newMutexDB := MutexDB{
		repository: make(map[string]int),
	}
	return &newMutexDB
}

func (m *MutexDB) Add(key string, value int) error {
	m.gate.Lock()
	defer m.gate.Unlock()
	_, ok := m.repository[key]
	if ok {
		return fmt.Errorf("key %s already exist", key)
	}
	m.repository[key] = value
	return nil
}

func (m *MutexDB) Delete(key string) error {
	m.gate.Lock()
	defer m.gate.Unlock()
	_, ok := m.repository[key]
	if !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	delete(m.repository, key)
	return nil
}

func (m *MutexDB) Update(key string, value int) error {
	m.gate.Lock()
	defer m.gate.Unlock()
	_, ok := m.repository[key]
	if !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	m.repository[key] = value
	return nil
}

func (m *MutexDB) Find(key string) (value int, ok bool) {
	m.gate.Lock()
	defer m.gate.Unlock()
	value, ok = m.repository[key]
	return
}
