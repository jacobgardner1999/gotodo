package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type JsonStore struct {
	Users     map[string]*User
	storePath string
}

func NewJsonStore(storagePath string) (JsonStore, error) {
	return JsonStore{
		Users:     make(map[string]*User),
		storePath: storagePath,
	}, nil
}

func (s JsonStore) getUsersFromJson() (map[string]*User, error) {
	users := make(map[string]*User)
	file := s.storePath + "/users.json"

	jsonFile, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*User), nil
		} else {
			return make(map[string]*User), err
		}
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &users)

	return users, nil
}

func (s JsonStore) GetTodoLists(userID string) (map[string]*TodoList, error) {
	todos := make(map[string]*TodoList)
	file := s.storePath + "/" + userID + "lists.json"

	jsonFile, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return todos, nil
		} else {
			return todos, err
		}
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &todos)
	if err != nil {
		return todos, err
	}

	return todos, nil
}

func (s JsonStore) GetTodoList(userID string, listID string) (TodoList, error) {
	todoLists, err := s.GetTodoLists(userID)
	if err != nil {
		return TodoList{}, err
	}

	if _, exists := todoLists[listID]; !exists {
		return TodoList{}, fmt.Errorf("list with ID %s doesn't exist for user ID %s", listID, userID)
	}

	return *todoLists[listID], nil
}

func (s JsonStore) CreateUser(username string) (id string, e error) {
	users, err := s.getUsersFromJson()
	if err != nil {
		return "json error", fmt.Errorf("%s", err)
	}

	s.Users = users

	userID := fmt.Sprintf("%04d", len(s.Users)+1)
	user := NewUser(userID, username)

	s.Users[userID] = &user
	byteValue, err := json.MarshalIndent(s.Users, "", "  ")

	if err != nil {
		return "Error Marshalling", fmt.Errorf("%s", err)
	}

	err = os.WriteFile(s.storePath+"users.json", byteValue, 0644)

	if err != nil {
		return "Error Writing", fmt.Errorf("%s", err)
	}

	return userID, nil
}

func (s JsonStore) GetUser(userID string) (User, error) {
	users, err := s.getUsersFromJson()
	if err != nil {
		return User{}, fmt.Errorf("%s", err)
	}

	s.Users = users
	if _, exists := s.Users[userID]; !exists {
		return User{}, fmt.Errorf("no user found with ID %s", userID)
	}

	return *s.Users[userID], nil
}

func (s JsonStore) UpdateTodoList(list TodoList, userID string) error {
	todos, err := s.GetTodoLists(userID)
	if err != nil {
		return err
	}

	todos[list.ID] = &list
	byteValue, err := json.MarshalIndent(todos, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(s.storePath+"/"+userID+"lists.json", byteValue, 0644)

	if err != nil {
		return err
	}
	return nil
}

func (s JsonStore) DeleteTodoList(userID string, listID string) error {
	todos, err := s.GetTodoLists(userID)
	if err != nil {
		return err
	}

	delete(todos, listID)

	byteValue, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.storePath+"/"+userID+"lists.json", byteValue, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s JsonStore) AddTodo(todo Todo, listID string, userID string) error {
	list, err := s.GetTodoList(userID, listID)
	if err != nil {
		return err
	}

	list.Todos[todo.ID] = &todo

	err = s.UpdateTodoList(list, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s JsonStore) ToggleTodo(userID string, listID string, todoID string) error {
	list, err := s.GetTodoList(userID, listID)
	if err != nil {
		return err
	}

	list.Todos[todoID].Completed = !list.Todos[todoID].Completed

	err = s.UpdateTodoList(list, userID)
	if err != nil {
		return err
	}

	return nil
}
