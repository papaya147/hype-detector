package main

import (
	"log/slog"
	"os"

	"github.com/papaya147/stonks/news_scraper"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	scrapers := map[string]news_scraper.Scraper{}

	scrapers["economic-times-articles-formatted"] = news_scraper.NewEconomicTimesScraper(
		logger.WithGroup("econimic times scraper"),
		"https://economictimes.indiatimes.com/lazy_list_tech.cms",
		map[string]string{
			"information-tech": "78570530",
			"technology":       "78570561",
			"finance/banking":  "13358319",
			"energy/power":     "13358361",
		})

	scrapers["money-control-articles-formatted"] = news_scraper.NewMoneyControlScraper(logger.WithGroup("money control scraper"))

	for folder, scraper := range scrapers {
		scraper.ScrapeAndSave(1, 30, folder)
	}
}
