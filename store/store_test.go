package store

import "testing"

func TestAddUser(t *testing.T) {
	store := NewInMemoryStore()

	user := User{ID: "0001", Name: "Steve"}
	store.CreateUser(user)

	got := store.users["0001"].Name
	want := "Steve"

	if got != want {
		t.Errorf("got %q want %q given, %q", got, want, "test")
	}
}

func TestAddDuplicateUserID(t *testing.T) {
	store := NewInMemoryStore()

	user := User{ID: "001", Name: "Steve"}
	user2 := User{ID: "001", Name: "Stephen"}

	store.CreateUser(user)
	err := store.CreateUser(user2)

	if err == nil {
		t.Fatal("expected an error")
	}

	want := "user with ID 001 already exists"

	if err.Error() != want {
		t.Errorf("got %q want %q", err.Error(), want)
	}
}
