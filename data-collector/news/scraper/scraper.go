package scraper

type Scraper interface {
	ScrapeAndSave(start, end int, folder string)
}
