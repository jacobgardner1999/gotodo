package main

import (
	store "ToDo/store"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state     string
	width     int
	height    int
	store     *store.InMemoryStore
	userID    string
	user      *store.User
	toDoLists []*store.TodoList
	list      *store.TodoList
	toDoList  []*store.Todo
	title     string
	cursor    int
}

func InitialModel() model {
	return model{
		state:     "login",
		width:     0,
		height:    0,
		store:     store.NewInMemoryStore(),
		userID:    "",
		user:      &store.User{},
		toDoLists: []*store.TodoList{},
		list:      nil,
		toDoList:  []*store.Todo{},
		title:     "",
		cursor:    0,
	}
}

type Msg string

func (m model) Init() tea.Cmd {
	user := store.User{ID: "0001", Name: "Jacob", TodoLists: make(map[string]*store.TodoList)}

	m.store.CreateUser(user)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case "login":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				user, err := m.store.GetUser(m.userID)
				if err != nil {
					fmt.Println("Error obtaining user with ID ", m.userID)
					break
				}
				m.user = &user
				m.toDoLists = slices.Collect(maps.Values(m.user.TodoLists))
				m.state = "main"
				m.cursor = 0
			case "backspace":
				m.userID = m.userID[:len(m.userID)-1]
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.userID += msg.String()
			}
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
		}
	case "main":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "j", "down":
				if m.cursor < len(m.toDoLists)-1 {
					m.cursor++
				}
			case "h", "left":
				m.state = "login"
			case "a":
				m.state = "addList"
			case "enter", "l", "right":
				m.toDoList = slices.Collect(maps.Values(m.toDoLists[m.cursor].Todos))
				m.list = m.toDoLists[m.cursor]
				m.state = "list"
				m.cursor = 0
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	case "addList":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.store.AddTodoList(store.NewTodoList(strconv.Itoa(len(m.toDoLists)), m.title), m.user.ID)
				m.toDoLists = slices.Collect(maps.Values(m.user.TodoLists))
				m.title = ""
				m.state = "main"
				m.cursor = 0
			case "backspace":
				m.title = m.title[:len(m.title)-1]
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.title += msg.String()
			}
		}
	case "addTodo":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				todo := store.Todo{ID: strconv.Itoa(len(m.toDoList)), Title: m.title, Completed: false}
				m.store.AddTodo(todo, m.list.ID, m.user.ID)
				m.toDoList = slices.Collect(maps.Values(m.list.Todos))
				m.title = ""
				m.state = "list"
				m.cursor = 0
			case "backspace":
				m.title = m.title[:len(m.title)-1]
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.title += msg.String()
			}
		}
	case "list":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "j", "down":
				if m.cursor < len(m.toDoList)-1 {
					m.cursor++
				}
			case "h", "left":
				m.state = "main"
			case "a":
				m.state = "addTodo"
			case "enter", "l", "right":
				m.store.CompleteTodo(m.user.ID, m.list.ID, m.toDoList[m.cursor].ID)
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	lineBreak := "\n--------------------\n"

	if m.width == 0 {
		return "loading..."
	}

	switch m.state {
	case "login":
		return "Enter your User ID to log in: " + m.userID + lineBreak + "\n (Press Enter to continue) \n"
	case "main":
		s := "Todo Lists:" + lineBreak
		for i, list := range m.toDoLists {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, list.Name)
		}
		s += lineBreak
		s += "Press Enter to select, q to quit, a to add list"
		return s
	case "list":
		s := "Todo list: " + m.list.Name
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
	case "addTodo":
		return "What do you need to do? " + m.title + lineBreak + "\n (Press Enter to continue)"
	case "addList":
		return "Enter the name of your new list: " + m.title + lineBreak + "\n (Press Enter to continue)"
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
