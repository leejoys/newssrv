package api

//mongotest need rafactoring
// func TestAPI_posts(t *testing.T) {
// 	// Создаём чистый объект API для теста.
// 	pwd := os.Getenv("Cloud0pass")
// 	connstr := fmt.Sprintf(
// 		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/dbtest?retryWrites=true&w=majority",
// 		pwd)
// 	dbase, err := mongodb.New("dbtest", connstr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer dbase.Close()
// 	defer func() {
// 		if err = dbase.DropDB(); err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	posts := []storage.Post{
// 		{ID: 3,
// 			Title:   "Вышел Microsoft Linux",
// 			Content: "Как сообщают непроверенные источники, новая ОС будет бесплатной.",
// 			PubTime: time.Now().Unix(),
// 			Link:    "https://github.com/microsoft/CBL-Mariner"},
// 		{ID: 2,
// 			Title:   "Инженеры Google не желают возвращаться в офисы",
// 			Content: "Инженеры Google не желают возвращаться в офисы, заявляя, что они не менее продуктивны на удалёнке.",
// 			PubTime: time.Now().Unix(),
// 			Link:    "https://habr.com/ru/news/t/568128/"},
// 		{ID: 1,
// 			Title:   "Название",
// 			Content: "Контент.",
// 			PubTime: time.Now().Unix(),
// 			Link:    "https://google.com"}}
// 	for _, p := range posts {
// 		err = dbase.AddPost(p)
// 		if err != nil {
// 			t.Fatalf("AddPost error: %s", err)
// 		}
// 	}

// 	api := New(dbase)
// 	// Создаём HTTP-запрос.
// 	req := httptest.NewRequest(http.MethodGet, "/news/2", nil)
// 	// Создаём объект для записи ответа обработчика.
// 	rr := httptest.NewRecorder()
// 	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
// 	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
// 	api.r.ServeHTTP(rr, req)

// 	// Проверяем код ответа.
// 	if !(rr.Code == http.StatusOK) {
// 		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
// 	}
// 	// Читаем тело ответа.
// 	b, err := ioutil.ReadAll(rr.Body)
// 	if err != nil {
// 		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
// 	}
// 	// Раскодируем JSON в массив новостей.
// 	var data []storage.Post
// 	err = json.Unmarshal(b, &data)
// 	if err != nil {
// 		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
// 	}
// 	// Проверяем, что в массиве ровно две новости.
// 	wantArr := []storage.Post{posts[2], posts[1]}
// 	if !reflect.DeepEqual(data, wantArr) {
// 		t.Fatalf("получено %v, ожидалось %v", data, wantArr)
// 	}
// }
