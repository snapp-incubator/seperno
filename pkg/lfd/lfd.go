package lfd

type DetectedNumber struct {
	Number     int64
	StartIndex int
	EndIndex   int
}

type NumberDetector interface {
	DetectNumbers(text string) []DetectedNumber
}
