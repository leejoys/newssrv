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
	Post(int) (Post, error)                       //Post - получение одной публикации
	PostsN(int, int) ([]Post, int, error)         // получение n-ной страницы публикаций
	Filter(int, int, string) ([]Post, int, error) // получение n-ной страницы публикаций с ключевым словом
	AddPost(Post) error                           // создание новой публикации
	UpdatePost(Post) error                        // обновление публикации
	DeletePost(Post) error                        // удаление публикации по ID
	Close()                                       // освобождение ресурса
}
