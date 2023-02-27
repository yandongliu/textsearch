package entities

import "time"

type Index map[string][]int

const (
	DocTypeHN = "hacker_news"
	DocTypeGN = "google_news"
)

// const (
// 	HackerNews = 0
// 	GoogleNews = 1
// )

type Document struct {
	Type  string
	Title string
	URL   string
	Text  string
	ID    int
}

type HNDocument struct {
	By    string
	Id    int
	Score int
	Time  int64
	Title string
	Type  string
	Url   string
	Date  time.Time
}
