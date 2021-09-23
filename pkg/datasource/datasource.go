package datasource

import (
	"fmt"
	"io"
	"net/http"
	"newssrv/pkg/storage"
	"time"
)

// Config - данные для обработчика.
type Source struct {
	Links         []string          `json:"rss"`
	RequestPeriod int               `json:"request_period"`
	PostChan      chan storage.Post `json:"-"`
	ErrorChan     chan error        `json:"-"`
	Parser        Interface         `json:"-"`
}

// Interface задаёт контракт на работу с обработчиком данных.
type Interface interface {
	Parse(io.Reader) ([]storage.Post, error) // запуск обработчика данных
}

//Run запускает опрос заданных адресов с заданным периодом
func (s *Source) Run() {
	for {
		for _, link := range s.Links {
			go func(link string) {
				resp, err := http.Get(link)
				if err != nil {
					s.ErrorChan <- fmt.Errorf("datasource.Run_http.Get error: %s", err)
					return
				}
				defer resp.Body.Close()
				posts, err := s.Parser.Parse(resp.Body)
				if err != nil {
					s.ErrorChan <- fmt.Errorf("datasource.Run_s.Parser.Parse error: %s", err)
					return
				}
				for _, p := range posts {
					s.PostChan <- p
				}
			}(link)
		}
		time.Sleep(time.Minute * time.Duration(s.RequestPeriod))
	}
}
