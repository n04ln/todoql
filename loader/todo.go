package loader

import (
	"context"
	"errors"

	"github.com/NoahOrberg/todoql/model"
	"github.com/NoahOrberg/todoql/repository"
	"github.com/graph-gophers/dataloader"
)

func newTodoLoader(repo repository.Repo) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		results := make([]*dataloader.Result, len(keys))
		todoIDs := make([]string, 0, len(keys))
		userIDs := make([]string, 0, len(keys))
		for _, key := range keys {
			switch key := key.(type) {
			case todoIDKey:
				todoIDs = append(todoIDs, key.id)
			case userIDKey:
				userIDs = append(userIDs, key.id)
			}
		}
		todos, err := repo.FindTodoByIDs(todoIDs)
		if err != nil {
			return nil
		}
		todosByUserIDs, err := repo.FindTodoByUserIDs(userIDs)
		if err != nil {
			return nil
		}
		for i, key := range keys {
			results[i] = &dataloader.Result{Data: nil, Error: nil}
			switch key := key.(type) {
			case todoIDKey:
				for _, todo := range todos {
					if key.id == todo.ID {
						results[i].Data = todo
						continue
					}
				}
				if results[i].Data == nil {
					results[i].Error = errors.New("bookmark not found")
				}
			case userIDKey:
				results[i].Data = todosByUserIDs[key.id]
			}
		}
		return results
	}
}

const todoLoaderKey = "todoLoader"

type userIDKey struct {
	id string
}

func (key userIDKey) String() string {
	return key.id
}

func (key userIDKey) Raw() interface{} {
	return key.id
}

type todoIDKey struct {
	id string
}

func (key todoIDKey) String() string {
	return key.id
}

func (key todoIDKey) Raw() interface{} {
	return key.id
}

func LoadTodoByUserID(ctx context.Context, id string) ([]model.Todo, error) {
	ldr, err := getLoader(ctx, todoLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, userIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	d, ok := data.([]*model.Todo)
	if !ok {
		return nil, errors.New("not match type")
	}
	res := make([]model.Todo, 0, len(d))
	for _, v := range d {
		res = append(res, *v)
	}
	return res, nil
}
