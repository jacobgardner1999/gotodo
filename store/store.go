package store

import "fmt"

type InMemoryStore struct {
	users map[string]User
}

type User struct {
	ID        string
	Name      string
	TodoLists map[string]TodoList
}

type TodoList struct {
	ID    string
	Name  string
	Todos map[string]Todo
}

type Todo struct {
	ID        string
	Title     string
	Completed bool
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: make(map[string]User),
	}
}

func NewUser(id, name string) User {
	return User{
		ID:        id,
		Name:      name,
		TodoLists: make(map[string]TodoList),
	}
}

func (s *InMemoryStore) CreateUser(user User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (u *User) AddTodoList(list TodoList) error {
	if _, exists := u.TodoLists[list.ID]; exists {
		return fmt.Errorf("list with ID %s for user %s already exists", list.ID, u.ID)
	}

	u.TodoLists[list.ID] = list
	return nil
}
