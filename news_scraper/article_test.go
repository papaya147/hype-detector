package news_scraper

import (
	"fmt"
	"testing"
	"time"
)

func TestNextMarketDay(t *testing.T) {
	day := time.Date(2024, time.August, 23, 17, 43, 0, 0, time.Now().Location())
	nextDay := nextMarketDay(day)
	fmt.Println(nextDay)
}
