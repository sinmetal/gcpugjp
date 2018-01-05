package example_appengine

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func setupTodo(swPlugin *swagger.Plugin) {
	s := &TodoService{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "TODO", Description: "TODO list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(s.Get)
	ucon.Handle("GET", "/todo/{id}", hInfo)
	hInfo.Description, hInfo.Tags = "get todo entity", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(s.List)
	ucon.Handle("GET", "/todo", hInfo)
	hInfo.Description, hInfo.Tags = "get todo list", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(s.Insert)
	ucon.Handle("POST", "/todo", hInfo)
	hInfo.Description, hInfo.Tags = "post new todo entity", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(s.Update)
	ucon.Handle("PUT", "/todo/{id}", hInfo)
	hInfo.Description, hInfo.Tags = "update todo entity", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(s.Delete)
	ucon.Handle("DELETE", "/todo/{id}", hInfo)
	hInfo.Description, hInfo.Tags = "delete todo entity", []string{tag.Name}
}

type TodoService struct {
}

// +jwg
// +qbg
type Todo struct {
	ID        int64  `datastore:"-" goon:"id"`
	Text      string `swagger:",req"`
	Done      bool
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (todo *Todo) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(todo, ps); err != nil {
		return err
	}

	return nil
}

func (todo *Todo) Save() ([]datastore.Property, error) {
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = time.Now()
	}
	todo.UpdatedAt = time.Now()

	ps, err := datastore.SaveStruct(todo)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

type ListOpts struct {
	Offset int `json:"offset" swagger:",in=query"`
	Limit  int `json:"limit" swagger:",in=query"`
}

type HttpError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("status %d: %s", err.Code, err.Text)
}

func (err *HttpError) StatusCode() int {
	return err.Code
}

func (err *HttpError) ErrorMessage() interface{} {
	return err
}

func (s *TodoService) Get(c context.Context, req *IntIDRequest) (*TodoJSON, error) {
	if req.ID == 0 {
		return nil, &HttpError{http.StatusBadRequest, "ID is required"}
	}

	g := goon.FromContext(c)

	todo := &Todo{ID: req.ID}
	err := g.Get(todo)
	if err == datastore.ErrNoSuchEntity {
		return nil, &HttpError{http.StatusNotFound, fmt.Sprintf("ID: %d is not found", req.ID)}
	} else if err != nil {
		return nil, err
	}

	resp, err := NewTodoJSONBuilder().AddAll().Convert(todo)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TodoService) List(c context.Context, r *http.Request, opts *ListOpts) ([]*TodoJSON, error) {
	// ucon can return []*Xxx. but CloudEndpoints can't it.

	// ucon can take 3 args. but go-endpoints can't it.
	{
		c := appengine.NewContext(r)
		appengine.AppID(c)
	}

	g := goon.FromContext(c)

	qb := NewTodoQueryBuilder()
	if opts.Limit != 0 {
		qb.Limit(opts.Limit)
	}
	if opts.Offset != 0 {
		qb.Offset(opts.Offset)
	}
	qb.CreatedAt.Asc()
	q := qb.Query()

	var todoList []*Todo
	_, err := g.GetAll(q, &todoList)
	if err != nil {
		return nil, err
	}

	resp, err := NewTodoJSONBuilder().AddAll().ConvertList(todoList)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TodoService) Insert(c context.Context, req *TodoJSON) (*TodoJSON, error) {
	if req == nil {
		return nil, &HttpError{http.StatusBadRequest, "Payload is required"}
	}
	if req.ID != 0 {
		return nil, &HttpError{http.StatusBadRequest, "ID should be 0"}
	}

	todo, err := req.Convert()
	if err != nil {
		return nil, err
	}

	g := goon.FromContext(c)
	_, err = g.Put(todo)
	if err != nil {
		return nil, err
	}

	resp, err := NewTodoJSONBuilder().AddAll().Convert(todo)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TodoService) Update(c context.Context, req *TodoJSON) (*TodoJSON, error) {
	if req.ID == 0 {
		return nil, &HttpError{http.StatusBadRequest, "ID is required"}
	}

	todo, err := req.Convert()
	if err != nil {
		return nil, err
	}

	old := &Todo{ID: todo.ID}

	g := goon.FromContext(c)

	err = g.RunInTransaction(func(g *goon.Goon) error {
		err = g.Get(old)
		if err == datastore.ErrNoSuchEntity {
			return &HttpError{http.StatusNotFound, fmt.Sprintf("ID: %d is not found", req.ID)}
		} else if err != nil {
			return err
		}

		_, err = g.Put(todo)
		if err != nil {
			return err
		}

		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	resp, err := NewTodoJSONBuilder().AddAll().Convert(todo)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TodoService) Delete(c context.Context, req *IntIDRequest) (*TodoJSON, error) {
	if req.ID == 0 {
		return nil, &HttpError{http.StatusBadRequest, "ID is required"}
	}

	todo := &Todo{ID: req.ID}

	g := goon.FromContext(c)

	err := g.Get(todo)
	if err == datastore.ErrNoSuchEntity {
		return nil, &HttpError{http.StatusNotFound, fmt.Sprintf("ID: %d is not found", req.ID)}
	} else if err != nil {
		return nil, err
	}

	key := g.Key(todo)

	err = g.Delete(key)
	if err != nil {
		return nil, err
	}

	resp, err := NewTodoJSONBuilder().AddAll().Convert(todo)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
