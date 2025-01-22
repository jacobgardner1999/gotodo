package store

type Store interface {
	CreateUser(username string) (id string, e error)
	GetUser(id string) (User, error)
	GetTodoList(userID string, listID string) (TodoList, error)
	GetTodoLists(userID string) (map[string]*TodoList, error)
	UpdateTodoList(list TodoList, userID string) error
	DeleteTodoList(userID string, listID string) error
	AddTodo(todo Todo, listID string, userID string) error
	ToggleTodo(userID string, listID string, todoID string) error
}

type User struct {
	ID        string
	Name      string
	TodoLists map[string]*TodoList
}

type TodoList struct {
	ID    string
	Name  string
	Todos map[string]*Todo
}

type Todo struct {
	ID        string
	Title     string
	Completed bool
}

func NewUser(id, name string) User {
	return User{
		ID:        id,
		Name:      name,
		TodoLists: make(map[string]*TodoList),
	}
}

func NewTodoList(id, name string) TodoList {
	return TodoList{
		ID:    id,
		Name:  name,
		Todos: make(map[string]*Todo),
	}
}

func (t *Todo) Toggle() {
	t.Completed = !t.Completed
}
