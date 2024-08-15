package scraper

import (
	"fmt"
	"regexp"
	"testing"
)

func TestArticle_FormatContent(t *testing.T) {
	percentReg := regexp.MustCompile(`(\d+)\s*(%|percent|per\s*cent)`)
	content := "65 percent"
	parts := percentReg.FindStringSubmatch(content)
	fmt.Println(NumberToToken(parts[1], true))
}
