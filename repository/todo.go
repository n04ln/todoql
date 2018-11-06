package repository

import (
	"database/sql"
	"fmt"

	"github.com/MasterMinds/squirrel"
	"github.com/NoahOrberg/todoql/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Repo interface {
	StoreTodo(m *model.Todo) error
	FindTodoByIDs(ids []string) ([]*model.Todo, error)
	FindTodoByUserIDs(ids []string) (map[string][]*model.Todo, error)
	FindTodoByID(id string) (*model.Todo, error)
}

type repository struct {
	db *sql.DB
}

func New() (Repo, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("root:pass@tcp(todoql_mysql:3306)/todoql?charset=utf8&parseTime=true&loc=Asia%%2FTokyo"))
	if err != nil {
		return nil, err
	}

	return &repository{
		db: db,
	}, nil
}

func (r *repository) FindTodoByIDs(ids []string) ([]*model.Todo, error) {
	q, args, err := squirrel.Select("id", "text", "done", "user_id").
		From("todos").
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	var gErr error
	res := make([]*model.Todo, 0, len(ids))
	for rows.Next() {
		t := new(model.Todo)
		if err := rows.Scan(&t.ID, &t.Text, &t.Done, &t.UserID); err != nil {
			errors.Wrap(gErr, err.Error())
			continue
		}
		res = append(res, t)
	}
	if gErr != nil {
		return nil, gErr
	}

	return res, nil
}

func (r *repository) FindTodoByUserIDs(ids []string) (map[string][]*model.Todo, error) {
	q, args, err := squirrel.Select("id", "text", "done", "user_id").
		From("todos").
		Where(squirrel.Eq{"user_id": ids}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	var gErr error
	res := map[string][]*model.Todo{}
	for rows.Next() {
		t := new(model.Todo)
		if err := rows.Scan(&t.ID, &t.Text, &t.Done, &t.UserID); err != nil {
			errors.Wrap(gErr, err.Error())
			continue
		}
		res[t.UserID] = append(res[t.UserID], t)
	}
	if gErr != nil {
		return nil, gErr
	}

	return res, nil
}

func (r *repository) FindTodoByID(id string) (*model.Todo, error) {
	q, args, err := squirrel.Select("id", "text", "done", "user_id").
		From("todos").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	t := new(model.Todo)
	if err = r.db.QueryRow(q, args...).Scan(&t.ID, &t.Text, &t.Done, &t.UserID); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *repository) StoreTodo(m *model.Todo) error {
	q, args, err := squirrel.Select("name").
		From("users").
		Where(m.UserID).
		Limit(1).
		ToSql()
	if err != nil {
		return err
	}
	var name string
	if err = r.db.QueryRow(q, args...).Scan(&name); err != nil {
		return errors.New("no such user")
	}

	_, err = squirrel.Insert("todos").
		Columns("id", "text", "done", "user_id").
		Values(m.ID, m.Text, m.Done, m.UserID).
		Exec()
	return err
}
