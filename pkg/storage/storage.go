package storage

// Post - публикация.
type Post struct {
	ID      int    `xml:"-" json:"ID"`                // номер записи
	Title   string `xml:"title" json:"Title"`         // заголовок публикации
	Content string `xml:"description" json:"Content"` // содержание публикации
	PubDate string `xml:"pubDate" json:"-"`           // время публикации из RSS
	PubTime int64  `xml:"-" json:"PubTime"`           //время публикации для БД и фронта
	Link    string `xml:"link" json:"Link"`           // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error)          // получение всех публикаций
	PostsN(int, int) ([]Post, error) // получение n-ной страницы публикаций
	AddPost(Post) error              // создание новой публикации
	UpdatePost(Post) error           // обновление публикации
	DeletePost(Post) error           // удаление публикации по ID
	Close()                          // освобождение ресурса
	//	DropDB() error              //удаление БД (при работе с монго)
}
