package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/papaya147/stonks/news_scraper"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dataFolder := "data/news"

	news_scraper.NewLiveMintScraper(
		logger.WithGroup("live mint scraper"),
		"https://www.livemint.com/listing/subsection/market~stock-market-news",
	).ScrapeAndSave(1, 50, fmt.Sprintf("%s/live-mint-articles-formatted", dataFolder))

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
	).ScrapeAndSave(1, 30, fmt.Sprintf("%s/economic-times-articles-formatted", dataFolder))

	news_scraper.NewMoneyControlScraper(
		logger.WithGroup("money control scraper"),
	).ScrapeAndSave(1, 30, fmt.Sprintf("%s/money-control-articles-formatted", dataFolder))
}
