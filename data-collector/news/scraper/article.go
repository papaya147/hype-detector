package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func holidays() []time.Time {
	return []time.Time{
		// republic day
		time.Date(2024, time.January, 26, 0, 0, 0, 0, time.Local),

		// independence day
		time.Date(2024, time.August, 15, 0, 0, 0, 0, time.Local),

		// gandhi's bday
		time.Date(2024, time.October, 2, 0, 0, 0, 0, time.Local),

		// diwali
		time.Date(2024, time.November, 1, 0, 0, 0, 0, time.Local),

		// holi
		time.Date(2024, time.March, 25, 0, 0, 0, 0, time.Local),

		// dusshera
		time.Date(2024, time.October, 12, 0, 0, 0, 0, time.Local),
	}
}

func nextMarketDay(day time.Time) time.Time {
	// before market hours, move time to 09:15
	if day.Hour() < 9 || (day.Hour() == 9 && day.Minute() < 15) {
		nextDay := time.Date(day.Year(), day.Month(), day.Day(), 9, 15, 0, 0, day.Location())
		return nextMarketDay(nextDay)
	}

	// after market hours, move time to next day
	if day.Hour() > 15 || (day.Hour() == 15 && day.Minute() > 30) {
		nextDay := day.AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 15, 0, 0, nextDay.Location())
		return nextMarketDay(nextDay)
	}

	// weekend, move time to next monday
	if day.Weekday() == time.Saturday {
		nextDay := day.AddDate(0, 0, 2)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 15, 0, 0, nextDay.Location())
		return nextMarketDay(nextDay)
	} else if day.Weekday() == time.Sunday {
		nextDay := day.AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 15, 0, 0, nextDay.Location())
		return nextMarketDay(nextDay)
	}

	// holiday, move time to next day
	for _, holiday := range holidays() {
		if day.Day() == holiday.Day() && day.Month() == holiday.Month() {
			nextDay := day.AddDate(0, 0, 1)
			nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 15, 0, 0, nextDay.Location())
			return nextMarketDay(nextDay)
		}
	}

	// market day as normal
	return day
}

type Article struct {
	Url             string    `json:"url,omitempty"`
	Title           string    `json:"title,omitempty"`
	Description     string    `json:"description,omitempty"`
	Content         string    `json:"content,omitempty"`
	CleanedContent  string    `json:"cleaned_content,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
	MarketTimestamp time.Time `json:"market_timestamp"`
	OffMarketHours  bool      `json:"off_market_hours"`
}

func NewArticle(url, title, description, content string, timestamp time.Time) *Article {
	marketTimestamp := nextMarketDay(timestamp)

	return &Article{
		Url:             url,
		Title:           title,
		Description:     description,
		Content:         content,
		Timestamp:       timestamp,
		MarketTimestamp: marketTimestamp,
		OffMarketHours: timestamp.Day() != marketTimestamp.Day() ||
			timestamp.Hour() != marketTimestamp.Hour() ||
			timestamp.Minute() != marketTimestamp.Minute(),
	}
}

func (a *Article) FormatContent() {
	// remove everything but number, % symbols, alphabets and spaces
	reg := regexp.MustCompile(`[^a-zA-Z0-9%,.' ]+`)
	cleanedText := reg.ReplaceAllString(a.Content, " ")
	a.CleanedContent = strings.ToLower(cleanedText)

	// replacing number percentages
	percentReg := regexp.MustCompile(`\s+(\d[\d,]*\.?[\d,]*)\s*(%|percent|per\s*cent)`)
	a.CleanedContent = percentReg.ReplaceAllStringFunc(a.CleanedContent, func(s string) string {
		parts := percentReg.FindStringSubmatch(s)
		token, err := NumberToToken(parts[1], true)
		if err != nil {
			panic(err)
		}
		return " " + token + " "
	})

	// replacing numbers
	numberReg := regexp.MustCompile(`\s+(\d[\d,]*\.?[\d,]*)`)
	a.CleanedContent = numberReg.ReplaceAllStringFunc(a.CleanedContent, func(s string) string {
		parts := numberReg.FindStringSubmatch(s)
		token, err := NumberToToken(parts[1], false)
		if err != nil {
			panic(err)
		}
		return " " + token + " "
	})

	a.CleanedContent = strings.ReplaceAll(a.CleanedContent, ",", " ")
	a.CleanedContent = strings.ReplaceAll(a.CleanedContent, ".", " ")
	a.CleanedContent = strings.ReplaceAll(a.CleanedContent, "'", "")
	a.CleanedContent = strings.TrimSpace(a.CleanedContent)

	re := regexp.MustCompile(`\s+`)
	a.CleanedContent = re.ReplaceAllString(a.CleanedContent, " ")
}

func (a *Article) String() string {
	return a.Timestamp.String() + " - " + a.Title + " - " + a.Description + " - " + a.Content
}

var nonLetterRegexp = regexp.MustCompile(`[^\w]+`)

func (a *Article) safeFileName() string {
	safeTitle := nonLetterRegexp.ReplaceAllString(a.Title, "-")
	return fmt.Sprintf("%d-%s", a.Timestamp.Unix(), safeTitle)
}

func (a *Article) Save(folder string) error {
	j, err := json.Marshal(a)
	if err != nil {
		return err
	}

	f := fmt.Sprintf("%s/%s.json", folder, a.safeFileName())
	_ = os.Mkdir(folder, os.ModePerm)
	return os.WriteFile(f, j, 0644)
}

type Articles []*Article

func (a Articles) Save(folder string) error {
	for _, article := range a {
		err := article.Save(folder)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadArticles(fileName string) (Articles, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var articles Articles
	err = json.NewDecoder(file).Decode(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a Articles) FormatContent() {
	for _, article := range a {
		article.FormatContent()
	}
}
