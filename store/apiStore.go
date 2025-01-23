package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ApiStore struct {
	Users      map[string]*User
	serverPort string
	storePath  string
}

func NewApiStore(serverPort, storePath string) ApiStore {
	return ApiStore{
		Users:      make(map[string]*User),
		serverPort: serverPort,
		storePath:  storePath,
	}
}

func (s ApiStore) getUsersFromJson() (map[string]*User, error) {
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

func (s ApiStore) CreateUser(username string) (id string, e error) {
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

func (s ApiStore) GetUser(id string) (User, error) {
	users, err := s.getUsersFromJson()
	if err != nil {
		return User{}, fmt.Errorf("%s", err)
	}

	s.Users = users
	if _, exists := s.Users[id]; !exists {
		return User{}, fmt.Errorf("no user found with ID %s", id)
	}

	return *s.Users[id], nil
}

func (s ApiStore) GetTodoList(userID string, listID string) (TodoList, error) {
	requestURL := fmt.Sprintf("http://localhost:%s/lists/%s/%s", s.serverPort, userID, listID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return TodoList{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TodoList{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return TodoList{}, err
	}

	var list TodoList

	if err = json.Unmarshal(resBody, list); err != nil {
		return TodoList{}, err
	}

	return list, nil
}

func (s ApiStore) GetTodoLists(userID string) (map[string]*TodoList, error) {
	lists := make(map[string]*TodoList)

	requestURL := fmt.Sprintf("http://localhost:%s/lists/%s", s.serverPort, userID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Panic("1: ", err)
		return lists, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("2: ", err)
		return lists, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic("3: ", err)
		return lists, err
	}

	if err = json.Unmarshal(resBody, &lists); err != nil {
		log.Panic("4: ", err, resBody, res.Status)
		return lists, err
	}

	return lists, nil
}

func (s ApiStore) UpdateTodoList(list TodoList, userID string) error {
	todos, err := s.GetTodoLists(userID)
	if err != nil {
		return err
	}

	todos[list.ID] = &list
	byteValue, err := json.MarshalIndent(todos, "", "  ")

	bodyReader := bytes.NewReader(byteValue)

	requestURL := fmt.Sprintf("http://localhost:%s/lists/%s", s.serverPort, userID)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	}

	return nil
}

func (s ApiStore) DeleteTodoList(userID string, listID string) error {
	return nil
}

func (s ApiStore) AddTodo(todo Todo, listID string, userID string) error {
	return nil
}

func (s ApiStore) ToggleTodo(userID string, listID string, todoID string) error {
	return nil
}
