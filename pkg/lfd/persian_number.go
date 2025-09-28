package lfd

import (
	"regexp"
	"strconv"
	"strings"
)

// DetectedNumber word mappings
var (
	persianNumberMap = map[string]int64{
		"صفر": 0, "یک": 1, "دو": 2, "سه": 3, "چهار": 4, "پنج": 5,
		"شش": 6, "هفت": 7, "هشت": 8, "نه": 9, "ده": 10,
		"یازده": 11, "دوازده": 12, "سیزده": 13, "چهارده": 14, "پانزده": 15,
		"شانزده": 16, "هفده": 17, "هجده": 18, "نوزده": 19,
		"بیست": 20, "سی": 30, "چهل": 40, "پنجاه": 50, "شصت": 60,
		"هفتاد": 70, "هشتاد": 80, "نود": 90,
		"صد": 100, "یکصد": 100, "دویست": 200, "سیصد": 300, "چهارصد": 400,
		"پانصد": 500, "ششصد": 600, "هفتصد": 700, "هشتصد": 800, "نهصد": 900,
		"هزار": 1000,
	}

	ordinalNumberMap = map[string]int64{
		"اول": 1, "دوم": 2, "سوم": 3,
	}

	ordinalSuffixes = []string{"مین", "ام", "وم", "م", "ین"}
	multipliers     = map[string]int64{"صد": 100, "هزار": 1000}

	// Compiled regexes
	lettersRegex = regexp.MustCompile(`^[\p{L}]+$`)
	tokenRegex   = regexp.MustCompile(`([\p{L}]+|[\p{N}]+|\s+)`)
)

type PersianNumberDetector struct{}

// DetectNumbers converts Persian number words to digits
func (f *PersianNumberDetector) DetectNumbers(text string) []DetectedNumber {
	if text == "" {
		return []DetectedNumber{}
	}

	preprocessed := preprocessConjunctions(text)
	tokens := tokenize(preprocessed)

	result := make([]DetectedNumber, 0)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if isWhitespace(token) {
			continue
		}

		if _, numVal, isNumber := parseToken(token, tokens, &i); isNumber {
			result = append(result, DetectedNumber{Number: numVal})
		}
	}
	return result
}

// preprocessConjunctions handles concatenated conjunctions
func preprocessConjunctions(input string) string {
	numberWords := getNumberWordList()
	pattern := `(` + strings.Join(numberWords, "|") + `)`

	replacements := []struct {
		pattern     string
		replacement string
	}{
		{pattern + `و` + pattern, "$1 و $2"},
		{pattern + `و`, "$1 و"},
		{`و` + pattern, "و $1"},
	}

	result := input
	for _, r := range replacements {
		re := regexp.MustCompile(r.pattern)
		result = re.ReplaceAllString(result, r.replacement)
	}

	return result
}

func getNumberWordList() []string {
	words := make([]string, 0, len(persianNumberMap)+len(ordinalNumberMap))

	for word := range persianNumberMap {
		words = append(words, regexp.QuoteMeta(word))
	}
	for word := range ordinalNumberMap {
		words = append(words, regexp.QuoteMeta(word))
	}

	return words
}

func tokenize(input string) []string {
	return tokenRegex.FindAllString(input, -1)
}

func isWhitespace(token string) bool {
	return strings.TrimSpace(token) == ""
}

// parseToken processes a single token and returns its numeric representation if applicable
func parseToken(token string, tokens []string, index *int) (string, int64, bool) {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return token, 0, false
	}

	// Handle existing digits
	if val, ok := parseDigits(trimmed); ok {
		return strconv.FormatInt(val, 10), val, true
	}

	// Handle ordinals and compound numbers
	if val, ok := parseNumberWord(trimmed, tokens, index); ok {
		return strconv.FormatInt(val, 10), val, true
	}

	return token, 0, false
}

func parseDigits(s string) (int64, bool) {
	if !isNumeric(s) {
		return 0, false
	}

	normalized := normalizeDigits(s)
	val, err := strconv.ParseInt(normalized, 10, 64)
	return val, err == nil
}

func parseNumberWord(word string, tokens []string, index *int) (int64, bool) {
	// Try direct lookup
	if val, exists := persianNumberMap[word]; exists {
		return parseCompoundNumber(val, tokens, index)
	}

	if val, exists := ordinalNumberMap[word]; exists {
		return val, true
	}

	// Try ordinal with suffix
	if val, ok := parseOrdinalWithSuffix(word); ok {
		return val, true
	}

	return 0, false
}

func parseOrdinalWithSuffix(word string) (int64, bool) {
	for _, suffix := range ordinalSuffixes {
		if strings.HasSuffix(word, suffix) {
			base := strings.TrimSuffix(word, suffix)

			if val, exists := persianNumberMap[base]; exists {
				return val, true
			}
			if val, exists := ordinalNumberMap[base]; exists {
				return val, true
			}
		}
	}
	return 0, false
}

func parseCompoundNumber(initial int64, tokens []string, index *int) (int64, bool) {
	total := initial
	pos := *index

	for {
		next := pos + 1

		// Skip whitespace
		for next < len(tokens) && isWhitespace(tokens[next]) {
			next++
		}

		if next >= len(tokens) {
			break
		}

		token := tokens[next]

		// Handle multipliers (هزار، صد)
		if mult, isMultiplier := multipliers[token]; isMultiplier {
			total *= mult
			pos = next
			next++
			continue
		}

		// Handle separated hundreds (یک صد)
		if token == "صد" {
			total *= 100
			pos = next
			next++
			continue
		}

		// Expect conjunction "و"
		if token != "و" {
			break
		}

		next++ // Skip "و"

		// Skip whitespace after "و"
		for next < len(tokens) && isWhitespace(tokens[next]) {
			next++
		}

		if next >= len(tokens) || !lettersRegex.MatchString(tokens[next]) {
			break
		}

		// Parse next number
		nextVal, ok := parseNextNumber(tokens[next])
		if !ok {
			break
		}

		total += nextVal
		pos = next
	}

	*index = pos
	return total, true
}

func parseNextNumber(word string) (int64, bool) {
	// Direct number word
	if val, exists := persianNumberMap[word]; exists {
		return val, true
	}

	// Irregular ordinal
	if val, exists := ordinalNumberMap[word]; exists {
		return val, true
	}

	// Ordinal with suffix
	return parseOrdinalWithSuffix(word)
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= '۰' && r <= '۹')) {
			return false
		}
	}
	return true
}

func normalizeDigits(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if r >= '۰' && r <= '۹' {
			result.WriteRune('0' + (r - '۰'))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
