package lfd

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/snapp-incubator/seperno/internal"
	"github.com/snapp-incubator/seperno/pkg/options"
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
		"صد": 100, "یکصد": 100, "دویست": 200, "سیصد": 300, "چهارصد": 400, "چارصد": 400,
		"پانصد": 500, "پونصد": 500, "ششصد": 600, "شونصد": 600, "هفتصد": 700, "هشتصد": 800, "نهصد": 900,
		"هزار": 1000,
	}

	ordinalNumberMap = map[string]int64{
		"اول": 1, "دوم": 2, "سوم": 3,
	}

	ordinalSuffixes = []string{"مین", "ام", "وم", "م", "ین"}
	multipliers     = map[string]int64{"صد": 100, "هزار": 1000}

	// Compiled regexes
	lettersRegex = regexp.MustCompile(`^[\p{L}]+$`)              // Matches strings containing Unicode letters (for Persian word validation)
	tokenRegex   = regexp.MustCompile(`([\p{L}]+|[\p{N}]+|\s+)`) // Splits text into tokens: letter sequences, number sequences, or whitespace
)

type Token struct {
	Value      string
	StartIndex int
	EndIndex   int
}

type PersianNumberDetector struct{}

// DetectNumbers converts Persian number words to digits
func (f *PersianNumberDetector) DetectNumbers(text string) []DetectedNumber {
	if text == "" {
		return []DetectedNumber{}
	}

	normalizer := internal.NewNormalizer(options.DefaultOptions)
	normalizedCharacters := normalizer.NormalizeCharacters(text)
	preprocessed, addedSpacesPSumArray := preprocessConjunctions(string(normalizedCharacters))
	tokens := tokenizeWithPositions(preprocessed)

	return processTokensToNumbers(tokens, addedSpacesPSumArray)
}

// processTokensToNumbers processes tokens and converts detected numbers to DetectedNumber structs
func processTokensToNumbers(tokens []Token, addedSpacesPSumArray []int) []DetectedNumber {
	result := make([]DetectedNumber, 0)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if isWhitespace(token.Value) {
			continue
		}

		if _, numVal, startIdx, endIdx, isNumber := parseTokenWithPositions(token, tokens, &i); isNumber {
			result = append(result, DetectedNumber{
				Number:     numVal,
				StartIndex: startIdx - addedSpacesPSumArray[startIdx],
				EndIndex:   endIdx - addedSpacesPSumArray[endIdx],
			})
		}
	}
	return result
}

func tokenizeWithPositions(input string) []Token {
	matches := tokenRegex.FindAllStringSubmatch(input, -1)
	indexes := tokenRegex.FindAllStringIndex(input, -1)

	tokens := make([]Token, len(matches))
	for i, match := range matches {
		// Convert byte indices to rune indices
		startRuneIndex := utf8.RuneCountInString(input[:indexes[i][0]])
		endRuneIndex := utf8.RuneCountInString(input[:indexes[i][1]]) - 1

		tokens[i] = Token{
			Value:      match[0],
			StartIndex: startRuneIndex,
			EndIndex:   endRuneIndex,
		}
	}
	return tokens
}

// parseTokenWithPositions processes a single token and returns its numeric representation with positions
func parseTokenWithPositions(token Token, tokens []Token, index *int) (string, int64, int, int, bool) {
	trimmed := strings.TrimSpace(token.Value)
	if trimmed == "" {
		return token.Value, 0, 0, 0, false
	}

	// Handle existing digits
	if val, ok := parseDigits(trimmed); ok {
		return strconv.FormatInt(val, 10), val, token.StartIndex, token.EndIndex, true
	}

	// Handle ordinals and compound numbers
	if val, endIdx, ok := parseNumberWordWithPositions(token, tokens, index); ok {
		return strconv.FormatInt(val, 10), val, token.StartIndex, endIdx, true
	}

	return token.Value, 0, 0, 0, false
}

func parseNumberWordWithPositions(token Token, tokens []Token, index *int) (int64, int, bool) {
	word := token.Value

	// Try direct lookup
	if val, exists := persianNumberMap[word]; exists {
		return parseCompoundNumberWithPositions(val, token, tokens, index)
	}

	if val, exists := ordinalNumberMap[word]; exists {
		return val, token.EndIndex, true
	}

	// Try ordinal with suffix
	if val, ok := parseOrdinalWithSuffix(word); ok {
		return val, token.EndIndex, true
	}

	return 0, 0, false
}

func parseCompoundNumberWithPositions(initial int64, startToken Token, tokens []Token, index *int) (int64, int, bool) {
	total := initial
	pos := *index
	endIdx := startToken.EndIndex

	for {
		next := pos + 1

		// Skip whitespace
		for next < len(tokens) && isWhitespace(tokens[next].Value) {
			next++
		}

		if next >= len(tokens) {
			break
		}

		token := tokens[next]

		// Handle multipliers (هزار، صد)
		if mult, isMultiplier := multipliers[token.Value]; isMultiplier {
			total *= mult
			pos = next
			endIdx = token.EndIndex
			continue
		}

		// Handle separated hundreds (یک صد)
		if token.Value == "صد" {
			total *= 100
			pos = next
			endIdx = token.EndIndex
			continue
		}

		// Expect conjunction "و"
		if token.Value != "و" {
			break
		}

		next++ // Skip "و"

		// Skip whitespace after "و"
		for next < len(tokens) && isWhitespace(tokens[next].Value) {
			next++
		}

		if next >= len(tokens) || !lettersRegex.MatchString(tokens[next].Value) {
			break
		}

		// Parse next number
		nextVal, ok := parseNextNumber(tokens[next].Value)
		if !ok {
			break
		}

		total += nextVal
		pos = next
		endIdx = tokens[next].EndIndex
	}

	*index = pos
	return total, endIdx, true
}

// Keep the rest of the functions unchanged
func preprocessConjunctions(input string) (string, []int) {
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

	addedSpaces := make([]int, 0)
	ptr := 0
	resultRunes := []rune(result)
	inputRunes := []rune(input)
	for i := 0; i < len(resultRunes) && ptr < len(inputRunes); i++ {
		if resultRunes[i] != ' ' && inputRunes[ptr] != ' ' {
			ptr++
		} else if resultRunes[i] == ' ' && inputRunes[ptr] == ' ' {
			ptr++
		} else if resultRunes[i] == ' ' && inputRunes[ptr] != ' ' {
			addedSpaces = append(addedSpaces, i)
		}
	}

	psumArray := make([]int, len(resultRunes))
	for _, i := range addedSpaces {
		psumArray[i] = 1
	}
	for i := 1; i < len(psumArray); i++ {
		psumArray[i] += psumArray[i-1]
	}

	return result, psumArray
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

func isWhitespace(token string) bool {
	return strings.TrimSpace(token) == ""
}

func parseDigits(s string) (int64, bool) {
	if !isNumeric(s) {
		return 0, false
	}
	val, err := strconv.ParseInt(s, 10, 64)
	return val, err == nil
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
