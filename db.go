package main

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStruct struct {
	Tasks map[int]Task `json:"tasks"`
}

func NewDB(path string) (*DB, error) {
	slog.Info("Initializing DB connection...")
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	slog.Info("Creating DB...")
	dbStructure := DBStruct{
		Tasks: map[int]Task{},
	}

	db.mu.RLock()
	err := db.writeDB(dbStructure)
	db.mu.RUnlock()

	return err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		slog.Warn("JSON database not found, initializing")
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStruct, error) {
	dbStructure := DBStruct{}

	db.mu.RLock()
	dat, err := os.ReadFile(db.path)
	db.mu.RUnlock()

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Could not read file: %v", db.path)
		return dbStructure, err
	}

	err = json.Unmarshal(dat, &dbStructure)

	if err != nil {
		log.Printf("Could not unmarshal data: %v", dat)
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStruct) error {
	slog.Info("Attempting to write to DB...")

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		log.Printf("Could not marshal data: %v", dat)
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		log.Printf("Could not write new data to db: %v", err)
		return err
	}
	return nil
}

func (db *DB) AddTask(t Task) (Task, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return Task{}, err
	}

	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Description == t.Description {
			return Task{}, errors.New("task already added")
		}
	}

	id := len(dbStruct.Tasks) + 1

	dbStruct.Tasks[id] = t

	db.mu.RLock()
	err = db.writeDB(dbStruct)
	db.mu.RUnlock()

	if err != nil {
		log.Println("Failed to write db")
		return Task{}, err
	}

	return dbStruct.Tasks[id], nil
}

func (db *DB) GetAllTasks() (map[int]string, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return nil, err
	}
	returnVals := make(map[int]string)
	for i := range dbStruct.Tasks {
		returnVals[i] = dbStruct.Tasks[i].Description
	}
	return returnVals, nil
}

func (db *DB) GetTask(id int) (Task, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return Task{}, err
	}
	var returnVals Task
	if dbStruct.Tasks[id].Description != "" {
		returnVals = dbStruct.Tasks[id]
		return returnVals, nil
	}
	return Task{}, errors.New("task not found")
}

func (db *DB) CompleteTask(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return err
	}
	t := dbStruct.Tasks[id]
	if t.Complete {
		return errors.New("already complete")
	}

	t.Complete = true
	dbStruct.Tasks[id] = t

	db.mu.Lock()
	err = db.writeDB(dbStruct)
	defer db.mu.Unlock()

	if err != nil {
		return err
	}
	if dbStruct.Tasks[id].Complete {
		log.Printf("Task %v completed", id)
		return nil
	}
	return errors.New("an error occurred")

}

func (db *DB) DeleteTask(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return err
	}
	_, ok := dbStruct.Tasks[id]
	if !ok {
		log.Println("task not found")
		return errors.New("task not found")
	}
	delete(dbStruct.Tasks, id)
	db.mu.Lock()
	defer db.mu.Unlock()
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	log.Printf("Task %v successfully deleted\n", id)
	return nil
}

func (db *DB) EditTask(id int, desc string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Could not load db")
		return err
	}
	_, ok := dbStruct.Tasks[id]
	if !ok {
		log.Println("task not found")
		return errors.New("task not found")
	}
	t := Task{Description: desc}
	dbStruct.Tasks[id] = t
	db.mu.Lock()
	defer db.mu.Unlock()
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	log.Printf("task %v successfully updated", id)
	return nil
}
