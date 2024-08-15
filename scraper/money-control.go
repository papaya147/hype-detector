package scraper

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/papaya147/parallelize"
	"log/slog"
	"strings"
	"time"
)

type MoneyControlScraper struct {
	logger  *slog.Logger
	baseUrl string
}

func NewMoneyControlScraper(logger *slog.Logger) *MoneyControlScraper {
	return &MoneyControlScraper{
		logger:  logger,
		baseUrl: "https://www.moneycontrol.com/news/tags/companies/news",
	}
}

func (mcs *MoneyControlScraper) articleLinks(start, end int) []string {
	var articleLinks []string
	c := colly.NewCollector()

	c.OnHTML("ul#cagetory li.clearfix h2 a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		articleLinks = append(articleLinks, link)
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

	return articleLinks
}

func (mcs *MoneyControlScraper) Scrape(start, end int) Articles {
	links := mcs.articleLinks(start, end)
	var articles []*Article

	batchSize := 5

	for i := 0; i < len(links); i += batchSize {
		end = i + batchSize
		if end > len(links) {
			end = len(links)
		}

		batchLinks := links[i:end]

		g := parallelize.NewGroup()
		channels := make([]parallelize.WithOutputWithArgsChannels[*Article], len(batchLinks))
		for j, link := range batchLinks {
			ch := parallelize.AddWithOutputWithArgs(g, mcs.scrape, context.Background(), link)
			channels[j] = ch
		}

		g.Execute()

		for _, ch := range channels {
			article, err := ch.Read()
			if err != nil {
				mcs.logger.Error(err.Error())
			}

			if article != nil {
				articles = append(articles, article)
			}
		}
	}

	return articles
}

func (mcs *MoneyControlScraper) scrape(_ context.Context, link string) (*Article, error) {
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
		return nil, err
	}

	c.Wait()

	if content == "" {
		return nil, nil
	}

	return NewArticle(link, title, desc, content, timestamp), nil
}
