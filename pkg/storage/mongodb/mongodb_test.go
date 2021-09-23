package mongodb

import (
	"fmt"
	"newssrv/pkg/storage"
	"os"
	"reflect"
	"testing"
	"time"
)

//Функциональный тест базы данных
func TestMongo(t *testing.T) {
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/dbtest?retryWrites=true&w=majority",
		pwd)
	db, err := New("dbtest", connstr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer func() {
		if err = db.DropDB(); err != nil {
			t.Error(err)
		}
	}()

	posts := []storage.Post{
		{ID: 1,
			Title:   "Вышел Microsoft Linux",
			Content: "Как сообщают непроверенные источники, новая ОС будет бесплатной.",
			PubTime: time.Now().Unix(),
			Link:    "https://github.com/microsoft/CBL-Mariner"},
		{ID: 2,
			Title:   "Инженеры Google не желают возвращаться в офисы",
			Content: "Инженеры Google не желают возвращаться в офисы, заявляя, что они не менее продуктивны на удалёнке.",
			PubTime: time.Now().Unix(),
			Link:    "https://habr.com/ru/news/t/568128/"}}
	for _, p := range posts {
		err = db.AddPost(p)
		if err != nil {
			t.Fatalf("AddPost error: %s", err)
		}
	}

	received, err := db.Posts()
	if err != nil {
		t.Fatalf("Posts error: %s", err)
	}
	if !reflect.DeepEqual(posts, received) {
		t.Errorf("received %v, wanted %v", received, posts)
	}

	p := storage.Post{
		ID:      2,
		Title:   "Инженеры Google не желают возвращаться в офисы",
		Content: "Инженеры Google не желают возвращаться в офисы, заявляя, что они не менее продуктивны на удалёнке.",
		PubTime: time.Now().Unix(),
		Link:    "https://habr.com/ru/news/t/568128"}

	err = db.UpdatePost(p)
	if err != nil {
		t.Fatalf("UpdatePost error: %s", err)
	}
	received, err = db.Posts()
	if err != nil {
		t.Fatalf("Posts error: %s", err)
	}
	if !reflect.DeepEqual(p, received[1]) {
		t.Errorf("received %v, wanted %v", received[1], p)
	}

	for _, p := range posts {
		err := db.DeletePost(p)
		if err != nil {
			t.Fatalf("DeletePost error: %s", err)
		}
	}
	received, err = db.Posts()
	if err != nil {
		t.Fatalf("Posts error: %s", err)
	}
	posts = []storage.Post{}
	if !reflect.DeepEqual(posts, received) {
		t.Errorf("received %v, wanted %v", received, posts)
	}
}
