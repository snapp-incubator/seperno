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

// ordinalNumberMap handles irregular ordinals.
var ordinalNumberMap = map[string]string{
	"اول": "1",
	"دوم": "2",
	"سوم": "3",
}

// ordinalSuffixes maps common Persian ordinal suffixes.
var ordinalSuffixes = []string{"م", "ام", "وم", "مین"}

// Precompiled helpers
var lettersRe = regexp.MustCompile(`^\p{L}+$`)

// ConvertWordsToIntFa converts Persian number words in an input to digits
// and also returns the list of integers found.
func ConvertWordsToIntFa(input string) (string, []int64) {
	// Split the input into words, preserving delimiters (spaces, commas, etc.).
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

// splitWithDelimiters splits the input string into words, digits, punctuation, and whitespace.
func splitWithDelimiters(input string) []string {
	re := regexp.MustCompile(`([\p{L}]+|[\p{N}]+|\p{P}|\s+)`)
	return re.FindAllString(input, -1)
}

// convertNumberWord converts a single or compound number word to its digit equivalent.
func convertNumberWord(word string, words []string, index *int) (string, bool, int64) {
	trimmed := strings.TrimSpace(word)
	if trimmed == "" {
		return word, false, 0
	}

	// --- Case 0: already digits (English or Persian) ---
	if isNumeric(trimmed) {
		normalized := normalizeDigits(trimmed)
		val, _ := strconv.ParseInt(normalized, 10, 64)
		return normalized, true, val
	}

	// --- Case 1: irregular ordinals like "اول" / "دوم" / "سوم" ---
	if val, ok := ordinalNumberMap[trimmed]; ok {
		numVal, _ := strconv.ParseInt(val, 10, 64)
		return val, true, numVal
	}

	// --- Case 2: suffix-based ordinals like "چهارمین", "پنجمین" ---
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

	// --- Case 3: normal words (cardinal or compound) ---
	numStr, isNumber, numVal := parseCompoundNumber(trimmed, words, index)
	if !isNumber {
		return word, false, 0
	}
	return numStr, true, numVal
}

// Checks if a string is entirely numeric (English or Persian digits).
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

// Converts Persian digits to English.
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

// parseCompoundNumber parses single or compound number words (e.g., "سی و پنج").
func parseCompoundNumber(word string, words []string, index *int) (string, bool, int64) {
	cur := *index
	// First piece must be a known number word.
	first, ok := persianNumberMap[word]
	if !ok {
		return "", false, 0
	}
	total, _ := strconv.ParseInt(first, 10, 64)
	last := cur

	// Try to consume sequences: (spaces) "و" (spaces) <number-word>
	for {
		j := last + 1
		// skip whitespace
		for j < len(words) && strings.TrimSpace(words[j]) == "" {
			j++
		}
		// need "و"
		if j >= len(words) || words[j] != "و" {
			break
		}
		j++ // past "و"
		// skip whitespace after "و"
		for j < len(words) && strings.TrimSpace(words[j]) == "" {
			j++
		}
		if j >= len(words) || !lettersRe.MatchString(words[j]) {
			// not a valid next word -> stop (don't consume "و")
			break
		}
		// next must be a number word
		if nextStr, ok := persianNumberMap[words[j]]; ok {
			v, _ := strconv.ParseInt(nextStr, 10, 64)
			total += v
			last = j // consume up to this token
			continue
		}
		// next word isn't a number -> stop (don't consume)
		break
	}

	// Move outer loop index to the last consumed number token.
	*index = last
	return strconv.FormatInt(total, 10), true, total
}
