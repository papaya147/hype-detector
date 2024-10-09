package main

import (
	"log/slog"
	"os"

	"github.com/papaya147/stonks/news_scraper"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	news_scraper.NewLiveMintScraper(
		logger.WithGroup("live mint scraper"),
		"https://www.livemint.com/listing/subsection/market~stock-market-news",
	).ScrapeAndSave(1, 50, "live-mint-articles-formatted")

	news_scraper.NewEconomicTimesScraper(
		logger.WithGroup("econimic times scraper"),
		"https://economictimes.indiatimes.com/lazy_list_tech.cms",
		map[string]string{
			"information-tech":   "78570530",
			"technology":         "78570561",
			"banking":            "13358319",
			"power":              "13358361",
			"auto":               "64829342",
			"electric-vehicles":  "81585238",
			"two/three-wheelers": "64829323",
			"finance":            "13358311",
			"hotels":             "13357036",
		},
	).ScrapeAndSave(1, 30, "economic-times-articles-formatted")

	news_scraper.NewMoneyControlScraper(
		logger.WithGroup("money control scraper"),
	).ScrapeAndSave(1, 30, "money-control-articles-formatted")
}
