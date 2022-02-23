// Package database provides concrete data access types.
package database

import (
	"errors"
	"todo-svc/foundation/model"
)

var (
	// ErrBadInput is returned when the input to a function does not satisfy base sanity checks.
	ErrBadInput = errors.New("bad input")
	// ErrNotFound is returned when data referenced by ID cannot be found.
	ErrNotFound = errors.New("not found")
)

/*
Persistence is the interface between applications and the underlying database.
*/
type Persistence interface {
	// CreateTodo creates a new TodoItem in the data store.
	// May return ErrBadInput if the input item is nil
	CreateTodo(item *model.TodoItem) (*model.TodoItem, error)
	// GetAllTodos gets all TodoItems in the store.
	// Will return nil if there are no TodoItems in the store.
	GetAllTodos() ([]*model.TodoItem, error)
	// GetTodo gets a single TodoItem from the store.
	// Will return ErrNotFound if the id does not match a known TodoItem.
	GetTodo(id uint64) (*model.TodoItem, error)
	// UpdateTodo updates the state of a TodoItem in the store.
	// Will return ErrNotFound if the id does not match a known TodoItem.
	UpdateTodo(id uint64, newState *model.TodoItem) error
	// DeleteTodo removes a TodoItem from the store.
	// Will return ErrBadInput if the newState is nil.
	// Will return ErrNotFound if the id does not match a known TodoItem.
	DeleteTodo(id uint64) error
}
