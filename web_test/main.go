package main

import (
	"fmt"
	"net/http"

	"github.com/dairaga/webkit"
)

type category struct {
	ID     uint64 `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Parent uint64 `json:"parent,omitempty"`
}

var categories = make(map[uint64]*category)

func middle(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return r, true
}

func list(w http.ResponseWriter, r *http.Request) webkit.Result {
	lst := make([]*category, 0, len(categories))

	for _, v := range categories {
		lst = append(lst, v)
	}

	return webkit.OK(lst)
}

func find(w http.ResponseWriter, r *http.Request, id uint64) webkit.Result {

	category, ok := categories[id]
	if !ok {
		return webkit.NotFound(nil)

	}
	return webkit.OK(category)
}

func add(w http.ResponseWriter, r *http.Request, category *category) webkit.Result {

	id := uint64(len(categories) + 1)
	category.ID = id

	categories[id] = category

	w.Header().Add("Location", fmt.Sprintf("/test/categories/%d", id))
	return webkit.OK(nil)
}

func update(w http.ResponseWriter, r *http.Request, id uint64, category *category) webkit.Result {

	category.ID = id
	categories[id] = category

	return webkit.Custom(http.StatusNoContent, nil)
}

func del(w http.ResponseWriter, r *http.Request, id uint64) {

	_, ok := categories[id]
	if !ok {
		w.WriteHeader(404)
		return
	}
	delete(categories, id)
	w.WriteHeader(http.StatusNoContent)

}

func main() {
	fmt.Println("test")
	webkit.Use("/test", middle)

	webkit.Handle("/test/categories", list).Methods("GET")
	webkit.Handle("/test/categories", add).Methods("POST")
	webkit.Handle("/test/categories/{id:[0-9]+}", find, "id").Methods("GET")
	webkit.Handle("/test/categories/{id:[0-9]+}", update, "id").Methods("PUT")
	webkit.Handle("/test/categories/{id:[0-9]+}", del, "id").Methods("DELETE")

	fmt.Println("start ...")
	http.ListenAndServe(":8080", webkit.Router())

}
