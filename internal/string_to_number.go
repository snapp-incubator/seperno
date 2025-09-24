package internal

import (
	"regexp"
	"strconv"
	"strings"
)

// persianNumberMap maps Persian number words to their digit equivalents.
var persianNumberMap = map[string]string{
	"صفر":    "0",
	"یک":     "1",
	"دو":     "2",
	"سه":     "3",
	"چهار":   "4",
	"پنج":    "5",
	"شش":     "6",
	"هفت":    "7",
	"هشت":    "8",
	"نه":     "9",
	"ده":     "10",
	"یازده":  "11",
	"دوازده": "12",
	"سیزده":  "13",
	"چهارده": "14",
	"پانزده": "15",
	"شانزده": "16",
	"هفده":   "17",
	"هجده":   "18",
	"نوزده":  "19",
	"بیست":   "20",
	"سی":     "30",
	"چهل":    "40",
	"پنجاه":  "50",
	"شصت":    "60",
	"هفتاد":  "70",
	"هشتاد":  "80",
	"نود":    "90",
	"صد":     "100",
	"دویست":  "200",
	"سیصد":   "300",
	"چهارصد": "400",
	"پانصد":  "500",
	"ششصد":   "600",
	"هفتصد":  "700",
	"هشتصد":  "800",
	"نهصد":   "900",
	"هزار":   "1000",
}

// ordinalNumberMap maps irregular Persian ordinal words to digits.
var ordinalNumberMap = map[string]string{
	"اول": "1",
	"دوم": "2",
	"سوم": "3",
}

// ordinalSuffixes lists common Persian ordinal suffixes.
var ordinalSuffixes = []string{"م", "ام", "وم", "مین"}

// lettersRe matches strings containing only Persian letters.
var lettersRe = regexp.MustCompile(`^\p{L}+$`)

// ConvertWordsToIntFa converts Persian number words to digits and returns the converted string
// along with a slice of integers found in the input. It handles number words (e.g., "بیست"),
// existing digits (e.g., "22"), and compound numbers with concatenated "و" (e.g., "هزاروهفتصد").
func ConvertWordsToIntFa(input string) (string, []int64) {
	// Preprocess to handle concatenated "و" adjacent to number words.
	input = preprocessConjunctions(input)
	words := splitWithDelimiters(input)
	var result []string
	var numbers []int64

	for i := 0; i < len(words); i++ {
		word := words[i]
		if numStr, isNumber, numVal := convertNumberWord(word, words, &i); isNumber {
			result = append(result, numStr)
			numbers = append(numbers, numVal)
		} else {
			result = append(result, word)
		}
	}

	return strings.Join(result, ""), numbers
}

// preprocessConjunctions inserts spaces around "و" when it’s adjacent to one or two number words.
func preprocessConjunctions(input string) string {
	// List all number words for regex.
	var numberWords []string
	for word := range persianNumberMap {
		numberWords = append(numberWords, regexp.QuoteMeta(word))
	}
	for word := range ordinalNumberMap {
		numberWords = append(numberWords, regexp.QuoteMeta(word))
	}
	numberWordPattern := `(` + strings.Join(numberWords, "|") + `)`

	// Patterns to match:
	// 1. when <numberWord>و<numberWord> → <numberWord> و <numberWord>
	// 2. when  <numberWord>و → <numberWord> و
	// 3. when و<numberWord> → و <numberWord>
	patterns := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(numberWordPattern + `و` + numberWordPattern), "$1 و $2"},
		{regexp.MustCompile(numberWordPattern + `و`), "$1 و"},
		{regexp.MustCompile(`و` + numberWordPattern), "و $1"},
	}

	// Apply all replacements.
	result := input
	for _, p := range patterns {
		result = p.re.ReplaceAllString(result, p.repl)
	}
	return result
}

// splitWithDelimiters splits the input into tokens: Persian words, digits, and spaces.
// Assumes punctuation is pre-removed and special spaces are normalized to regular spaces.
func splitWithDelimiters(input string) []string {
	re := regexp.MustCompile(`([\p{L}]+|[\p{N}]+|\s+)`)
	return re.FindAllString(input, -1)
}

// convertNumberWord converts a single or compound number word (e.g., "سی و پنج") or digit
// to its numeric equivalent. It returns the converted string, a boolean indicating if it’s a number,
// and the numeric value. Updates the index for compound numbers.
func convertNumberWord(word string, words []string, index *int) (string, bool, int64) {
	trimmed := strings.TrimSpace(word)
	if trimmed == "" {
		return word, false, 0
	}

	// Case 1: Existing digits (English or Persian).
	if isNumeric(trimmed) {
		normalized := normalizeDigits(trimmed)
		val, _ := strconv.ParseInt(normalized, 10, 64)
		return normalized, true, val
	}

	// Case 2: Irregular ordinals (e.g., "اول", "دوم").
	if val, ok := ordinalNumberMap[trimmed]; ok {
		numVal, _ := strconv.ParseInt(val, 10, 64)
		return val, true, numVal
	}

	// Case 3: Ordinals with suffixes (e.g., "چهارمین", "پنجم").
	for _, suffix := range ordinalSuffixes {
		if strings.HasSuffix(trimmed, suffix) {
			base := strings.TrimSuffix(trimmed, suffix)
			if val, ok := persianNumberMap[base]; ok {
				numVal, _ := strconv.ParseInt(val, 10, 64)
				return val, true, numVal
			}
			if val, ok := ordinalNumberMap[base]; ok {
				numVal, _ := strconv.ParseInt(val, 10, 64)
				return val, true, numVal
			}
		}
	}

	// Case 4: Cardinal or compound numbers (e.g., "سی و پنج").
	return parseCompoundNumber(trimmed, words, index)
}

// isNumeric checks if a string contains only digits (English 0-9 or Persian ۰-۹).
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if (r >= '0' && r <= '9') || (r >= '۰' && r <= '۹') {
			continue
		}
		return false
	}
	return true
}

// normalizeDigits converts Persian digits to English digits.
func normalizeDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '۰' && r <= '۹' {
			b.WriteRune('0' + (r - '۰'))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// parseCompoundNumber parses single or compound number words (e.g., "سی و پنج") into digits.
// Returns the converted string, a boolean indicating success, and the numeric value.
func parseCompoundNumber(word string, words []string, index *int) (string, bool, int64) {
	// Check if the first word is a known number word.
	first, ok := persianNumberMap[word]
	if !ok {
		return "", false, 0
	}
	total, _ := strconv.ParseInt(first, 10, 64)
	last := *index

	// Process compound numbers (e.g., "سی و پنج").
	for {
		j := last + 1
		// Skip spaces.
		for j < len(words) && strings.TrimSpace(words[j]) == "" {
			j++
		}
		// Expect "و" (and) for compound numbers.
		if j >= len(words) || words[j] != "و" {
			break
		}
		j++ // Move past "و".
		// Skip spaces after "و".
		for j < len(words) && strings.TrimSpace(words[j]) == "" {
			j++
		}
		// Ensure the next token is a valid word.
		if j >= len(words) || !lettersRe.MatchString(words[j]) {
			break
		}
		// The next word must be a number word.
		if nextStr, ok := persianNumberMap[words[j]]; ok {
			v, _ := strconv.ParseInt(nextStr, 10, 64)
			total += v
			last = j // Update to the last consumed token.
			continue
		}
		// Not a number word; stop processing.
		break
	}

	// Update the outer loop index to the last consumed token.
	*index = last
	return strconv.FormatInt(total, 10), true, total
}
