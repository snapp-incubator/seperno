package internal

import (
	"fmt"
	"strings"
)

var (
	iranianMegas    = []string{"", "هزار", "میلیون", "میلیارد", "بیلیون", "بیلیارد", "تریلیون", "تریلیارد"}
	iranianUnits    = []string{"", "یک", "دو", "سه", "چهار", "پنج", "شش", "هفت", "هشت", "نه"}
	iranianTens     = []string{"", "ده", "بیست", "سی", "چهل", "پنجاه", "شصت", "هفتاد", "هشتاد", "نود"}
	iranianTeens    = []string{"ده", "یازده", "دوازده", "سیزده", "چهارده", "پانزده", "شانزده", "هفده", "هجده", "نوزده"}
	iranianHundreds = []string{"", "صد", "دویست", "سیصد", "چهارصد", "پانصد", "ششصد", "هفتصد", "هشتصد", "نهصد"}
)

func IntegerToPersian(input int) string {
	if input == 0 {
		return "صفر"
	}

	words := make([]string, 0, 10) // Pre-allocating slice for performance
	if input < 0 {
		words = append(words, "منفی")
		input = -input
	}

	triplets := integerToTriplets(input)
	for idx := len(triplets) - 1; idx >= 0; idx-- {
		if triplet := triplets[idx]; triplet > 0 {
			words = append(words, convertTriplet(triplet))
			if mega := iranianMegas[idx]; mega != "" {
				words = append(words, mega)
			}
		}
	}

	return strings.Join(words, " ")
}

func convertTriplet(triplet int) string {
	if triplet == 0 {
		return ""
	}

	hundreds := triplet / 100
	tens := (triplet / 10) % 10
	units := triplet % 10

	parts := make([]string, 0, 3)
	if hundreds > 0 {
		parts = append(parts, iranianHundreds[hundreds])
	}

	switch {
	case tens == 1:
		parts = append(parts, iranianTeens[units])
	case tens > 1:
		if units > 0 {
			parts = append(parts, fmt.Sprintf("%s و %s", iranianTens[tens], iranianUnits[units]))
		} else {
			parts = append(parts, iranianTens[tens])
		}
	case units > 0:
		parts = append(parts, iranianUnits[units])
	}

	return strings.Join(parts, " و ")
}

func integerToTriplets(number int) []int {
	var triplets []int

	for number > 0 {
		triplets = append(triplets, number%1_000)
		number = number / 1_000
	}

	return triplets
}
