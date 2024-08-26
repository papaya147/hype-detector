package news_scraper

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

var magnitudeMap = map[int]string{
	10:                "TEN",
	100:               "HUNDRED",
	1_000:             "THOUSAND",
	10_000:            "TEN_THOUSAND",
	100_000:           "HUNDRED_THOUSAND",
	1_000_000:         "MILLION",
	10_000_000:        "TEN_MILLION",
	100_000_000:       "HUNDRED_MILLION",
	1_000_000_000:     "BILLION",
	10_000_000_000:    "TEN_BILLION",
	100_000_000_000:   "HUNDRED_BILLION",
	1_000_000_000_000: "TRILLION",
}

var magnitudeMapKeys []int

func reverse(slice []int) {
	start := 0
	end := len(slice) - 1

	for start < end {
		slice[start], slice[end] = slice[end], slice[start] // Swap elements
		start++
		end--
	}
}

func init() {
	magnitudeMapKeys = make([]int, len(magnitudeMap))
	i := 0
	for k := range magnitudeMap {
		magnitudeMapKeys[i] = k
		i++
	}
	sort.Ints(magnitudeMapKeys)
	reverse(magnitudeMapKeys)
}

var numberMap = map[int]string{
	0: "ZERO",
	1: "ONE",
	2: "TWO",
	3: "THREE",
	4: "FOUR",
	5: "FIVE",
	6: "SIX",
	7: "SEVEN",
	8: "EIGHT",
	9: "NINE",
}

func NumberToToken(number string, isPercentage bool) (string, error) {
	number = strings.ReplaceAll(number, ",", "")
	floatNum, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return "", err
	}
	num := int(floatNum)

	absNum := int(math.Abs(float64(num)))
	token := ""

	if absNum != num {
		token += "NEGATIVE_"
	}

	// getting MSB
	n := num
	for n >= 10 {
		n /= 10
	}
	token += numberMap[n] + "_"

	// getting magnitude
	for _, mag := range magnitudeMapKeys {
		if num >= mag {
			token += magnitudeMap[mag]
			if isPercentage {
				token += "_"
			}
			break
		}
	}

	if isPercentage {
		token += "PERCENT"
	}

	return "<" + strings.TrimSpace(token) + ">", nil
}
