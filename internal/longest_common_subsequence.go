package internal

// lcsSimilarity calculates the length of the longest common subsequence between two strings.
func lcsSimilarity(s1, s2 string) float64 {
	lenS1 := len(s1)
	lenS2 := len(s2)

	// If one string is empty, the LCS length is 0
	if lenS1 == 0 || lenS2 == 0 {
		return 0
	}

	// Initialize two slices for dynamic programming to optimize space complexity
	prev := make([]float64, lenS2+1)
	curr := make([]float64, lenS2+1)

	// Fill the DP table
	for i := 1; i <= lenS1; i++ {
		for j := 1; j <= lenS2; j++ {
			if s1[i-1] == s2[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}
		// Swap current and previous rows
		prev, curr = curr, prev
	}

	// LCS length is in the last element of the previous row
	return prev[lenS2]
}

// longestCommonSubsequence returns the length of the longest common subsequence.
func longestCommonSubsequence(s1, s2 string) float64 {
	return lcsSimilarity(s1, s2)
}

// synonymLongestCommonSubsequence calculates the maximum LCS length between a query and a set of synonyms.
func synonymLongestCommonSubsequence(query string, synonymSet []string) float64 {
	lcsSeqArray := make([]float64, len(synonymSet))
	for i, synonym := range synonymSet {
		lcsSeqArray[i] = longestCommonSubsequence(query, synonym)
	}
	return maxSlice(lcsSeqArray)
}

// CalculateLongestCommonSubsequence calculates LCS similarities for a list of queries against a list of synonym sets.
func CalculateLongestCommonSubsequence(queries []string, synonymSets [][]string) []float64 {
	similarities := make([]float64, len(queries))
	for i := 0; i < len(queries); i++ {
		similarities[i] = synonymLongestCommonSubsequence(queries[i], synonymSets[i])
	}
	return similarities
}

// maxSlice returns the maximum value in a slice of integers.
func maxSlice(slice []float64) float64 {
	maxVal := slice[0]
	for _, v := range slice {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}
