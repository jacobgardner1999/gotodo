package store

import "testing"

func TestAddUser(t *testing.T) {
	store := NewInMemoryStore()

	user := NewUser("0001", "Steve")
	store.CreateUser(user)

	got := store.users["0001"].Name
	want := "Steve"

	if got != want {
		t.Errorf("got %q want %q given, %q", got, want, "test")
	}
}

func TestAddDuplicateUserID(t *testing.T) {
	store := NewInMemoryStore()

	user := NewUser("0001", "Steve")
	user2 := NewUser("0001", "Stephen")

	store.CreateUser(user)
	err := store.CreateUser(user2)

	if err == nil {
		t.Fatal("expected an error")
	}

	want := "user with ID 0001 already exists"

	if err.Error() != want {
		t.Errorf("got %q want %q", err.Error(), want)
	}
}

func TestAddList(t *testing.T) {
	store := NewInMemoryStore()

	user := NewUser("0001", "Steve")
	store.CreateUser(user)

	list := NewTodoList("0001", "test list")
	store.AddTodoList(list, "0001")

	got := len(user.TodoLists)
	want := 1

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestAddDuplicateListID(t *testing.T) {
	store := NewInMemoryStore()

	user := NewUser("0001", "Steve")
	store.CreateUser(user)

	list := NewTodoList("0001", "test list")
	list2 := NewTodoList("0001", "duplicate list")
	store.AddTodoList(list, "0001")

	err := store.AddTodoList(list2, "0001")

	if err == nil {
		t.Fatal("expected an error")
	}

	want := "list with ID 0001 for user 0001 already exists"

	if err.Error() != want {
		t.Errorf("got %q want %q", err.Error(), want)
	}
}

func TestAddTodo(t *testing.T) {
	store := NewInMemoryStore()

	user := NewUser("0001", "Steve")
	store.CreateUser(user)

	list := NewTodoList("0001", "test list")
	store.AddTodoList(list, "0001")
	todo := Todo{"0001", "Make Todo App", false}
	store.AddTodo(todo, "0001", "0001")

	got := user.TodoLists["0001"].Todos["0001"].Title
	want := "Make Todo App"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
