package repository

import (
	todo "awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) CreateItem(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItem := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItem, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItem(userId, itemId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id "+
		"INNER JOIN %s ul ON ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemsTable, userListsTable)

	if err := r.db.Select(&items, query, itemId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetItemById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id "+
		"INNER JOIN %s ul ON ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemsTable, userListsTable)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) DeleteItem(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul "+
		"WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2",
		todoItemsTable, listsItemsTable, userListsTable)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoItemsTable, setQuery, listsItemsTable, userListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
