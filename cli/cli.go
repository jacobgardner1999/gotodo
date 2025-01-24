package main

import (
	"ToDo/store"
	"fmt"
	"maps"
	"os"
	"regexp"
	"slices"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state      string
	page       string
	store      store.Store
	user       *store.User
	toDoLists  []*store.TodoList
	listID     string
	list       *store.TodoList
	toDoList   []*store.Todo
	input      string
	cursor     int
	loginError string
}

func InitialModel() model {
	apiStore := store.NewApiStore("8080", "../data/")
	return model{
		state:      "userInput",
		page:       "login",
		store:      apiStore,
		user:       &store.User{},
		toDoLists:  []*store.TodoList{},
		listID:     "",
		list:       nil,
		toDoList:   []*store.Todo{},
		input:      "",
		cursor:     0,
		loginError: "",
	}
}

type Msg string

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	re := regexp.MustCompile(`[0-9]+$`)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "main":
			switch msg.String() {
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "j", "down":
				switch m.page {
				case "lists":
					if m.cursor < len(m.toDoLists)-1 {
						m.cursor++
					}
				case "todos":
					if m.cursor < len(m.toDoList)-1 {
						m.cursor++
					}
				}
			case "d":
				switch m.page {
				case "lists":
					m.store.DeleteTodoList(m.user.ID, m.toDoLists[m.cursor].ID)
					todos, _ := m.store.GetTodoLists(m.user.ID)
					m.toDoLists = slices.Collect(maps.Values(todos))
					m.cursor = 0
				}
			case "h", "left":
				lists, _ := m.store.GetTodoLists(m.user.ID)
				m.toDoLists = slices.Collect(maps.Values(lists))
				switch m.page {
				case "lists":
					m.state = "userInput"
					m.page = "login"
					m.cursor = 0
				case "todos":
					m.page = "lists"
					m.cursor = 0
				}
			case "a":
				m.state = "userInput"
			case "enter", "l", "right":
				switch m.page {
				case "lists":
					m.toDoList = slices.Collect(maps.Values(m.toDoLists[m.cursor].Todos))
					m.listID = m.toDoLists[m.cursor].ID
					m.list = m.toDoLists[m.cursor]
					m.page = "todos"
					m.cursor = 0
				case "todos":
					m.store.ToggleTodo(m.user.ID, m.list.ID, m.toDoList[m.cursor].ID)
					list, _ := m.store.GetTodoList(m.user.ID, m.listID)
					m.toDoList = slices.Collect(maps.Values(list.Todos))
				case "addUser":
					m.state = "userInput"
					m.page = "login"
					m.cursor = 0
					m.input = ""
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		case "userInput":
			switch msg.String() {
			case "enter":
				switch m.page {
				case "login":
					user, err := m.store.GetUser(m.input)
					if err != nil {
						m.loginError = "failed to get user with ID " + m.input
						break
					}
					m.user = &user
					todos, _ := m.store.GetTodoLists(m.user.ID)
					m.toDoLists = slices.Collect(maps.Values(todos))
					m.state = "main"
					m.page = "lists"
					m.input = ""
					m.cursor = 0
				case "lists":
					m.store.UpdateTodoList(store.NewTodoList(strconv.Itoa(len(m.toDoLists)), m.input), m.user.ID)
					todos, _ := m.store.GetTodoLists(m.user.ID)
					m.toDoLists = slices.Collect(maps.Values(todos))
					m.input = ""
					m.state = "main"
					m.cursor = 0
				case "todos":
					todo := store.Todo{ID: strconv.Itoa(len(m.toDoList)), Title: m.input, Completed: false}
					m.store.AddTodo(todo, m.list.ID, m.user.ID)
					todos, _ := m.store.GetTodoList(m.user.ID, m.listID)
					m.toDoList = slices.Collect(maps.Values(todos.Todos))
					m.input = ""
					m.state = "main"
					m.cursor = 0
				case "addUser":
					id, _ := m.store.CreateUser(m.input)
					m.input = id
					m.state = "main"
					m.cursor = 0
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case "a":
				if m.page == "login" {
					m.input = ""
					m.page = "addUser"
				} else {
					m.input += msg.String()
				}
			case "ctrl+c":
				return m, tea.Quit
			default:
				if m.page == "login" && !re.MatchString(msg.String()) {
					return m, nil
				}
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	lineBreak := "\n--------------------\n"

	s := "Woah! Another Todo App!" + lineBreak

	switch m.state {
	case "main":
		switch m.page {
		case "lists":
			s += "Todo Lists:" + lineBreak
			for i, list := range m.toDoLists {
				cursor := " "
				if i == m.cursor {
					cursor = ">"
				}
				s += fmt.Sprintf("%s %s\n", cursor, list.Name)
			}
			s += lineBreak
			s += "Press Enter to select, q to quit, a to add list, d to delete list"
			return s
		case "todos":
			s += "Todo list: " + m.list.Name
			s += lineBreak
			if len(m.toDoList) == 0 {
				s += "--list is empty, press a to add a todo--"
			}
			for i, todo := range m.toDoList {
				cursor := " "
				if i == m.cursor {
					cursor = ">"
				}
				check := " "
				if todo.Completed {
					check = "X"
				}
				s += fmt.Sprintf("%s [%s] %s\n", cursor, check, todo.Title)
			}
			s += lineBreak
			s += "Press Enter to complete task, q to quit, a to add todo"
			return s
		case "addUser":
			s += "User added! Your ID is: " + m.input + lineBreak + "\n (Press Enter to continue)"
			return s
		}
	case "userInput":
		switch m.page {
		case "login":
			s += "Enter your User ID to log in: " + m.input + lineBreak + m.loginError + "\n (Press Enter to continue, ctrl+c to quit, a to add a user)\n"
			return s
		case "todos":
			s += "What do you need to do? " + m.input + lineBreak + "\n (Press Enter to continue)"
			return s
		case "lists":
			s += "Enter the name of your new list: " + m.input + lineBreak + "\n (Press Enter to continue)"
			return s
		case "addUser":
			s += "Enter your name: " + m.input + lineBreak + "\n (Press Enter to continue)"
			return s
		}
	}
	return ""
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		os.Exit(1)
	}
}
