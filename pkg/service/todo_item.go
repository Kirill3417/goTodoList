package service

import (
	todo "awesomeProject"
	"awesomeProject/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetListById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateItem(listId, item)
}

func (s *TodoItemService) GetAllItem(userId, itemId int) ([]todo.TodoItem, error) {
	return s.repo.GetAllItem(userId, itemId)
}

func (s *TodoItemService) GetItemById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.UpdateItem(userId, itemId, input)
}
