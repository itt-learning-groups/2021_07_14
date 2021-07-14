package storage

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestInMemoryDB_SaveTodo(t *testing.T) {
	testData := []struct {
		testName        string
		todoName        string
		todoDescription string
		db              DB
		expectedResult  Todo
		wantErr         bool
		expectedErr     error
	}{
		{
			testName:        "success",
			todoName:        "shopping",
			todoDescription: "get milk and eggs",
			db:              &InMemoryDB{},
			expectedResult: Todo{
				Name:        "shopping",
				Description: "get milk and eggs",
			},
		},
		{
			testName:        "failure: todo already in list",
			todoName:        "shopping",
			todoDescription: "get milk and eggs",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
			expectedResult: Todo{},
			wantErr:        true,
			expectedErr:    ErrAlreadyInList,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			result, err := db.SaveTodo(context.Background(), td.todoName, td.todoDescription)

			if !td.wantErr && err != nil {
				t.Fatalf("Save got unexpected error: %v", err)
			}

			if td.wantErr && !errors.Is(err, td.expectedErr) {
				t.Fatalf("SaveTodo expected error '%v'; got %v", td.expectedErr, err)
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
				t.Error("SaveTodo got a Todo with missing ID")
			}
		})
	}
}

func TestInMemoryDB_GetTodoList(t *testing.T) {
	testData := []struct {
		testName       string
		db             DB
		expectedResult []Todo
		wantErr        bool
	}{
		{
			testName: "success",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:   "22222bbb-bbbb-2222-b2bb-111aa1a11a1a",
						Name: "wash car",
					},
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
					{
						ID:          "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
						Name:        "walk dog",
						Description: "take dog to park",
					},
				},
			},
			expectedResult: []Todo{
				{
					ID:   "22222bbb-bbbb-2222-b2bb-111aa1a11a1a",
					Name: "wash car",
				},
				{
					ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
					Name:        "shopping",
					Description: "get milk and eggs",
				},
				{
					ID:          "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
					Name:        "walk dog",
					Description: "take dog to park",
				},
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			result, err := db.GetTodoList(context.Background())

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

func TestInMemoryDB_GetTodoByName(t *testing.T) {
	testData := []struct {
		testName       string
		todoName       string
		db             DB
		expectedResult Todo
		wantErr        bool
		expectedErr    error
	}{
		{
			testName: "success",
			todoName: "shopping",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
						Name:        "walk dog",
						Description: "take dog to park",
					},
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
			expectedResult: Todo{
				ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
				Name:        "shopping",
				Description: "get milk and eggs",
			},
		},
		{
			testName: "failure: todo not found",
			todoName: "wash car",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "33333ccc-cccc-3333-c3cc-111aa1a11a1a",
						Name:        "walk dog",
						Description: "take dog to park",
					},
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
			expectedResult: Todo{},
			wantErr:        true,
			expectedErr:    ErrNotFound,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			result, err := db.GetTodoByName(context.Background(), td.todoName)

			if !td.wantErr && err != nil {
				t.Fatalf("GetTodoByName got unexpected error: %+v", err)
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
				t.Errorf("GetTodoByName expected vs actual results don't match: %v", diff)
			}
		})
	}
}

func TestInMemoryDB_EditTodo(t *testing.T) {
	testData := []struct {
		testName       string
		db             DB
		todoID         string
		todoEdit       Todo
		expectedResult Todo
		wantErr        bool
	}{
		{
			testName: "success",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
			todoID: "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
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
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			result, err := db.EditTodo(context.Background(), td.todoID, td.todoEdit)

			if !td.wantErr && err != nil {
				t.Fatalf("EditTodo got unexpected error: %+v", err)
			}

			if td.wantErr && err == nil {
				t.Fatal("EditTodo expected an error; got none")
			}

			resultsCmp := cmp.Comparer(func(expected, actual Todo) bool {
				nameMatch := expected.Name == actual.Name
				descMatch := expected.Description == actual.Description
				return nameMatch && descMatch
			})

			if diff := cmp.Diff(td.expectedResult, result, resultsCmp); diff != "" {
				t.Errorf("EditTodo expected vs actual results don't match: %v", diff)
			}
		})
	}
}

func TestInMemoryDB_DeleteTodo(t *testing.T) {
	testData := []struct {
		testName    string
		db          DB
		todoID      string
		wantErr     bool
		expectedErr error
	}{
		{
			testName: "success",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
			todoID: "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
		},
		{
			testName: "failure: todo not found",
			db: &InMemoryDB{
				todoList: []Todo{},
			},
			todoID:  "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
			wantErr: true,
			expectedErr: ErrNotFound,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			err := db.DeleteTodo(context.Background(), td.todoID)

			if !td.wantErr && err != nil {
				t.Fatalf("EditTodo got unexpected error: %+v", err)
			}

			if td.wantErr && !errors.Is(err, td.expectedErr) {
				t.Fatalf("GetTodoByName expected error '%v'; got %v", td.expectedErr, err)
			}
		})
	}
}

func TestInMemoryDB_ClearTodoList(t *testing.T) {
	testData := []struct {
		testName    string
		db          DB
		wantErr     bool
	}{
		{
			testName: "success",
			db: &InMemoryDB{
				todoList: []Todo{
					{
						ID:          "11111aaa-aaaa-1111-a1aa-111aa1a11a1a",
						Name:        "shopping",
						Description: "get milk and eggs",
					},
				},
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			db := td.db
			err := db.ClearTodoList(context.Background())

			if !td.wantErr && err != nil {
				t.Fatalf("EditTodo got unexpected error: %+v", err)
			}

			if td.wantErr && err == nil {
				t.Fatal("EditTodo expected an error; got none")
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
