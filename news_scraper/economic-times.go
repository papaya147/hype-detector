package news_scraper

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type EconomicTimesScraper struct {
	logger          *slog.Logger
	baseUrl         string
	msids           map[string]string
	articleLinkChan chan string
	articleChan     chan *Article
}

func NewEconomicTimesScraper(logger *slog.Logger, baseUrl string, msids map[string]string) *EconomicTimesScraper {
	return &EconomicTimesScraper{
		logger:          logger,
		baseUrl:         baseUrl,
		msids:           msids,
		articleLinkChan: make(chan string, 10),
		articleChan:     make(chan *Article, 10),
	}
}

func (ets *EconomicTimesScraper) fetchArticleLinks(start, end int) {
	c := colly.NewCollector()

	c.OnHTML(".story-box .desc h4 a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		ets.articleLinkChan <- link
	})

	c.OnError(func(r *colly.Response, err error) {
		ets.logger.Error(err.Error())
	})

	for fieldPage, msid := range ets.msids {
		for i := start; i <= end; i++ {
			url := fmt.Sprintf("%s?msid=%s&curpg=%d", ets.baseUrl, msid, i)
			ets.logger.Info(fmt.Sprintf("visiting %s (msid %s) page %d", fieldPage, msid, i))

			err := c.Visit(url)
			if err != nil {
				ets.logger.Error(fmt.Sprintf("error visiting page: %s", err))
			}

			ets.logger.Info(fmt.Sprintf("msid %s page %d completed", msid, i))
		}
	}

	c.Wait()
}

func (ets *EconomicTimesScraper) ScrapeAndSave(start, end int, folder string) {
	go func() {
		defer close(ets.articleLinkChan)
		ets.fetchArticleLinks(start, end)
	}()

	go func() {
		defer close(ets.articleChan)
		for link := range ets.articleLinkChan {
			go func() {
				ets.scrape(context.Background(), link)
			}()
		}
	}()

	for article := range ets.articleChan {
		article.FormatContent()
		err := article.Save(folder)
		if err != nil {
			panic(err)
		}
	}
}

var styleTagRegexp = regexp.MustCompile("<style>[^<]*</style>")
var aTagRegexp = regexp.MustCompile("<a[^<]*>|</a[^<]*>")
var divTagWithNoContentRegexp = regexp.MustCompile("<div[^<>]*></div>")
var extraContentRegexp = regexp.MustCompile("<div.*")
var selfClosedTagRegexp = regexp.MustCompile("<[^>]*>")

func (ets *EconomicTimesScraper) scrape(_ context.Context, link string) {
	ets.logger.Info(fmt.Sprintf("scraping link %s", link))
	c := colly.NewCollector()

	var title, desc, content string
	var timestamp time.Time
	var err error

	c.OnHTML("time.jsdtTime", func(e *colly.HTMLElement) {
		epochMillisString := strings.TrimSpace(e.Attr("data-dt"))
		epochMillis, err := strconv.Atoi(epochMillisString)
		if err != nil {
			ets.logger.Error(err.Error())
		}
		timestamp = time.Unix(int64(epochMillis/1e3), 0)
	})

	c.OnHTML("h1.artTitle", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	c.OnHTML("h2.summary", func(e *colly.HTMLElement) {
		desc = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.artText", func(e *colly.HTMLElement) {
		innerHtml, _ := e.DOM.Html()
		innerHtml = styleTagRegexp.ReplaceAllString(innerHtml, "")
		innerHtml = aTagRegexp.ReplaceAllString(innerHtml, "")

		for i := 0; i < 3; i++ {
			innerHtml = divTagWithNoContentRegexp.ReplaceAllString(innerHtml, "")
		}

		innerHtml = extraContentRegexp.ReplaceAllString(innerHtml, "")
		innerHtml = selfClosedTagRegexp.ReplaceAllString(innerHtml, "")

		content = html.UnescapeString(innerHtml)
	})

	err = c.Visit(link)
	if err != nil {
		ets.logger.Error(err.Error())
	}

	c.Wait()

	if content != "" {
		ets.articleChan <- NewArticle(link, title, desc, content, timestamp)
	}
}
