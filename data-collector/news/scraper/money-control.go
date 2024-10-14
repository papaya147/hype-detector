package scraper

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type MoneyControlScraper struct {
	logger          *slog.Logger
	baseUrl         string
	articleLinkChan chan string
	articleChan     chan *Article
}

func NewMoneyControlScraper(logger *slog.Logger) *MoneyControlScraper {
	return &MoneyControlScraper{
		logger:          logger,
		baseUrl:         "https://www.moneycontrol.com/news/tags/companies/news",
		articleLinkChan: make(chan string, 10),
		articleChan:     make(chan *Article, 10),
	}
}

func (mcs *MoneyControlScraper) fetchArticleLinks(start, end int) {
	c := colly.NewCollector()

	c.OnHTML("ul#cagetory li.clearfix h2 a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		mcs.articleLinkChan <- link
	})

	c.OnError(func(r *colly.Response, err error) {
		mcs.logger.Error(err.Error())
	})

	for i := start; i <= end; i++ {
		url := fmt.Sprintf("%s/page-%d/", mcs.baseUrl, i)
		mcs.logger.Info(fmt.Sprintf("visiting page %d", i))

		err := c.Visit(url)
		if err != nil {
			mcs.logger.Error(fmt.Sprintf("error visiting page %d: %s", i, err))
		}

		c.Wait()
		mcs.logger.Info(fmt.Sprintf("page %d completed", i))
	}
}

func (mcs *MoneyControlScraper) ScrapeAndSave(start, end int, folder string) {
	go func() {
		defer close(mcs.articleLinkChan)
		mcs.fetchArticleLinks(start, end)
	}()

	go func() {
		defer close(mcs.articleChan)
		wg := sync.WaitGroup{}
		for link := range mcs.articleLinkChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mcs.scrape(context.Background(), link)
			}()
		}
		wg.Wait()
	}()

	for article := range mcs.articleChan {
		article.FormatContent()
		err := article.Save(folder)
		if err != nil {
			panic(err)
		}
	}
}

func (mcs *MoneyControlScraper) scrape(_ context.Context, link string) {
	mcs.logger.Info(fmt.Sprintf("scraping link %s", link))
	c := colly.NewCollector()

	var title, desc, content string
	var timestamp time.Time
	var err error

	c.OnHTML("div.page_left_wrapper", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.ChildText("h1"))
		desc = strings.TrimSpace(e.ChildText("h2"))

		timeDiv := e.DOM.Find("div.article_schedule")
		date := timeDiv.Find("span").Text()
		timeText := strings.TrimSpace(strings.Split(timeDiv.Text(), "/")[1])
		timestamp, err = time.Parse("January 02, 2006 15:04 MST", fmt.Sprintf("%s %s", date, timeText))
		if err != nil {
			mcs.logger.Error(err.Error())
		}

		contentDiv := e.DOM.Find("div#contentdata")
		if contentDiv.Length() == 0 {
			mcs.logger.Info(fmt.Sprintf("article %s contains no content, skipping", link))
			return
		}

		contentDiv.Find("p").Each(func(i int, selection *goquery.Selection) {
			attr, _ := selection.Attr("class")
			if attr == "" {
				content += strings.TrimSpace(selection.Text()) + " "
			}
		})
	})

	err = c.Visit(link)
	if err != nil {
		mcs.logger.Error(err.Error())
	}

	c.Wait()

	if content != "" {
		mcs.articleChan <- NewArticle(link, title, desc, content, timestamp)
	}
}
