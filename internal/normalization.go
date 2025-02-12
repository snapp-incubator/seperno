package internal

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/snapp-incubator/seperno/pkg/options"
)

type Normalize struct {
	convertHalfSpaceToSpace bool
	urlRemover              bool
	outerSpaceRemover       bool
	spaceCombiner           bool
	normalizePunctuations   bool
	endsWithEndOfLineChar   bool
	intToWord               bool
	convertNumberLang       string
}

func NewNormalizer(conf options.NormalizerOptions) *Normalize {
	return &Normalize{
		convertHalfSpaceToSpace: conf.ConvertHalfSpaceToSpace,
		urlRemover:              conf.URLRemover,
		outerSpaceRemover:       conf.OuterSpaceRemover,
		spaceCombiner:           conf.SpaceCombiner,
		normalizePunctuations:   conf.NormalizePunctuations,
		endsWithEndOfLineChar:   conf.EndsWithEndOfLineChar,
		intToWord:               conf.IntToWord,
		convertNumberLang:       string(conf.ConvertNumberLang),
	}
}

var (
	multiSpaceRegex = regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
	urlRemovalRegex = regexp.MustCompile(`https?://[^\s]+`)
	numberRegex     = regexp.MustCompile(`\b\d+\b`)
	// Create a regex pattern to remove spaces from start and end
	outerSpaceRegex = regexp.MustCompile(`^\s+|\s+$`)
)

// FindHalfSpace replaces a specific Unicode half-space with the given string representation
func (n Normalize) FindHalfSpace(input, halfSpace string) string {
	return strings.ReplaceAll(input, "\u200c", halfSpace)
}

// SpaceNormalizer normalizes spaces in the given input string based on the provided flag
func (n Normalize) spaceNormalizer(input string) string {
	// Replace specific HTML representation
	input = strings.ReplaceAll(input, "&zwnj;", " ")

	// Convert input to a rune slice
	inputRunes := []rune(input)

	// Iterate through the runes
	for i := 0; i < len(inputRunes); i++ {
		if n.convertHalfSpaceToSpace {
			switch inputRunes[i] {
			case spaceZeroWidthNonJoiner,
				space, noBreakSpace,
				zeroWidthNoBreakSpace,
				zeroWidthSpace:
				inputRunes[i] = ' ' // Replace with a space

			case zeroWidthJoiner:
				inputRunes[i] = nullChar // Replace with nullChar
			}
		} else {
			switch inputRunes[i] {
			case space, noBreakSpace,
				zeroWidthNoBreakSpace,
				zeroWidthSpace:
				inputRunes[i] = ' ' // Replace with a space

			case zeroWidthJoiner:
				inputRunes[i] = nullChar // Replace with nullChar
			}
		}
	}

	// Create a new string, trim it, and replace nullChar
	output := strings.TrimSpace(string(inputRunes))
	output = strings.ReplaceAll(output, string(nullChar), "")
	output = strings.ToLower(output)

	return output
}

// BasicNormalizer normalizes a Persian input string.
// If input is nil, it returns nil. Applies specific transformations to the input string.
func (n Normalize) BasicNormalizer(input string) string {
	if input == "" {
		return ""
	}

	// Call SpecialYehNormalizer
	input = n.specialYehNormalizer(input)

	// Apply SpaceNormalizer
	input = n.spaceNormalizer(input)

	inputRunes := n.normalizeCharacters(input)
	// Convert runes back to a string
	// Convert the rune slice back to a string, trim spaces, replace null strings, and convert to lowercase
	s := strings.TrimSpace(string(inputRunes))
	// s = strings.ReplaceAll(s, dot, nullString)
	s = strings.ReplaceAll(s, NewLine, nullString)
	s = strings.ReplaceAll(s, nullString, emptyString)
	stringInput := strings.ToLower(s)

	if n.spaceCombiner {
		stringInput = replaceMultiSpace(stringInput)
	}
	if n.urlRemover {
		stringInput = removeURLs(stringInput)
	}
	if n.normalizePunctuations {
		stringInput = normalizePunctuations(stringInput)
	}
	if n.endsWithEndOfLineChar {
		stringInput = normalizeEndsWithEndOfLineChar(stringInput)
	}
	if n.outerSpaceRemover { // should be last normalization step
		stringInput = removeOuterSpace(stringInput)
	}
	if n.intToWord {
		stringInput = replaceNumberToWords(stringInput)
	}
	return stringInput
}

func (n Normalize) normalizeCharacters(input string) []rune {
	// Convert input to a rune slice for character-by-character processing
	inputRunes := []rune(input)

	for i := 0; i < len(inputRunes); i++ {
		switch inputRunes[i] {

		// "الف" group replacements break
		case 'أ', 'ﺍ', 'إ', 'ﺁ', 'ا', 'آ', 'ٵ', 'ﴽ', 'ﺂ', 'ﺄ', 'ﺈ', 'ﺎ', 'Ĩ', 'ٱ', 'ٲ', 'ﭐ', 'ﭑ', 'ﺇ':
			inputRunes[i] = basicCharacters[2]

		// "ب" group replacements break
		case 'ﺒ', 'ٮ', 'ݕ', 'ﺏ', 'ﺐ', 'ﺑ':
			inputRunes[i] = basicCharacters[10]

		// "ژ" group replacements break
		case 'ژ', 'ﮊ':
			inputRunes[i] = basicCharacters[8]

		// "ی" group replacements break
		case 'ى', 'ئ', 'ي', 'ﯧ', 'ﯿ', 'ﻴ', 'ۍ', 'ې', 'ۑ', 'ﯤ', 'ﯼ', 'ﯽ', 'ﯾ',
			'ﺉ', 'ﺊ', 'ﻯ', 'ﻰ', 'ﻱ', 'ﻲ', 'ﻳ', 'ے':
			inputRunes[i] = basicCharacters[1]

		// group replacements break
		case 'ك', 'ڪ', 'ﮑ', 'ﻜ', 'ػ', 'ګ', 'ڬ', 'ڭ', 'ڮ',
			'ݢ', 'ݣ', 'ݤ', 'ﻙ', 'ﻛ', 'ﮏ', 'ﮐ':
			inputRunes[i] = basicCharacters[4]

		// "ه" group replacements break
		case 'ە', 'ہ', 'ﮭ', 'ھ', 'ۿ', 'ﮪ', 'ﮫ', 'ﮬ',
			'ﻪ', 'ﻫ', 'ﻬ', 'ﻩ':
			inputRunes[i] = basicCharacters[5]

		// "م" group replacements break
		case 'ﻤ', '۾', 'ݥ', 'ﻡ', 'ﻢ', 'ﻣ':
			inputRunes[i] = basicCharacters[30]

		// "ن" group replacements break
		case 'ﻨ', 'ڹ', 'ں', 'ڻ', 'ݧ', 'טּ', 'ﮟ', 'ﻥ', 'ﻦ', 'ﻧ':
			inputRunes[i] = basicCharacters[31]

		// "و" group replacements break
		case 'ﻮ', 'ؤ', 'ٷ', 'ﯣ', 'ﺆ', 'ٶ', 'ۄ', 'ۅ', 'ۆ', 'ۇ', 'ۈ', 'ۉ', 'ۊ', 'ۋ', 'ۏ',
			1928, 'ﯗ', 'ﯙ', 'ﯚ', 'ﯛ', 'ﯝ', 'ﯡ', 'ﯢ', 'ﺅ', 'ﻭ':
			inputRunes[i] = basicCharacters[3]

		// "ی" additional group replacements break
		case 'ٸ', 'ﺌ':
			inputRunes[i] = basicCharacters[1]

		// "ة به ه" group replacements break
		case 'ة', 'ۀ', 'ﺔ', 'ۂ', 'ۃ', 'ﺓ':
			inputRunes[i] = basicCharacters[5]

		// "پ" group replacements break
		case 'ﭙ', 'ݐ', 'ݒ', 'ﭖ', 'ﭗ', 'ﭘ':
			inputRunes[i] = basicCharacters[6]

		// "چ" group replacements break
		case 'ﭽ':
			inputRunes[i] = basicCharacters[7]

		// "گ" group replacements break
		case 'ﮕ', 'ڰ', 'ڲ', 'ڳ', 'ڴ', 'ﮒ', 'ﮓ', 'ﮔ', 'ﮗ', 'ﮘ', 'ﮝ':
			inputRunes[i] = basicCharacters[9]

		// "ت" group replacements break
		case 'ٹ', 'ٺ', 'ټ', 'ٿ', 'ݓ', 'ﺘ', 'ﺕ', 'ﺖ', 'ﺗ',
			'ﭞ', 'ﭟ', 'ﭠ', 'ﭡ', 'ﭥ', 'ﭦ':
			inputRunes[i] = basicCharacters[11]

		// "ث" group replacements break
		case 'ﺜ', 'ٽ', 'ݑ', 'ﺙ', 'ﺚ':
			inputRunes[i] = basicCharacters[12]

		// "چ" group replacements break
		case 'ڇ', 'ڿ', 'ݘ', 'ﭺ', 'ﭼ':
			inputRunes[i] = basicCharacters[7]

		// "ج" group replacements break
		case 'ﺠ', 'ڃ', 'ﺝ', 'ﺞ', 'ﺟ':
			inputRunes[i] = basicCharacters[13]

		// "ح" group replacements break
		case 'ﺤ', 'ځ', 'ﺡ', 'ﺢ', 'ﺣ':
			inputRunes[i] = basicCharacters[14]

		// "خ" group replacements break
		case 'ﺨ', 'ڂ', 'ݗ', 'ﺥ', 'ﺦ', 'ﺧ':
			inputRunes[i] = basicCharacters[15]

		// "د" group replacements break
		case 'ﺪ', 'ڈ', 'ډ', 'ڊ', 'ڋ', 'ڌ', 'ڍ', 'ڐ', 'ۮ', 'ﮈ', 'ﺩ':
			inputRunes[i] = basicCharacters[17]

		// "ذ" group replacements break
		case 'ﺬ':
			inputRunes[i] = basicCharacters[16]

		// "ر" group replacements break
		case 'ﺮ', 'ڑ', 'ڒ', 'ړ', 'ڔ', 'ڕ', 'ږ', 'ۯ', 'ݛ', 'ݬ', 'ﮍ', 'ﺭ':
			inputRunes[i] = basicCharacters[18]

		// "ز" group replacements break
		case 'ﺰ', 'ڗ', 'ݫ', 'ﺯ':
			inputRunes[i] = basicCharacters[19]

		// "س" group replacements break
		case 'ﺴ', 'ښ', 'ڛ', 'ݭ', 'ﺱ', 'ﺲ', 'ﺳ':
			inputRunes[i] = basicCharacters[20]

		// "ش" group replacements break
		case 'ﺸ', 'ڜ', 'ۺ', 'ݜ', 'ﺵ', 'ﺶ', 'ﺷ':
			inputRunes[i] = basicCharacters[21]

		// "ص" group replacements break
		case 'ﺼ', 'ڝ', 'ﺹ', 'ﺺ', 'ﺻ':
			inputRunes[i] = basicCharacters[22]

		// "ض" group replacements break
		case 'ﻀ', 'ۻ', 'ﺽ', 'ﺾ', 'ﺿ':
			inputRunes[i] = basicCharacters[23]

		// "ط" group replacements break
		case 'ﻄ', 'ﻁ', 'ﻂ', 'ﻃ':
			inputRunes[i] = basicCharacters[24]

		// "ظ" group replacements break
		case 'ﻈ', 'ڟ', 'ﻅ', 'ﻆ', 'ﻇ':
			inputRunes[i] = basicCharacters[25]

		// "ع" group replacements break
		case 'ﻌ', '؏', 'ڠ', 'ﻉ', 'ﻊ', 'ﻋ':
			inputRunes[i] = basicCharacters[26]

		// "غ" group replacements break
		case 'ﻐ', 'ۼ', 'ݞ', 'ݟ', 'ﻍ', 'ﻎ', 'ﻏ':
			inputRunes[i] = basicCharacters[27]

		// "ق" group replacements break
		case 'ﻘ', 'ڦ', 'ڧ', 'ڨ', 'ﻕ', 'ﻖ', 'ﻗ':
			inputRunes[i] = basicCharacters[28]

		// "ف" group replacements break
		case '؋', 'ف', 'ڢ', 'ڣ', 'ڤ', 'ڥ', 5317, 'ﻔ', 'ﻓ', 'ﻑ', 'ﻒ':
			inputRunes[i] = basicCharacters[32]

		// "ل" group replacements break
		case 'ﻠ', 'ڵ', 'ڶ', 'ڷ', 'ڸ', 'ݪ', 'ﻝ', 'ﻞ', 'ﻟ':
			inputRunes[i] = basicCharacters[29]

		// Replace specific characters with null char break
		case sokun, 65150, 65151, fatheh, 65142,
			65143, zameh, 65144, 65145, kasreh, 65147, tashdid, 64607, 64608, 64609,
			tanvinFatheh, tanvinZameh, tanvinKasreh, alefLittle, persianHamza, 1620,
			1652, 1789, 64420, 64421, 65163, 65152, rtl, ltr, cc, arabicSubscriptAlef,
			tilde, leftHalfRingBelow, dotBelow, diaeresisBelow, ringBelow,
			ogonek, verticalLineBelow, breveBelow, invertedBreveBelow,
			longStrokeOverlay, fermata, doubleBreveBelow,
			doubleRightwardsArrowBelow, plusSignBelow, lowLine,
			DIAERESIS, seagullBelow, leftAngleAbove, acuteAccent,
			upTackBelow, candrabindu, caronBelow, snakeBelow, dotAbove,
			doubleMacronBelow, zigzagAbove, graveAccent, upwardsArrowBelow,
			tildeBelow, turnedCommaAbove, caron, overline, graphemeJoiner,
			macron, enclosingDiamond, circumflexAccentBelow, asteriskBelow,
			doubleBreve, palatalizedHookBelow, longSolidusOverlay, doubleMacron,
			graveMacron, verticalTilde, equalSignBelow, latinSmallLetterX,
			clockwiseRingOverlay, anticlockwiseRingOverlay, verticalLineAbove,
			invertedBreve, commaBelow, cedilla, invertedBridgeBelow,
			macronAcute, 64610, 65148, 65149:
			inputRunes[i] = nullChar

		// "ی" group replacements break
		case arabicYehWithThreeDotsAbove, arabicYehWithInvertedV:
			inputRunes[i] = basicCharacters[1]

		// Normalize Persian and other numeral variations to English numerals break
		// Persian 0 and Thai 0 (๐)
		case faD0, arD0, enD0, 3664:
			inputRunes[i] = convertToDestNumber(enD0, n.convertNumberLang)

		// Persian 1, Thai 1 (๑), superscript 1 (¹), subscript 1 (₁), full-width 1, and others
		case faD1, arD1, enD1, 3665, 185, 8321, 65297, 3793:
			inputRunes[i] = convertToDestNumber(enD1, n.convertNumberLang)

		// Persian 2, subscript 2 (₂), superscript 2 (²), full-width 2, and others
		case faD2, arD2, enD2, 8322, 178, 2536, 3666, 1399, 65298:
			inputRunes[i] = convertToDestNumber(enD2, n.convertNumberLang)

		// Normalize Persian and other numeral variations to English numerals break
		// Persian 3, Thai 3 (๓), and others
		case faD3, arD3, enD3, 3667, 2537, 65299:
			inputRunes[i] = convertToDestNumber(enD3, n.convertNumberLang)

		// Persian 4, Thai 4 (๔), superscript 4 (⁴), and others
		case faD4, arD4, enD4, 3668, 8308, 3178, 65300:
			inputRunes[i] = convertToDestNumber(enD4, n.convertNumberLang)

		// Persian 5, full-width 5 (５)
		case faD5, arD5, enD5, 65301:
			inputRunes[i] = convertToDestNumber(enD5, n.convertNumberLang)

		// Persian 6, Thai 6 (๖)
		case faD6, arD6, enD6, 3670:
			inputRunes[i] = convertToDestNumber(enD6, n.convertNumberLang)

		// Normalize Persian and other numeral variations to English numerals break
		// Persian 7, Devanagari 7 (७), Thai 7 (๗), and others
		case faD7, arD7, enD7, 2925, 2797, 2669, 2413, 3671:
			inputRunes[i] = convertToDestNumber(enD7, n.convertNumberLang)

		// Persian 8, subscript 8 (₈)
		case faD8, arD8, enD8, 8328:
			inputRunes[i] = convertToDestNumber(enD8, n.convertNumberLang)

		// Persian 9, and other numeral variations
		case faD9, arD9, enD9, 3881, 2543, 6121:
			inputRunes[i] = convertToDestNumber(enD9, n.convertNumberLang)

		// Replace punctuation marks with Persian equivalentsbreak
		case '?':
			inputRunes[i] = '؟'
		case '%':
			inputRunes[i] = '٪'
		case ';':
			inputRunes[i] = '؛'
		case '：':
			inputRunes[i] = ':'

		// Replace commas with Persian commasbreak
		case ',':
			inputRunes[i] = '،'
		case '٬':
			inputRunes[i] = '،'
		// Normalize various dash-like characters to underscorebreak
		case '–', '˗', '־', '­', '━', '—', '─', '_', '➖', '-', 'ـ':
			inputRunes[i] = '_'
		// Replace various dash-like characters with ellipsisbreak
		case '┅', '┄', '┈':
			inputRunes[i] = '…'
		}

	}
	return inputRunes
}

func (n Normalize) VariationSelectorsRemover(input []string) []string {
	output := make([]string, len(input))
	for i, str := range input {
		output[i] = removeVariationSelectors(str)
	}
	return output
}

func removeVariationSelectors(input string) string {
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case vs0, vs1, vs2, vs3, vs4, vs5, vs6, vs7, vs8, vs9,
			vs10, vs11, vs12, vs13, vs14, vs15:
			runes[i] = ' ' // Replace with space
		}
	}
	return string(runes)
}

func (n Normalize) specialYehNormalizer(input string) string {
	var builder strings.Builder
	for _, c := range input {
		switch c {
		case 'ے', 'ﮮ', 'ﮯ', 'ۓ', 'ﮱ': // Special "yeh" characters
			builder.WriteRune(basicCharacters[1]) // StandarD "ی"
			builder.WriteRune(basicCharacters[0]) // Spacebreak
		case 'ﻩ', 'ﮦ': // Special "heh" characters
			builder.WriteRune(basicCharacters[5]) // StandarD "ه"
			builder.WriteRune(basicCharacters[0]) // space
		default:
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

// BasicNormalizerArray Normalize each string in an array with attention to Persian language.
func (n Normalize) BasicNormalizerArray(input []string) []string {
	for i := range input {
		input[i] = n.BasicNormalizer(input[i])
	}
	return input
}

// BasicNormalizerSlice Normalize each string in a slice (ArrayList equivalent in Go) with attention to Persian language.
func (n Normalize) BasicNormalizerSlice(input []string) []string {
	result := make([]string, len(input))
	for i, str := range input {
		result[i] = n.BasicNormalizer(str)
	}
	return result
}

func replaceMultiSpace(input string) string {
	return multiSpaceRegex.ReplaceAllString(input, " ")
}

func removeOuterSpace(input string) string {
	// Remove leading and trailing spaces
	return outerSpaceRegex.ReplaceAllString(input, "")
}

func removeURLs(input string) string {
	// Replace all URLs with an empty string
	return urlRemovalRegex.ReplaceAllString(input, "")
}

func replaceNumberToWords(input string) string {
	return numberRegex.ReplaceAllStringFunc(input, func(match string) string {
		if num, err := strconv.Atoi(match); err == nil {
			return IntegerToPersian(num)
		}
		return match
	})
}

func normalizePunctuations(input string) string {
	var builder strings.Builder

	// Convert input string to runes for Unicode-safe processing
	for _, r := range input {
		if containsRune(punctuations, r) {
			// Replace punctuations with space (or remove by not appending)
			builder.WriteRune(' ') // Replace with space
		} else {
			builder.WriteRune(r) // Keep other characters
		}
	}

	return builder.String()
}

func normalizeEndsWithEndOfLineChar(input string) string {
	runes := []rune(input)
	if len(runes) == 0 {
		return input // Empty string remains unchanged
	}

	lastChar := runes[len(runes)-1]
	if containsRune(endOfLinesChar, lastChar) {
		// Remove the last character by slicing the rune slice
		return string(runes[:len(runes)-1])
	}
	return input // No end-of-line character to remove
}

func containsRune(slice []rune, r rune) bool {
	for _, s := range slice {
		if s == r {
			return true
		}
	}
	return false
}
