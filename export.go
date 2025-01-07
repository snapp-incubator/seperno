package main

/*
#include <stdbool.h>
*/
import "C"
import (
	"github.com/snapp-incubator/seperno/internal"
)

// Exported normalization function using C.bool
//
//export NormalizeText
func NormalizeText(input *C.char, convertHalfSpace C.bool, combineSpace C.bool, removeOuterSpace C.bool, removeURL C.bool, normalizePunctuations C.bool, endWithEOL C.bool) *C.char {
	// Directly map C.bool values to Go's bool type for clarity
	options := internal.NormalizerOptions{
		ConvertHalfSpaceToSpace: bool(convertHalfSpace),
		SpaceCombiner:           bool(combineSpace),
		OuterSpaceRemover:       bool(removeOuterSpace),
		URLRemover:              bool(removeURL),
		NormalizePunctuations:   bool(normalizePunctuations),
		EndsWithEndOfLineChar:   bool(endWithEOL),
	}

	// Initialize the normalizer with the options
	normalizer := internal.NewNormalizer(options)

	// Call the appropriate normalization method
	result := normalizer.BasicNormalizer(C.GoString(input))
	return C.CString(result)
}

func main() {}
