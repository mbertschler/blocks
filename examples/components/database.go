package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type DB struct {
	lock     sync.Mutex
	sessions map[string]*Session
	counters map[string]*StoredCounter
	todos    map[string]*StoredTodo
}

func NewDB() *DB {
	return &DB{
		sessions: make(map[string]*Session),
		counters: make(map[string]*StoredCounter),
		todos:    make(map[string]*StoredTodo),
	}
}

type Session struct {
	ID       string
	Created  time.Time
	LastSeen time.Time
	New      bool
}

type StoredCounter struct {
	ID    string
	Count int
}

func (db *DB) newSession() (*Session, error) {
	id, err := RandomString()
	if err != nil {
		return nil, fmt.Errorf("RandomString error: %w", err)
	}
	s := &Session{
		ID:       id,
		Created:  time.Now(),
		LastSeen: time.Now(),
		New:      true,
	}
	db.sessions[id] = s
	return s, nil
}

func (db *DB) GetSession(id string) (*Session, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	if id == "" {
		return db.newSession()
	}
	s, ok := db.sessions[id]
	if !ok {
		return db.newSession()
	}
	return s, nil
}

func (db *DB) SetSession(s *Session) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.sessions[s.ID] = s
	return nil
}

func (db *DB) GetCounter(id string) (*StoredCounter, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	c, ok := db.counters[id]
	if !ok {
		c = &StoredCounter{
			ID:    id,
			Count: 0,
		}
		db.counters[id] = c
	}
	return c, nil
}

func (db *DB) SetCounter(c *StoredCounter) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.counters[c.ID] = c
	return nil
}

func RandomString() (string, error) {
	// AES 192 == 24 bytes, so that should be enough
	// 24 bytes *8/6 = 32 bytes base64 encoded
	const length = 24
	buf := make([]byte, length)
	n, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("rand.Read failed: %w", err)
	}
	if n != length {
		return "", fmt.Errorf("short rand.Read: %d", n)
	}
	str := base64.URLEncoding.EncodeToString(buf)
	return str, nil
}

type StoredTodo struct {
	SessionID string
	Items     []StoredTodoItem
}

type StoredTodoItem struct {
	ID   int
	Done bool
	Text string
}

func (db *DB) GetTodo(id string) (*StoredTodo, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	t, ok := db.todos[id]
	if !ok {
		t = &StoredTodo{
			SessionID: id,
		}
		db.todos[id] = t
	}
	return t, nil
}

func (db *DB) SetTodo(t *StoredTodo) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.todos[t.SessionID] = t
	return nil
}
