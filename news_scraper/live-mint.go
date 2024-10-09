package news_scraper

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LiveMintScraper struct {
	logger          *slog.Logger
	baseUrl         string
	articleLinkChan chan string
	articleChan     chan *Article
}

func NewLiveMintScraper(logger *slog.Logger, baseUrl string) *LiveMintScraper {
	return &LiveMintScraper{
		logger:          logger,
		baseUrl:         baseUrl,
		articleLinkChan: make(chan string, 50),
		articleChan:     make(chan *Article, 50),
	}
}

func (lms *LiveMintScraper) fetchArticleLinks(start, end int) {
	c := colly.NewCollector(
		colly.UserAgent(""),
	)

	c.OnHTML("h2.headline a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		lms.articleLinkChan <- "https://www.livemint.com" + link
	})

	c.OnError(func(r *colly.Response, err error) {
		lms.logger.Error(err.Error())
	})

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	for i := start; i <= end; i++ {
		url := fmt.Sprintf("%s/%d", lms.baseUrl, i)
		lms.logger.Info(fmt.Sprintf("visiting page %d", i))

		err := c.Visit(url)
		if err != nil {
			lms.logger.Error(fmt.Sprintf("error visiting page %d: %s", i, err))
		}

		c.Wait()
		lms.logger.Info(fmt.Sprintf("page %d completed", i))
	}
}

func (lms *LiveMintScraper) ScrapeAndSave(start, end int, folder string) {
	go func() {
		defer close(lms.articleLinkChan)
		lms.fetchArticleLinks(start, end)
	}()

	go func() {
		defer close(lms.articleChan)
		wg := sync.WaitGroup{}
		for link := range lms.articleLinkChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lms.scrape(context.Background(), link)
			}()
		}
		wg.Wait()
	}()

	for article := range lms.articleChan {
		article.FormatContent()
		err := article.Save(folder)
		if err != nil {
			panic(err)
		}
	}
}

func (lms *LiveMintScraper) scrape(_ context.Context, link string) {
	lms.logger.Info(fmt.Sprintf("scraping link %s", link))
	c := colly.NewCollector(
		colly.UserAgent(""),
	)

	c.WithTransport(&http.Transport{
		IdleConnTimeout: 60 * time.Second,
	})

	var title, desc, content string
	var timestamp time.Time
	var err error

	c.OnHTML("h1#article-0", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	c.OnHTML("h2.storyPage_summary__Ge5SX", func(e *colly.HTMLElement) {
		desc = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.storyPage_date__JS9qJ span", func(e *colly.HTMLElement) {
		timeStr := strings.TrimSpace(e.Text)
		timestamp, err = time.Parse("2 Jan 2006, 3:04 PM MST", timeStr)
		if err != nil {
			lms.logger.Error(err.Error())
		}
	})

	c.OnHTML("div.storyParagraph p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if !strings.HasPrefix(text, "Disclaimer") {
			content += strings.TrimSpace(text)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		lms.logger.Error(err.Error())
	})

	err = c.Visit(link)
	if err != nil {
		lms.logger.Error(err.Error())
	}

	c.Wait()

	if timestamp.Unix() > 0 {
		lms.articleChan <- NewArticle(link, title, desc, content, timestamp)
	}
}
