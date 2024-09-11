package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mu   *sync.RWMutex
	log  *Slogger
}

type DBStruct struct {
	Tasks []Task `json:"tasks"`
}

func NewDB(path string) (*DB, error) {
	log := &Slogger{}
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
		log:  log,
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	db.log.Info("Creating DB.")
	dbStructure := DBStruct{
		Tasks: []Task{},
	}

	db.mu.RLock()
	err := db.writeDB(dbStructure)
	db.mu.RUnlock()

	return err
}

func (db *DB) ensureDB() error {
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		db.log.Warn("JSON database not found, initializing.")
		return db.createDB()
	}
	if len(data) == 0 {
		os.WriteFile(db.path, []byte("{}"), 0666)
	}
	return err
}

func (db *DB) loadDB() (DBStruct, error) {
	dbStructure := DBStruct{}
	db.mu.RLock()
	dat, err := os.ReadFile(db.path)
	db.mu.RUnlock()

	if errors.Is(err, os.ErrNotExist) {
		db.log.Printf("Could not read file: %v\n", db.path)
		return dbStructure, err
	}

	err = json.Unmarshal(dat, &dbStructure)

	if err != nil {
		db.log.Printf("Could not unmarshal data: %v\n", dat)
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStruct) error {

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddTask(task Task) (Task, error) {
	dbStruct, err := db.loadDB()
	now := time.Now()
	if err != nil {
		db.log.Println("Could not load db")
		return Task{}, err
	}

	var idx int = 0

	for _, ts := range dbStruct.Tasks {
		if ts.ID > idx {
			idx = ts.ID
		}
	}

	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Description == task.Description {
			return Task{}, errors.New("task already added")
		}
	}

	task.ID = idx + 1
	task.CreatedAt = now
	task.UpdatedAt = now

	dbStruct.Tasks = append(dbStruct.Tasks, task)

	db.mu.Lock()
	err = db.writeDB(dbStruct)
	db.mu.Unlock()

	if err != nil {
		db.log.Println("Failed to write db")
		return Task{}, err
	}

	return task, nil
}

func (db *DB) GetAllTasks() ([]Task, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		db.log.Println("Could not load db")
		return nil, err
	}
	return dbStruct.Tasks, nil
}

func (db *DB) GetTask(id int) (Task, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Task{}, err
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil {
		return Task{}, err
	}
	return dbStruct.Tasks[i], nil
}

func (db *DB) CompleteTask(id int) error {
	now := time.Now()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil {
		return err
	}
	if dbStruct.Tasks[i].Complete {
		return errors.New("already complete")
	}
	dbStruct.Tasks[i].Complete = true
	dbStruct.Tasks[i].UpdatedAt = now
	db.mu.Lock()
	defer db.mu.Unlock()
	if err = db.writeDB(dbStruct); err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteTask(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil {
		return err
	}
	dbStruct.Tasks = append(dbStruct.Tasks[:i], dbStruct.Tasks[i+1:]...)
	db.mu.Lock()
	defer db.mu.Unlock()
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) EditTask(id int, desc string) error {
	now := time.Now()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil {
		return err
	}
	dbStruct.Tasks[i].UpdatedAt = now
	dbStruct.Tasks[i].Description = desc
	dbStruct.Tasks[i].Complete = false
	db.mu.Lock()
	defer db.mu.Unlock()
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetCompletedTasks() ([]Task, error) {
	var tasks []Task
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Task{}, err
	}
	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Complete {
			tasks = append(tasks, dbStruct.Tasks[i])
		}
	}
	return tasks, nil
}

func (db *DB) GetIncompleteTasks() ([]Task, error) {
	var tasks []Task
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Task{}, err
	}
	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Complete == false {
			tasks = append(tasks, dbStruct.Tasks[i])
		}
	}
	return tasks, nil
}

func (db *DB) IncompleteTask(id int) error {
	now := time.Now()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil {
		return err
	}
	dbStruct.Tasks[i].UpdatedAt = now
	dbStruct.Tasks[i].Complete = false
	err = db.save(dbStruct)
	return err
}

func index(tasks []Task, id int) (int, error) {
	for i, t := range tasks {
		if t.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func (db *DB) save(payload DBStruct) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.writeDB(payload)
	return err
}
