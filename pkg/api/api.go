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
	// получить страницу page последних новостей при количестве quantity на страницу
	api.r.HandleFunc("/news/{page}/{quantity}", api.posts).Methods(http.MethodGet)
	// получить страницу page последних новостей при количестве quantity на страницу
	// содержащих keyword в заголовке
	api.r.HandleFunc("/filter/{page}/{quantity}/{keyword}", api.filter).Methods(http.MethodGet)
	// получить детальную новость n
	api.r.HandleFunc("/detailed/{n}", api.detailed).Methods(http.MethodGet)
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

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Получение публикаций по ключевому слову.
func (api *API) filter(w http.ResponseWriter, r *http.Request) {

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
	key := mux.Vars(r)["keyword"]

	type result struct {
		Count int
		Posts []storage.Post
	}
	res := result{}
	res.Posts, res.Count, err = api.db.Filter(n, q, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Filter error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Получение одной публикации.
func (api *API) detailed(w http.ResponseWriter, r *http.Request) {

	ns := mux.Vars(r)["n"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := storage.Post{}
	p, err = api.db.Post(n)
	if err != nil {
		http.Error(w, fmt.Sprintf("Post error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
