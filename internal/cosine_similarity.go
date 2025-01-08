package internal

import "C"
import (
	"math"
	"strings"
)

// Function to tokenize and build a consistent vocabulary
func buildVocabulary(query string, synonyms []string) map[string]int {
	vocab := make(map[string]int)
	index := 0

	// Add query words to the vocabulary
	words := strings.Fields(query)
	for _, word := range words {
		if _, exists := vocab[word]; !exists {
			vocab[word] = index
			index++
		}
	}

	// Add synonym words to the vocabulary
	for _, synonym := range synonyms {
		synWords := strings.Fields(synonym)
		for _, word := range synWords {
			if _, exists := vocab[word]; !exists {
				vocab[word] = index
				index++
			}
		}
	}

	return vocab
}

// Function to create a count vector based on a vocabulary
func createVector(text string, vocab map[string]int) []float64 {
	vector := make([]float64, len(vocab))

	// Tokenize the text and populate the vector
	words := strings.Fields(text)
	for _, word := range words {
		if index, exists := vocab[word]; exists {
			vector[index]++
		}
	}

	return vector
}

// Function to calculate cosine similarity
func cosineSimilarity(vecA, vecB []float64) float64 {
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0

	for i := range vecA {
		dotProduct += vecA[i] * vecB[i]
		normA += vecA[i] * vecA[i]
		normB += vecB[i] * vecB[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Function to compute cosine similarity for a query and a set of synonyms
func synonymCosineSimilarity(query string, synonymSet []string) float64 {
	// Build a consistent vocabulary
	vocab := buildVocabulary(query, synonymSet)

	// Create the query vector
	queryVector := createVector(query, vocab)

	// Compute cosine similarity for each synonym
	maxSim := 0.0
	for _, synonym := range synonymSet {
		synonymVector := createVector(synonym, vocab)
		sim := cosineSimilarity(queryVector, synonymVector)
		if sim > maxSim {
			maxSim = sim
		}
	}

	return maxSim
}

func CalculateCosineSimilarity(text []string, list_text [][]string) []float64 {
	var simialrities []float64
	for i, query := range text {
		simialrities = append(simialrities, synonymCosineSimilarity(query, list_text[i]))
	}
	return simialrities
}
