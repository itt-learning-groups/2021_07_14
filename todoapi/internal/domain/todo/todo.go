package todo

import (
	"context"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/internal/domain/storage"
)

func Save(ctx context.Context, db storage.DB, name, desc string) (Todo, error) {
	todo, err := db.SaveTodo(ctx, name, desc)
	if err != nil {
		return Todo{}, err
	}

	return Todo{
		ID: todo.ID,
		Name: todo.Name,
		Description: todo.Description,
	}, nil
}

func GetAll(ctx context.Context, db storage.DB) ([]Todo, error) {
	list, err := db.GetTodoList(ctx)
	if err != nil {
		return []Todo{}, err
	}

	var returnList []Todo
	for _, todo := range list {
		returnList = append(returnList, Todo{
			ID: todo.ID,
			Name: todo.Name,
			Description: todo.Description,
		})
	}

	return returnList, nil
}

func Edit(ctx context.Context, db storage.DB, name string, todo Todo) (Todo, error) {
	match, err := db.GetTodoByName(ctx, name)
	if err != nil {
		return Todo{}, err
	}

	editedTodo, err := db.EditTodo(ctx, match.ID, storage.Todo{
		Name: todo.Name,
		Description: todo.Description,
	})
	if err != nil {
		return Todo{}, err
	}

	return Todo{
		ID:          editedTodo.ID,
		Name:        editedTodo.Name,
		Description: editedTodo.Description,
	}, nil
}

func Delete(ctx context.Context, db storage.DB, name string) error {
	match, err := db.GetTodoByName(ctx, name)
	if err != nil {
		return err
	}

	err = db.DeleteTodo(ctx, match.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAll(ctx context.Context, db storage.DB) error {
	return db.ClearTodoList(ctx)
}
