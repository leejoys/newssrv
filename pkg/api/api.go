package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newssrv/pkg/storage"
	"strconv"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера GoNews
type API struct {
	db storage.Interface
	r  *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{page}/{quantity}", api.posts).Methods(http.MethodGet)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

// Получение публикаций.
func (api *API) posts(w http.ResponseWriter, r *http.Request) {

	ns := mux.Vars(r)["page"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qs := mux.Vars(r)["quantity"]
	q, err := strconv.Atoi(qs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type result struct {
		Count int
		Posts []storage.Post
	}
	res := result{}
	res.Posts, res.Count, err = api.db.PostsN(n, q)
	if err != nil {
		http.Error(w, fmt.Sprintf("PostsN error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	//uncomment for mongo
	// postsWithIDs := []storage.Post{}
	// for i, post := range posts {
	// 	post.ID = i + 1
	// 	postsWithIDs = append(postsWithIDs, post)
	// }

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
