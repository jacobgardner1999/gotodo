package store

import "fmt"

type InMemoryStore struct {
	users map[string]*User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: make(map[string]*User),
	}
}

func (s *InMemoryStore) CreateUser(username string) (id string, e error) {
	userID := fmt.Sprintf("%04d", len(s.users)+1)
	user := NewUser(userID, username)

	err := s.addUser(user)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return userID, nil
}

func (s *InMemoryStore) addUser(user User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}
	s.users[user.ID] = &user
	return nil
}

func (s *InMemoryStore) GetUser(userID string) (User, error) {
	if _, exists := s.users[userID]; !exists {
		return User{}, fmt.Errorf("no user found with ID %s", userID)
	}

	return *s.users[userID], nil
}

func (s *InMemoryStore) AddTodoList(list TodoList, userID string) error {
	user := s.users[userID]
	if _, exists := user.TodoLists[list.ID]; exists {
		return fmt.Errorf("list with ID %s for user %s already exists", list.ID, userID)
	}

	user.TodoLists[list.ID] = &list
	return nil
}

func (s *InMemoryStore) AddTodo(todo Todo, listID string, userID string) error {
	list := s.users[userID].TodoLists[listID]
	if _, exists := list.Todos[todo.ID]; exists {
		return fmt.Errorf("todo with ID %s in list ID %s for user ID %s already exists", todo.ID, listID, userID)
	}
	list.Todos[todo.ID] = &todo
	return nil
}

func (s *InMemoryStore) ToggleTodo(userID string, listID string, todoID string) error {
	user := s.users[userID]
	list := user.TodoLists[listID]

	if _, exists := list.Todos[todoID]; !exists {
		return fmt.Errorf("todo with ID %s in list ID %s for user ID %s does not exist", todoID, listID, userID)
	}

	todo := list.Todos[todoID]

	todo.Toggle()

	return nil
}
