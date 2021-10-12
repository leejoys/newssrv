package pgdb

import (
	"context"
	"errors"
	"newssrv/pkg/storage"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrorDuplicatePost error = errors.New("SQLSTATE 23505")

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

//New - Конструктор объекта хранилища.
func New(connstr string) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

//Close - освобождение ресурса
func (s *Store) Close() {
	s.db.Close()
}

//Post - получение одной публикации
func (s *Store) Post(n int) (storage.Post, error) {
	p := storage.Post{}
	err := s.db.QueryRow(context.Background(),
		`SELECT 
	posts.id, 
	posts.title, 
	posts.content, 
	posts.pubdate, 
	posts.pubtime,
	posts.link
	FROM posts
	WHERE id=$1;`, n).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.PubDate,
		&p.PubTime,
		&p.Link,
	)

	if err != nil {
		return storage.Post{}, err
	}

	return p, err
}

//PostsN - получение n-ной страницы публикаций при q публикаций на страницу
func (s *Store) PostsN(n, q int) ([]storage.Post, int, error) {
	count := 0
	err := s.db.QueryRow(context.Background(),
		`SELECT count(*) FROM posts;`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	o := n * q
	rows, err := s.db.Query(context.Background(),
		`SELECT 
	posts.id, 
	posts.title, 
	posts.content, 
	posts.pubdate, 
	posts.pubtime,
	posts.link
	FROM posts
	OFFSET $1
	LIMIT $2;`, o, q)

	if err != nil {
		return nil, 0, err
	}

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubDate,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	return posts, count, rows.Err()
}

//Filter - получение n-ной страницы публикаций с ключевым словом key при q публикаций на страницу
func (s *Store) Filter(n, q int, key string) ([]storage.Post, int, error) {
	count := 0
	err := s.db.QueryRow(context.Background(),
		`SELECT count(*) FROM posts WHERE title ILIKE '%'||$1||'%';`, key).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	o := n * q
	rows, err := s.db.Query(context.Background(),
		`SELECT 
	posts.id, 
	posts.title, 
	posts.content, 
	posts.pubdate, 
	posts.pubtime,
	posts.link
	FROM posts
	WHERE title ILIKE '%'||$3||'%'
	OFFSET $1
	LIMIT $2;`, o, q, key)

	if err != nil {
		return nil, 0, err
	}

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubDate,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	return posts, count, rows.Err()
}

//AddPost - создание новой публикации
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO posts (
		title, 
		content, 
		pubdate, 
		pubtime,
		link) 
	VALUES ($1,$2,$3,$4,$5);`,
		p.Title,
		p.Content,
		p.PubDate,
		p.PubTime,
		p.Link)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				err = ErrorDuplicatePost
			}
		}
	}
	return err
}

//TODO ErrorDuplicatePost
//UpdatePost - обновление по id значения title,  pubdate, pubtime, и link
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE posts 
	SET title=$2,
	content=$3,
	pubdate=$4,
	pubtime=$5,
	link=$6
	WHERE id=$1;`,
		p.ID,
		p.Title,
		p.Content,
		p.PubDate,
		p.PubTime,
		p.Link)
	return err
}

//DeletePost - удаляет пост по id
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM posts 
	WHERE id=$1;`, p.ID)
	return err
}
