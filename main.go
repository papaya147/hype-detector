package main

import (
	"github.com/papaya147/stonks/scraper"
	"log/slog"
	"os"
)

func saveArticles() scraper.Articles {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := scraper.NewMoneyControlScraper(logger)

	articles := s.Scrape(1, 30)

	err := articles.Save("money-control-articles")
	if err != nil {
		panic(err)
	}

	return articles
}

func formatArticles(articles scraper.Articles) {
	articles.FormatContent()

	err := articles.Save("money-control-articles-formatted")
	if err != nil {
		panic(err)
	}
}

func main() {
	articles := saveArticles()
	formatArticles(articles)
}
