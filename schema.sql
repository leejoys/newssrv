DROP TABLE IF EXISTS posts;

CREATE TABLE IF NOT EXISTS posts(
id SERIAL PRIMARY KEY,
title TEXT UNIQUE NOT NULL,
content TEXT NOT NULL,
pubdate TEXT NOT NULL,
pubtime BIGINT NOT NULL,
link TEXT NOT NULL
);

INSERT INTO posts (id, title, content, pubdate, pubtime, link) 
VALUES (0, 'Статья', 'Содержание статьи', 0, 0, 'google.com');
