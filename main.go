package main

import (
	"log/slog"
	"os"

	"github.com/papaya147/stonks/news_scraper"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := news_scraper.NewEconomicTimesScraper(logger, "https://economictimes.indiatimes.com/lazy_list_tech.cms", map[string]string{
		"information-tech": "78570530",
		"technology":       "78570561",
		"finance/banking":  "13358319",
		"energy/power":     "13358361",
	})

	s.ScrapeAndSave(1, 1, "economic-times-articles-formatted")

	// mcs := news_scraper.NewMoneyControlScraper(logger)
	// mcs.ScrapeAndSave(1, 1)
}
