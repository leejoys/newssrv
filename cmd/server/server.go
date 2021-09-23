package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"newssrv/pkg/api"
	"newssrv/pkg/datasource"
	"newssrv/pkg/datasource/rss"
	"newssrv/pkg/storage"
	"newssrv/pkg/storage/mongodb"
	"os"
	"time"
)

// Сервер newssrv.
type server struct {
	ds  *datasource.Source
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера
	srv := server{}

	// Создаем источник данных
	cByte, err := os.ReadFile("./aconfig.json")
	if err != nil {
		log.Fatalf("main ioutil.ReadFile error: %s", err)
	}
	srv.ds = &datasource.Source{}
	err = json.Unmarshal(cByte, srv.ds)
	if err != nil {
		log.Fatalf("main json.Unmarshal error: %s", err)
	}
	srv.ds.PostChan = make(chan storage.Post)
	srv.ds.ErrorChan = make(chan error)
	p := rss.New()
	srv.ds.Parser = p

	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/newssrv?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("newssrv", connstr)
	if err != nil {
		log.Fatalf("mongo.New error: %s", err)
	}

	// Инициализируем хранилище сервера БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	go srv.poster()
	go srv.logger()
	go srv.ds.Run()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	log.Fatal(http.ListenAndServe(":8080", srv.api.Router()))
}

//обрабатывает ответы из каналов с постами
func (s *server) poster() {
	for post := range s.ds.PostChan {

		t, err := time.Parse(time.RFC1123, post.PubDate)
		if err != nil {
			s.ds.ErrorChan <- fmt.Errorf("poster time.Parse error: %s", err)
		}
		post.PubTime = t.Unix()
		err = s.db.AddPost(post)
		if err != nil {
			s.ds.ErrorChan <- err
		}
	}
}

//обрабатывает ответы из каналов с ошибками
func (s *server) logger() {
	for err := range s.ds.ErrorChan {
		if err == mongodb.ErrorDuplicatePost {
			log.Println(err)
			continue
		}
		log.Fatalln(err)
	}
}
