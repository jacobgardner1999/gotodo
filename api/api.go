package main

import (
	"ToDo/store"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func main() {
	store, _ := store.NewJsonStore("../data/")
	listHandler := NewListHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &HomeHandler{})
	mux.Handle("/lists/", listHandler)

	log.Fatalln("ListenAndServe: ", http.ListenAndServe(":8080", mux))
}

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

type ListHandler struct {
	store store.Store
}

func NewListHandler(s store.Store) *ListHandler {
	return &ListHandler{
		store: s,
	}
}

func (h *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list store.TodoList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Println("Create List - Error Decoding")
		InternalServerErrorHandler(w, r)
		return
	}

	matches := ListRe.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		log.Println("Create List - Not enough arguments")
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.AddTodoList(list, matches[1]); err != nil {
		log.Println("Create List - Error adding list")
		InternalServerErrorHandler(w, r)
		return
	}

	log.Println("Create List - Added \n", list)
	w.WriteHeader(http.StatusOK)
}

func (h *ListHandler) GetLists(w http.ResponseWriter, r *http.Request) {}

func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {}

func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {}

func (h *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {}

var (
	ListRe       = regexp.MustCompile(`^/lists/([^/]+)$`)
	ListReWithID = regexp.MustCompile(`^/lists/([^/]+)/([^/])$`)
)

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && ListRe.MatchString(r.URL.Path):
		h.CreateList(w, r)
		return
	case r.Method == http.MethodGet && ListRe.MatchString(r.URL.Path):
		h.GetLists(w, r)
		return
	case r.Method == http.MethodGet && ListReWithID.MatchString(r.URL.Path):
		h.GetList(w, r)
		return
	case r.Method == http.MethodPut && ListReWithID.MatchString(r.URL.Path):
		h.UpdateList(w, r)
		return
	case r.Method == http.MethodDelete && ListReWithID.MatchString(r.URL.Path):
		h.DeleteList(w, r)
		return
	default:
		return
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
