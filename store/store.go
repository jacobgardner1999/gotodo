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

func NewTodoList(id, name string) TodoList {
	return TodoList{
		ID:    id,
		Name:  name,
		Todos: make(map[string]Todo),
	}
}

func (s *InMemoryStore) CreateUser(user User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryStore) AddTodoList(list TodoList, userID string) error {
	user := s.users[userID]
	if _, exists := user.TodoLists[list.ID]; exists {
		return fmt.Errorf("list with ID %s for user %s already exists", list.ID, userID)
	}

	user.TodoLists[list.ID] = list
	return nil
}

func (s *InMemoryStore) AddTodo(todo Todo, listID string, userID string) error {
	list := s.users[userID].TodoLists[listID]
	if _, exists := list.Todos[todo.ID]; exists {
		return fmt.Errorf("todo with ID %s in list ID %s for user ID %s already exists", todo.ID, listID, userID)
	}
	list.Todos[todo.ID] = todo
	return nil
}
