package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"newssrv/pkg/storage"
)

type RSS struct{}

//создает объект парсера RSS с заданными параметрами
func New() *RSS {
	return &RSS{}
}

//читает RSS, парсит item в []storage.Post
func (s *RSS) Parse(body io.Reader) ([]storage.Post, error) {

	decoder := xml.NewDecoder(body)
	posts := []storage.Post{}
	// Чтение item по частям
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("RSS.Parse_decoder.Token error: %s", err)
		}
		//выбор токена по типу
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "item" {
				// Декодирование элемента в структуру
				var post storage.Post
				err = decoder.DecodeElement(&post, &tp)
				if err != nil {
					return nil, fmt.Errorf("RSS.Parse_decoder.DecodeElement error: %s", err)
				}
				posts = append(posts, post)
			}
		}
	}
	return posts, nil
}
