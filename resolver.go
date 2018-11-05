//go:generate gorunpkg github.com/99designs/gqlgen

package todoql

import (
	context "context"
	"fmt"
	"sync"

	model "github.com/NoahOrberg/todoql/model"
)

type Resolver struct {
	todos sync.Map
	next  int
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (model.Todo, error) {
	todo := model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("%d", r.next),
		UserID: input.UserID,
	}
	r.todos.Store(r.next, todo)
	r.next++
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]model.Todo, error) {
	todos := make([]model.Todo, 0, 0)
	r.todos.Range(func(k, v interface{}) bool {
		vv, ok := v.(model.Todo)
		if !ok {
			return true // continue
		}
		todos = append(todos, vv)
		return true
	})
	return todos, nil
}
func (r *queryResolver) Todo(ctx context.Context, userId string) ([]model.Todo, error) {
	todos := make([]model.Todo, 0, 0)
	r.todos.Range(func(k, v interface{}) bool {
		vv, ok := v.(model.Todo)
		if !ok {
			return true // continue
		}
		if vv.UserID == userId {
			todos = append(todos, vv)
		}
		return true
	})
	return todos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (User, error) {
	return User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}
