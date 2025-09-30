package main

/*
#include <stdbool.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/snapp-incubator/seperno/internal"
	"github.com/snapp-incubator/seperno/pkg/lfd"
	"github.com/snapp-incubator/seperno/pkg/options"
)

// Exported normalization function using C.bool
//
//export NormalizeText
func NormalizeText(input *C.char, convertHalfSpace C.bool, combineSpace C.bool, removeOuterSpace C.bool, removeURL C.bool, normalizePunctuations C.bool, endWithEOL C.bool, intToWord C.bool, language *C.char) *C.char {
	// Directly map C.bool values to Go's bool type for clarity
	normOptions := options.NormalizerOptions{
		ConvertHalfSpaceToSpace: bool(convertHalfSpace),
		SpaceCombiner:           bool(combineSpace),
		OuterSpaceRemover:       bool(removeOuterSpace),
		URLRemover:              bool(removeURL),
		NormalizePunctuations:   bool(normalizePunctuations),
		EndsWithEndOfLineChar:   bool(endWithEOL),
		IntToWord:               bool(intToWord),
		ConvertNumberLang:       options.Language(C.GoString(language)),
	}

	// Initialize the normalizer with the options
	normalizer := internal.NewNormalizer(normOptions)

	// Call the appropriate normalization method
	result := normalizer.BasicNormalizer(C.GoString(input))
	return C.CString(result)
}

//export DetectPersianNumbers
func DetectPersianNumbers(input *C.char, outNums **C.longlong, outStarts **C.int, outEnds **C.int, outLen *C.int) {
	// Convert C string -> Go string
	goInput := C.GoString(input)

	// Run your detector
	finder := &lfd.PersianNumberDetector{}
	numbers := finder.DetectNumbers(goInput)

	n := len(numbers)
	*outLen = C.int(n)

	// Allocate memory for C arrays
	nums := (*C.longlong)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.longlong(0)))))
	starts := (*C.int)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.int(0)))))
	ends := (*C.int)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.int(0)))))

	// Convert to Go slices backed by C memory
	numsSlice := unsafe.Slice(nums, n)
	startsSlice := unsafe.Slice(starts, n)
	endsSlice := unsafe.Slice(ends, n)

	// Fill arrays
	for i, number := range numbers {
		numsSlice[i] = C.longlong(number.Number)
		startsSlice[i] = C.int(number.StartIndex)
		endsSlice[i] = C.int(number.EndIndex)
	}

	// Set return pointers
	*outNums = nums
	*outStarts = starts
	*outEnds = ends
}

func main() {}
