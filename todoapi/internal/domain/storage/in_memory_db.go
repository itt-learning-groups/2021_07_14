package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type InMemoryDB struct {
	todoList []Todo
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{}
}

var ErrNotFound = errors.New("not found")
var ErrAlreadyInList = errors.New("todo already in list")

func (db *InMemoryDB) SaveTodo(ctx context.Context, name, description string) (Todo, error) {
	if _, err := db.GetTodoByName(ctx, name); err != ErrNotFound {
		return Todo{}, ErrAlreadyInList
	}

	todo := Todo{
		ID: createID(),
		Name: name,
		Description: description,
	}

	// save to memory
	db.todoList = append(db.todoList, todo)

	// in-memory DB never returns an error on save
	var err error = nil

	return todo, err
}

func (db *InMemoryDB) GetTodoList(ctx context.Context) ([]Todo, error) {
	// get list from memory
	list := db.todoList

	// in-memory DB never returns an error on get
	var err error = nil

	return list, err
}

func (db *InMemoryDB) GetTodoByName(ctx context.Context, name string) (Todo, error) {
	for _, todo := range db.todoList {
		if todo.Name == name {
			return todo, nil
		}
	}
	return Todo{}, ErrNotFound
}

func (db *InMemoryDB) EditTodo(ctx context.Context, id string, todo Todo) (Todo, error) {
	// find and edit matching Todo in memory
	var editedTodo Todo
	for i := range db.todoList {
		item := db.todoList[i]
		if item.ID == id {
			item.Name = todo.Name
			item.Description = todo.Description
			editedTodo = item
		}
	}

	// in-memory DB never returns an error on edit
	var err error = nil

	return editedTodo, err
}

func (db *InMemoryDB) DeleteTodo(ctx context.Context, id string) error {
	// find and delete matching Todo in memory
	for i := range db.todoList {
		item := db.todoList[i]
		if item.ID == id {
			if i == len(db.todoList) - 1 {
				db.todoList = db.todoList[:i]
			} else {
				db.todoList = append(db.todoList[:i], db.todoList[i+1:]...)
			}
			return nil
		}
	}

	return ErrNotFound
}

func (db *InMemoryDB) ClearTodoList(ctx context.Context) error {
	// clear the list in memory
	db.todoList = make([]Todo, 0)

	// in-memory DB never returns an error on clear-list
	var err error = nil

	return err
}

func createID() string {
	return uuid.NewString()
}
