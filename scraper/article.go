package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Article struct {
	Url            string    `json:"url,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	Content        string    `json:"content,omitempty"`
	CleanedContent string    `json:"cleaned_content,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
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

func (a *Article) FormatContent() {
	// remove everything but number, % symbols, alphabets and spaces
	reg := regexp.MustCompile(`[^a-zA-Z0-9%,. ]+`)
	cleanedText := reg.ReplaceAllString(a.Content, " ")
	cleanedText = strings.ToLower(cleanedText)

	// replacing number percentages
	percentReg := regexp.MustCompile(`\s+(\d[\d,]+)\s*(%|percent|per\s*cent)`)
	percentReplacedText := percentReg.ReplaceAllStringFunc(cleanedText, func(s string) string {
		parts := percentReg.FindStringSubmatch(s)
		token, err := NumberToToken(parts[1], true)
		if err != nil {
			panic(err)
		}
		return " " + token
	})

	// replacing numbers
	numberReg := regexp.MustCompile(`\s+(\d[\d,]+)`)
	numberReplacedText := numberReg.ReplaceAllStringFunc(percentReplacedText, func(s string) string {
		parts := numberReg.FindStringSubmatch(s)
		token, err := NumberToToken(parts[1], false)
		if err != nil {
			panic(err)
		}
		return " " + token
	})

	a.CleanedContent = strings.TrimSpace(numberReplacedText)
}

func (a *Article) String() string {
	return a.Timestamp.String() + " - " + a.Title + " - " + a.Description + " - " + a.Content
}

func (a *Article) Save(folder string) error {
	j, err := json.Marshal(a)
	if err != nil {
		return err
	}

	f := fmt.Sprintf("%s/%d %s.json", folder, a.Timestamp.Unix(), a.Title)
	_ = os.Mkdir(folder, os.ModePerm)
	return os.WriteFile(f, j, 0644)
}

type Articles []*Article

func (a Articles) Save(folder string) error {
	for _, article := range a {
		err := article.Save(folder)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadArticles(fileName string) (Articles, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var articles Articles
	err = json.NewDecoder(file).Decode(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a Articles) FormatContent() {
	for _, article := range a {
		article.FormatContent()
	}
}
