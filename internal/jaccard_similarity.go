package internal

// GenerateBigrams generates a set of bigrams from a given string
func GenerateBigrams(text string) map[string]struct{} {
	bigrams := make(map[string]struct{})
	textLen := len(text)

	if textLen < 2 {
		return bigrams
	}

	for i := 0; i < textLen-1; i++ {
		bigram := text[i : i+2]
		bigrams[bigram] = struct{}{}
	}

	return bigrams
}

// SynonymJaccardSimilarity calculates the maximum Jaccard similarity between a query and a set of synonyms
func SynonymJaccardSimilarity(query string, synonymSet []string) float64 {
	bigramsQuery := GenerateBigrams(query)

	// If the query has no bigrams, return 0
	if len(bigramsQuery) == 0 {
		return 0.0
	}

	// Calculate Jaccard similarity for each synonym and find the maximum
	maxSimilarity := 0.0
	for _, synonym := range synonymSet {
		bigramsSynonym := GenerateBigrams(synonym)
		intersection := 0
		for bigram := range bigramsQuery {
			if _, exists := bigramsSynonym[bigram]; exists {
				intersection++
			}
		}
		union := len(bigramsQuery) + len(bigramsSynonym) - intersection
		if union > 0 {
			similarity := float64(intersection) / float64(union)
			if similarity > maxSimilarity {
				maxSimilarity = similarity
			}
		}
	}

	return maxSimilarity
}

func CalculateJaccardSimilarity(text []string, list_text [][]string) []float64 {
	var simialrities []float64
	for i, query := range text {
		simialrities = append(simialrities, SynonymJaccardSimilarity(query, list_text[i]))
	}
	return simialrities
}
