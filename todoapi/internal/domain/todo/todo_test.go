package todo

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/internal/domain/storage"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/test/stubs"
	"reflect"
	"testing"
)

var simulatedDBError = errors.New("simulated DB error")

func TestSave(t *testing.T) {
	testData := []struct {
		testName        string
		todoName        string
		todoDescription string
		db              storage.DB
		expectedResult  Todo
		wantErr         bool
		expectedError   error
	}{
		{
			testName:        "success",
			todoName:        "shopping",
			todoDescription: "get milk and eggs",
			db: stubs.DBStub{
				SaveTodoFunc: func(ctx context.Context, name, description string) (storage.Todo, error) {
					return storage.Todo{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					}, nil
				},
			},
			expectedResult: Todo{
				ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
				Name:        "shopping",
				Description: "get milk and eggs",
			},
		},
		{
			testName:        "failure: unspecified DB error",
			todoName:        "shopping",
			todoDescription: "get milk and eggs",
			db: stubs.DBStub{
				SaveTodoFunc: func(ctx context.Context, name, description string) (storage.Todo, error) {
					return storage.Todo{}, simulatedDBError
				},
			},
			wantErr: true,
			expectedError: simulatedDBError,
		},
		{
			testName:        "failure: attempt to recreate todo with same name",
			todoName:        "shopping",
			todoDescription: "get milk and eggs",
			db: stubs.DBStub{
				SaveTodoFunc: func(ctx context.Context, name, description string) (storage.Todo, error) {
					return storage.Todo{}, storage.ErrAlreadyInList
				},
			},
			wantErr: true,
			expectedError: storage.ErrAlreadyInList,
		},
	}

	for _ , td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			result, err := Save(context.Background(), td.db, td.todoName, td.todoDescription)

			if !td.wantErr && err != nil {
				t.Fatalf("Save got unexpected error: %v", err)
			}

			if td.wantErr && !errors.Is(err, td.expectedError) {
				t.Fatalf("Save expected error '%v'; got %v", td.expectedError, err)
			}

			resultsCmp := cmp.Comparer(func(expected, actual Todo) bool {
				nameMatch := expected.Name == actual.Name
				descMatch := expected.Description == actual.Description
				return nameMatch && descMatch
			})

			if diff := cmp.Diff(td.expectedResult, result, resultsCmp); diff != "" {
				t.Errorf("Save expected vs actual results don't match: %v", diff)
			}

			if !td.wantErr && result.ID == "" {
				t.Error("Save got a Todo with missing ID")
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	testData := []struct {
		testName       string
		db             storage.DB
		expectedResult []Todo
		wantErr        bool
	}{
		{
			testName: "success",
			db: stubs.DBStub{
				GetTodoListFunc: func(ctx context.Context) ([]storage.Todo, error) {
					return []storage.Todo{
						{
							ID: "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
							Name: "shopping",
							Description: "get milk and eggs",
						},
						{
							ID: "22222bbb-bbbb-2222-b2bb-111aa1a11a1a",
							Name: "wash car",
						},
						{
							ID: "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
							Name: "walk dog",
							Description: "take dog to park",
						},
					}, nil
				},
			},
			expectedResult: []Todo{
				{
					ID: "22222bbb-bbbb-2222-b2bb-111aa1a11a1a",
					Name: "wash car",
				},
				{
					ID: "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
					Name: "shopping",
					Description: "get milk and eggs",
				},
				{
					ID: "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
					Name: "walk dog",
					Description: "take dog to park",
				},
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			result, err := GetAll(context.Background(), td.db)

			if !td.wantErr && err != nil {
				t.Fatalf("GetTodoList got unexpected error: %+v", err)
			}

			if td.wantErr && err == nil {
				t.Fatal("GetTodoList expected an error; got none")
			}

			if len(result) != len(td.expectedResult) {
				t.Errorf("GetTodoList expected a list with %d items; go %d", len(td.expectedResult), len(result))
			}

			for _, todo := range td.expectedResult {
				if !listContains(result, todo) {
					t.Errorf("GetTodoList result did not contain expected Todo item %+v", todo)
				}
			}
		})
	}
}

func TestEdit(t *testing.T) {
	testData := []struct {
		testName       string
		db             storage.DB
		todoName       string
		todoEdit       Todo
		expectedResult Todo
		expectedErr    error
		wantErr        bool
	}{
		{
			testName: "success",
			db: stubs.DBStub{
				GetTodoByNameFunc: func(ctx context.Context, name string) (storage.Todo, error) {
					return storage.Todo{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					}, nil
				},
				EditTodoFunc: func(ctx context.Context, id string, todo storage.Todo) (storage.Todo, error) {
					return storage.Todo{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "go shopping",
						Description: "milk, eggs",
					}, nil
				},
			},
			todoName: "shopping",
			todoEdit: Todo{
				Name:        "go shopping",
				Description: "milk, eggs",
			},
			expectedResult: Todo{
				ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
				Name:        "go shopping",
				Description: "milk, eggs",
			},
		},
		{
			testName: "failure: todo not already in list",
			db: stubs.DBStub{
				GetTodoByNameFunc: func(ctx context.Context, name string) (storage.Todo, error) {
					return storage.Todo{}, storage.ErrNotFound
				},
			},
			todoName: "shopping",
			todoEdit: Todo{
				Name:        "go shopping",
				Description: "milk, eggs",
			},
			expectedResult: Todo{},
			wantErr: true,
			expectedErr: storage.ErrNotFound,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			result, err := Edit(context.Background(), db, td.todoName, td.todoEdit)

			if !td.wantErr && err != nil {
				t.Fatalf("Edit got unexpected error: %+v", err)
			}

			if td.wantErr && !errors.Is(err, td.expectedErr) {
				t.Fatalf("GetTodoByName expected error '%v'; got %v", td.expectedErr, err)
			}

			resultsCmp := cmp.Comparer(func(expected, actual Todo) bool {
				nameMatch := expected.Name == actual.Name
				descMatch := expected.Description == actual.Description
				return nameMatch && descMatch
			})

			if diff := cmp.Diff(td.expectedResult, result, resultsCmp); diff != "" {
				t.Errorf("Edit expected vs actual results don't match: %v", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testData := []struct {
		testName    string
		db          storage.DB
		todoName    string
		wantErr     bool
		expectedErr error
	}{
		{
			testName: "success",
			db: stubs.DBStub{
				GetTodoByNameFunc: func(ctx context.Context, name string) (storage.Todo, error) {
					return storage.Todo{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					}, nil
				},
				DeleteTodoFunc: func(ctx context.Context, id string) error {
					return nil
				},
			},
			todoName: "shopping",
		},
		{
			testName: "failure: todo not already in list",
			db: stubs.DBStub{
				GetTodoByNameFunc: func(ctx context.Context, name string) (storage.Todo, error) {
					return storage.Todo{}, storage.ErrNotFound
				},
			},
			todoName: "shopping",
			wantErr: true,
			expectedErr: storage.ErrNotFound,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			err := Delete(context.Background(), db, td.todoName)

			if !td.wantErr && err != nil {
				t.Fatalf("EditTodo got unexpected error: %+v", err)
			}

			if td.wantErr && !errors.Is(err, td.expectedErr) {
				t.Fatalf("GetTodoByName expected error '%v'; got %v", td.expectedErr, err)
			}
		})
	}
}

func TestDeleteAll(t *testing.T) {
	testData := []struct {
		testName    string
		db          storage.DB
		wantErr     bool
	}{
		{
			testName: "success",
			db: stubs.DBStub{
				ClearTodoListFunc: func(ctx context.Context) error {
					return nil
				},
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			err := DeleteAll(context.Background(), db)

			if !td.wantErr && err != nil {
				t.Fatalf("EditTodo got unexpected error: %+v", err)
			}

			if td.wantErr && err == nil {
				t.Fatal("GetTodoList expected an error; got none")
			}
		})
	}
}

// listContains returns true if the given []Todo list contains a Todo item matching the one given as the 2nd argument
func listContains(list []Todo, match Todo) bool {
	for _, todo := range list {
		if reflect.DeepEqual(todo, match) {
			return true
		}
	}
	return false
}
