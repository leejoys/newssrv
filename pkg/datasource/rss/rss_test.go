package rss

import (
	"newssrv/pkg/storage"
	"reflect"
	"strings"
	"testing"
)

//юнит-тест для Parser_Parse
func TestParser_Parse(t *testing.T) {
	body := strings.NewReader("<rss><channel><title>GOSAMPLES - Learn Golang programming by example</title><link>https://gosamples.dev/</link><description>Learn Golang programming by example. </description><generator>Hugo -- gohugo.io</generator><language>en-us</language><image><url>https://gosamples.dev/apple-touch-icon.png</url><title>GOSAMPLES - Learn Golang programming by example</title><link>https://gosamples.dev/</link></image><lastBuildDate>Mon, 25 Jan 2021 00:00:00 +0000</lastBuildDate><item><title>title1</title><link>link1</link><pubDate>Fri, 23 Jul 2021 00:00:00 +0000</pubDate><guid>https://gosamples.dev/convert-int-to-string/</guid><description>description1</description></item><item><title>title2</title><link>link2</link><pubDate>Tue, 01 Jun 2021 00:00:00 +0000</pubDate><guid>https://gosamples.dev/write-csv/</guid><description>description2</description></item></channel></rss>")

	p := New()
	posts, err := p.Parse(body)
	if err != nil {
		t.Fatalf("Parse error:%s", err)
	}
	want := []storage.Post{{
		ID:      0,
		Title:   "title1",
		Content: "description1",
		PubDate: "Fri, 23 Jul 2021 00:00:00 +0000",
		PubTime: 0,
		Link:    "link1"}, {
		ID:      0,
		Title:   "title2",
		Content: "description2",
		PubDate: "Tue, 01 Jun 2021 00:00:00 +0000",
		PubTime: 0,
		Link:    "link2"}}
	if !reflect.DeepEqual(posts, want) {
		t.Errorf("posts=%v, want %v", posts, want)
	}

}
