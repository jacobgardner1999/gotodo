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

func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	var list store.TodoList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Println("Create List - Error Decoding ", err)
		InternalServerErrorHandler(w, r)
		return
	}
	log.Println(list)

	matches := ListRe.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		log.Println("Create List - Not enough arguments")
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.UpdateTodoList(list, matches[1]); err != nil {
		log.Println("Create List - ", err)
		InternalServerErrorHandler(w, r)
		return
	}
	log.Println("Update List - Success")
	w.WriteHeader(http.StatusOK)
}

func (h *ListHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	matches := ListRe.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		log.Println("Get Lists - Not enough arguments")
		InternalServerErrorHandler(w, r)
		return
	}

	lists, err := h.store.GetTodoLists(matches[1])
	if err != nil {
		log.Println("Get Lists - ", err)
		InternalServerErrorHandler(w, r)
		return
	}

	byteValue, err := json.MarshalIndent(lists, "", "  ")
	if err != nil {
		log.Println("Get  Lists - Marshal error ", err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Println("Get Lists - Success")

	w.WriteHeader(http.StatusOK)
	w.Write(byteValue)
}

func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	matches := ListReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 3 {
		log.Println("Get List - Not enough arguments")
		InternalServerErrorHandler(w, r)
		return
	}

	list, err := h.store.GetTodoList(matches[1], matches[2])
	if err != nil {
		log.Println("Get List - ", err)
		InternalServerErrorHandler(w, r)
		return
	}

	byteValue, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Println("Get List - Marshal Error ", err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Println("Get List - Success")

	w.WriteHeader(http.StatusOK)
	w.Write(byteValue)
}

func (h *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	matches := ListReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 3 {
		log.Println("Delete List - Not enough arguments")
		InternalServerErrorHandler(w, r)
		return
	}

	err := h.store.DeleteTodoList(matches[1], matches[2])
	if err != nil {
		log.Println("Delete list - ", err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Println("Delete List - Success")
	w.WriteHeader(http.StatusOK)
}

var (
	ListRe       = regexp.MustCompile(`^/lists/([^/]+)$`)
	ListReWithID = regexp.MustCompile(`^/lists/([^/]+)/([^/])$`)
)

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && ListRe.MatchString(r.URL.Path):
		h.UpdateList(w, r)
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
