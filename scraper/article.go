package scraper

import (
	"encoding/json"
	"os"
	"time"
)

type Article struct {
	Url         string    `json:"url,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewArticle(url, title, description, content string, timestamp time.Time) *Article {
	return &Article{
		Url:         url,
		Title:       title,
		Description: description,
		Content:     content,
		Timestamp:   timestamp,
	}
}

func (a *Article) String() string {
	return a.Timestamp.String() + " - " + a.Title + " - " + a.Description + " - " + a.Content
}

type Articles []*Article

func (a Articles) Save(fileName string) error {
	j, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, j, 0644)
}
