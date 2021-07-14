package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/internal/domain/storage"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/internal/domain/todo"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"net/http"
	"net/url"
)

type TodoListHandler struct {
	db       storage.DB
	validate *validator.Validate
}

func NewTodoListHandler() TodoListHandler {
	db := storage.NewInMemoryDB()
	return TodoListHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h TodoListHandler) EchoPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (h TodoListHandler) EchoPut(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	umBody := echoRequest{}
	err = json.Unmarshal(body, &umBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	param := vars["param"]

	echoMsg := echoPutResponse{
		Param:   param,
		ReqBody: umBody,
	}

	data, err := json.Marshal(echoMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (h TodoListHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	umBody := createRequest{}
	err = json.Unmarshal(body, &umBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.validate.Struct(umBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := todo.Save(ctx, h.db, umBody.Name, umBody.Description)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyInList) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(createResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (h TodoListHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	list, err := todo.GetAll(ctx, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var resp getAllResponse
	for _, item := range list {
		resp.List = append(resp.List, Todo{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h TodoListHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	vars := mux.Vars(r)
	todoName, err := url.QueryUnescape(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todoName == "" {
		http.Error(w, "missing 'name' parameter in request url", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	umBody := editRequest{}
	err = json.Unmarshal(body, &umBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.validate.Struct(umBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := todo.Edit(ctx, h.db, todoName, todo.Todo{
		Name: umBody.Name,
		Description: umBody.Description,
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(createResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (h TodoListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	vars := mux.Vars(r)
	todoName, err := url.QueryUnescape(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todoName == "" {
		http.Error(w, "missing 'name' parameter in request url", http.StatusBadRequest)
		return
	}

	err = todo.Delete(ctx, h.db, todoName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	return
}

func (h TodoListHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	err := todo.DeleteAll(ctx, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	return
}
