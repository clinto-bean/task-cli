package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

// Allocates a logger for the DB and calls ensureDB(), returning the DB and any error values
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

// Initializes the json DB file if it does not exist, returning errors if the call to writeDB fails
func (db *DB) createDB() error {
	db.log.Info("Creating DB.")
	dbStructure := DBStruct{
		Tasks: []Task{},
	}

	db.mu.RLock()
	err := db.writeDB(dbStructure)
	defer db.mu.RUnlock()

	return err
}

// Ensures the json DB file exists and if it does not, attempts to create it, returning any errors it encounters
func (db *DB) ensureDB() error {
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		db.log.Warn("JSON database not found, initializing.")
		return db.createDB()
	}
	if len(data) == 0 {
		os.WriteFile(db.path, []byte("{}"), 0600)
	}
	return err
}

// Attempts to read the file at the filepath for the DB, converts it to go struct and returns it as well as any error encountered
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

// Takes in a DBStruct struct and attempts to call writeDB with it, returning any errors that it may encounter
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

// Receives a task object and attempts to write it to the database, returning the database's task entry as well as any errors encountered
func (db *DB) AddTask(task Task) (Task, error) {
	dbStruct, err := db.loadDB()
	now := time.Now()
	if err != nil {
		db.log.Println("Could not load db")
		return Task{}, err
	}

	var idx int = 0
	// Iterating over the dbStruct.Tasks to ensure proper indexing even if a task is removed
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
	task.Status = "Incomplete"

	dbStruct.Tasks = append(dbStruct.Tasks, task)

	err = db.save(dbStruct)
	return task, err
}

// Attempts to load the DB and then returns all task data
func (db *DB) GetAllTasks() ([]Task, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		db.log.Println("Could not load db")
		return nil, err
	}
	return dbStruct.Tasks, nil
}

// Attempts to locate and return a specific task as well as any errors encountered
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

// Marks a task as complete and returns any errors encountered
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
	if dbStruct.Tasks[i].Status == "Complete" {
		return errors.New("already complete")
	}
	dbStruct.Tasks[i].Status = "Complete"
	dbStruct.Tasks[i].UpdatedAt = now
	err = db.save(dbStruct)
	return err
}

// Marks a task as in progress and returns any errors encountered
func (db *DB) StartTask(id int) error {
	now := time.Now()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	if dbStruct.Tasks == nil {
		return errors.New("no tasks")
	}
	i, err := index(dbStruct.Tasks, id)
	if err != nil || i < 0 {
		return err
	}
	if dbStruct.Tasks[i].Status != "Incomplete" {
		return errors.New("already started or completed")
	}
	dbStruct.Tasks[i].Status = "In Progress"
	dbStruct.Tasks[i].UpdatedAt = now
	err = db.save(dbStruct)
	return err
}

// Locates a task by .ID property and deletes it, returning any encountered errors
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
	err = db.save(dbStruct)
	return err
}

// Locates the task based on its .ID property and edits its description, returning any errors if encountered
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
	if dbStruct.Tasks[i].Status == "Complete" {
		return errors.New("already complete")
	}
	dbStruct.Tasks[i].UpdatedAt = now
	dbStruct.Tasks[i].Description = desc
	err = db.save(dbStruct)
	return err

}

// Returns any task where Status is Complete and also returns any error encountered
func (db *DB) GetCompletedTasks() ([]Task, error) {
	var tasks []Task
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Task{}, err
	}
	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Status == "Complete" {
			tasks = append(tasks, dbStruct.Tasks[i])
		}
	}
	return tasks, nil
}

// Returns all tasks where Status is not Complete and any errors encountered
func (db *DB) GetIncompleteTasks() ([]Task, error) {
	var tasks []Task
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Task{}, err
	}
	for i := range dbStruct.Tasks {
		if dbStruct.Tasks[i].Status != "Complete" {
			tasks = append(tasks, dbStruct.Tasks[i])
		}
	}
	return tasks, nil
}

// Marks a task as not started/incomplete via the undo command and returns an error if one occurs
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
	dbStruct.Tasks[i].Status = "Incomplete"
	err = db.save(dbStruct)
	return err
}

// Attempts to locate the task with a given ID, and will return its .ID property as well as any errors encountered
func index(tasks []Task, id int) (int, error) {
	for i, t := range tasks {
		if t.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

// Helper function to attempt to write to DB reducing the amount of code needed to be written. Takes in a DBStruct and returns error from writeDB
func (db *DB) save(payload DBStruct) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.writeDB(payload)
	if err != nil {
		fmt.Println("uh oh")
	}
	return err
}
