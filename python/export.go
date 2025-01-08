package main

/*
#include <stdbool.h>
*/
import "C"
import (
	"github.com/snapp-incubator/seperno/internal"
	"unsafe"
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

// Go function to calculate cosine similarity using the internal function
func calculateCosineSimilarity(text []string, list_text [][]string) []float64 {
	return internal.CalculateCosineSimilarity(text, list_text)
}

// Helper function to convert **C.char to []string in Go
func cStringArrayToSlice(cArray **C.char, length C.int) []string {
	goSlice := make([]string, length)
	for i := 0; i < int(length); i++ {
		// Corrected: Casting correctly for double pointer dereferencing
		cString := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(cArray)))
		goSlice[i] = C.GoString(cString)
	}
	return goSlice
}

// Helper function to convert ***C.char to [][]string in Go
func cStringMatrixToSlice(cMatrix ***C.char, outerLen, innerLen C.int) [][]string {
	goMatrix := make([][]string, outerLen)
	for i := 0; i < int(outerLen); i++ {
		// Corrected: Properly dereferencing ***C.char
		rowPointer := *(***C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cMatrix)) + uintptr(i)*unsafe.Sizeof(cMatrix)))
		goMatrix[i] = cStringArrayToSlice(rowPointer, innerLen)
	}
	return goMatrix
}

//export CosineSimilarity
func CosineSimilarity(input **C.char, list_text ***C.char, inputLen, listLen, listInnerLen C.int) *C.double {
	// Convert the input C arrays to Go slices
	goInput := cStringArrayToSlice(input, inputLen)
	goListText := cStringMatrixToSlice(list_text, listLen, listInnerLen)

	// Perform the cosine similarity calculation
	similarities := calculateCosineSimilarity(goInput, goListText)

	// Convert Go slice of float64 back to *C.double array
	cSimilarities := (*C.double)(C.malloc(C.size_t(len(similarities)) * C.size_t(unsafe.Sizeof(C.double(0)))))
	for i, value := range similarities {
		*((*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(cSimilarities)) + uintptr(i)*unsafe.Sizeof(C.double(0))))) = C.double(value)
	}

	return cSimilarities
}

func main() {}
