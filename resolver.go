//go:generate gorunpkg github.com/99designs/gqlgen

package todoql

import (
	context "context"

	"github.com/NoahOrberg/todoql/loader"
	model "github.com/NoahOrberg/todoql/model"
	"github.com/NoahOrberg/todoql/repository"
	"github.com/satori/uuid"
)

type Resolver struct {
	repo repository.Repo
}

func NewResolver(rep repository.Repo) (*Resolver, error) {
	return &Resolver{
		repo: rep,
	}, nil
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
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (model.Todo, error) {
	uuidv4, err := uuid.NewV4()
	if err != nil {
		return model.Todo{}, err
	}
	t := model.Todo{
		ID:     uuidv4.String(),
		Text:   input.Text,
		Done:   false,
		UserID: input.UserID,
	}
	if err := r.repo.StoreTodo(&t); err != nil {
		return model.Todo{}, nil
	}

	return t, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetUser(ctx context.Context, id string) (model.User, error) {
	return model.User{
		ID:   "qwer",
		Name: "noah",
	}, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (model.User, error) {
	return model.User{
		ID:   "qwer",
		Name: "noah",
	}, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Todos(ctx context.Context, obj *model.User) ([]model.Todo, error) {
	todos, err := loader.LoadTodoByUserID(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}
