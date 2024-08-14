package main

import (
	"github.com/papaya147/stonks/scraper"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := scraper.NewMoneyControlScraper(logger)
	articles := s.Scrape(1, 30)
	err := articles.Save("money-control.json")
	if err != nil {
		panic(err)
	}
}
