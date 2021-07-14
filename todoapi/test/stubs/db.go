package stubs

import (
	"context"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/internal/domain/storage"
)

type DBStub struct {
	SaveTodoFunc      func(ctx context.Context, name, description string) (storage.Todo, error)
	GetTodoListFunc   func(ctx context.Context) ([]storage.Todo, error)
	GetTodoByNameFunc func(ctx context.Context, name string) (storage.Todo, error)
	EditTodoFunc      func(ctx context.Context, id string, todo storage.Todo) (storage.Todo, error)
	DeleteTodoFunc    func(ctx context.Context, id string) error
	ClearTodoListFunc func(ctx context.Context) error
}

func (s DBStub) SaveTodo(ctx context.Context, name, description string) (storage.Todo, error) {
	return s.SaveTodoFunc(ctx, name, description)
}

func (s DBStub) GetTodoList(ctx context.Context) ([]storage.Todo, error) {
	return s.GetTodoListFunc(ctx)
}

func (s DBStub) GetTodoByName(ctx context.Context, name string) (storage.Todo, error) {
	return s.GetTodoByNameFunc(ctx, name)
}

func (s DBStub) EditTodo(ctx context.Context, id string, todo storage.Todo) (storage.Todo, error) {
	return s.EditTodoFunc(ctx, id, todo)
}

func (s DBStub) DeleteTodo(ctx context.Context, id string) error {
	return s.DeleteTodoFunc(ctx, id)
}

func (s DBStub) ClearTodoList(ctx context.Context) error {
	return s.ClearTodoListFunc(ctx)
}
