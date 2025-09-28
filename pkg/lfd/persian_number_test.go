package lfd

import (
	"reflect"
	"testing"
)

func TestConvertWordsToIntFa(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		numbers  []DetectedNumber
	}{
		// Simple single digits
		{
			name:     "ordinal_third",
			input:    "سوم",
			expected: "3",
			numbers:  []DetectedNumber{{Number: 3, StartIndex: 0, EndIndex: 2}},
		},
		{
			name:     "ordinal_first",
			input:    "اول",
			expected: "1",
			numbers:  []DetectedNumber{{Number: 1, StartIndex: 0, EndIndex: 2}},
		},
		{
			name:     "ordinal_second_with_suffix",
			input:    "دومین",
			expected: "2",
			numbers:  []DetectedNumber{{Number: 2, StartIndex: 0, EndIndex: 4}},
		},
		{
			name:     "ordinal_fourth_with_suffix",
			input:    "چهارمین",
			expected: "4",
			numbers:  []DetectedNumber{{Number: 4, StartIndex: 0, EndIndex: 6}},
		},
		{
			name:     "street_one",
			input:    "خیابان یک",
			expected: "خیابان 1",
			numbers:  []DetectedNumber{{Number: 1, StartIndex: 7, EndIndex: 8}},
		},
		{
			name:     "alley_two",
			input:    "کوچه دو",
			expected: "کوچه 2",
			numbers:  []DetectedNumber{{Number: 2, StartIndex: 5, EndIndex: 6}},
		},
		{
			name:     "street_number_nine",
			input:    "خیابان شماره نه",
			expected: "خیابان شماره 9",
			numbers:  []DetectedNumber{{Number: 9, StartIndex: 13, EndIndex: 14}},
		},

		// Teens
		{
			name:     "plate_eleven",
			input:    "پلاک یازده",
			expected: "پلاک 11",
			numbers:  []DetectedNumber{{Number: 11, StartIndex: 5, EndIndex: 9}},
		},
		{
			name:     "plate_thirteen",
			input:    "پلاک سیزده",
			expected: "پلاک 13",
			numbers:  []DetectedNumber{{Number: 13, StartIndex: 5, EndIndex: 9}},
		},

		// Tens
		{
			name:     "street_twenty",
			input:    "خیابان بیست",
			expected: "خیابان 20",
			numbers:  []DetectedNumber{{Number: 20, StartIndex: 7, EndIndex: 10}},
		},
		{
			name:     "street_thirty",
			input:    "خیابان سی",
			expected: "خیابان 30",
			numbers:  []DetectedNumber{{Number: 30, StartIndex: 7, EndIndex: 8}},
		},
		{
			name:     "street_forty",
			input:    "خیابان چهل",
			expected: "خیابان 40",
			numbers:  []DetectedNumber{{Number: 40, StartIndex: 7, EndIndex: 9}},
		},

		// Compound numbers
		{
			name:     "street_twenty_three",
			input:    "خیابان بیست و سه",
			expected: "خیابان 23",
			numbers:  []DetectedNumber{{Number: 23, StartIndex: 7, EndIndex: 15}},
		},
		{
			name:     "alley_thirty_five",
			input:    "کوچه سی و پنج",
			expected: "کوچه 35",
			numbers:  []DetectedNumber{{Number: 35, StartIndex: 5, EndIndex: 12}},
		},
		{
			name:     "plate_forty_seven",
			input:    "پلاک چهل و هفت",
			expected: "پلاک 47",
			numbers:  []DetectedNumber{{Number: 47, StartIndex: 5, EndIndex: 13}},
		},

		// Addresses with multiple numbers
		{
			name:     "address_multiple_numbers",
			input:    "خیابان بیست و چهار پلاک ده طبقه سوم",
			expected: "خیابان 24 پلاک 10 طبقه 3",
			numbers:  []DetectedNumber{{Number: 24, StartIndex: 7, EndIndex: 17}, {Number: 10, StartIndex: 24, EndIndex: 25}, {Number: 3, StartIndex: 32, EndIndex: 34}},
		},
		{
			name:     "address_ordinal_street",
			input:    "خیابان پنجم پلاک دوازده واحد هفت",
			expected: "خیابان 5 پلاک 12 واحد 7",
			numbers:  []DetectedNumber{{Number: 5, StartIndex: 7, EndIndex: 10}, {Number: 12, StartIndex: 17, EndIndex: 22}, {Number: 7, StartIndex: 29, EndIndex: 31}},
		},

		// Mixed digits and words
		{
			name:     "mixed_word_and_digit",
			input:    "پلاک بیست 22",
			expected: "پلاک 20 22",
			numbers:  []DetectedNumber{{Number: 20, StartIndex: 5, EndIndex: 8}, {Number: 22, StartIndex: 10, EndIndex: 11}},
		},
		{
			name:     "duplicate_numbers",
			input:    "واحد سه 3",
			expected: "واحد 3 3",
			numbers:  []DetectedNumber{{Number: 3, StartIndex: 5, EndIndex: 6}, {Number: 3, StartIndex: 8, EndIndex: 8}},
		},

		// Phone numbers
		{
			name:     "phone_number",
			input:    "تلفن صفر نهصد و بیست",
			expected: "تلفن 0 920",
			numbers:  []DetectedNumber{{Number: 0, StartIndex: 5, EndIndex: 7}, {Number: 920, StartIndex: 9, EndIndex: 19}},
		},

		// Edge cases
		{
			name:     "no_numbers",
			input:    "خیابان بدون شماره",
			expected: "خیابان بدون شماره",
			numbers:  []DetectedNumber{},
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
			numbers:  []DetectedNumber{},
		},
		// Larger numbers
		{
			name:     "hundred_twenty_three",
			input:    "صد و بیست و سه",
			expected: "123",
			numbers:  []DetectedNumber{{Number: 123, StartIndex: 0, EndIndex: 13}},
		},
		{
			name:     "two_hundred_fifty_six",
			input:    "دویست و پنجاه و شش",
			expected: "256",
			numbers:  []DetectedNumber{{Number: 256, StartIndex: 0, EndIndex: 17}},
		},
		{
			name:     "thousand_two_hundred",
			input:    "هزار و دویست",
			expected: "1200",
			numbers:  []DetectedNumber{{Number: 1200, StartIndex: 0, EndIndex: 11}},
		},
		{
			name:     "nineteen_ninety_nine",
			input:    "هزار و نهصد و نود و نه",
			expected: "1999",
			numbers:  []DetectedNumber{{Number: 1999, StartIndex: 0, EndIndex: 21}},
		},

		// Ordinals beyond 10
		{
			name:     "ordinal_twelfth",
			input:    "دوازدهم",
			expected: "12",
			numbers:  []DetectedNumber{{Number: 12, StartIndex: 0, EndIndex: 6}},
		},
		{
			name:     "ordinal_thirteenth",
			input:    "سیزدهم",
			expected: "13",
			numbers:  []DetectedNumber{{Number: 13, StartIndex: 0, EndIndex: 5}},
		},

		// Mix of thousands and smaller units
		{
			name:     "complex_address",
			input:    "خیابان هزار و بیست پلاک دویست و سه",
			expected: "خیابان 1020 پلاک 203",
			numbers:  []DetectedNumber{{Number: 1020, StartIndex: 7, EndIndex: 17}, {Number: 203, StartIndex: 24, EndIndex: 33}},
		},

		// Persian digits in bigger numbers
		{
			name:     "persian_digits",
			input:    "پلاک ۱۲۳۴",
			expected: "پلاک 1234",
			numbers:  []DetectedNumber{{Number: 1234, StartIndex: 5, EndIndex: 8}},
		},
		{
			name:     "compound_thousand_number",
			input:    "نمایندگی هزارو هفتصدوهفده کرمان موتور",
			expected: "نمایندگی 1717 کرمان موتور",
			numbers:  []DetectedNumber{{Number: 1717, StartIndex: 9, EndIndex: 24}},
		},
		{
			name:     "compound_thousand_number_2",
			input:    "پلاک صد و هفتادو چهار زنگ بیستوسه طبقه چهار",
			expected: "پلاک 174 زنگ 23 طبقه 4",
			numbers:  []DetectedNumber{{Number: 174, StartIndex: 5, EndIndex: 20}, {Number: 23, StartIndex: 26, EndIndex: 32}, {Number: 4, StartIndex: 39, EndIndex: 42}},
		},
		{
			name:     "thousand_twenty_combination",
			input:    "تعداد هزار و بیست",
			expected: "تعداد 1020",
			numbers:  []DetectedNumber{{Number: 1020, StartIndex: 6, EndIndex: 16}},
		},
		{
			name:     "separated_hundred_ten",
			input:    "یک صد و ده تا تکه",
			expected: "110 تا تکه",
			numbers:  []DetectedNumber{{Number: 110, StartIndex: 0, EndIndex: 9}},
		},
		{
			name:     "joined_hundred_fifty",
			input:    "یکصد و پنجاه",
			expected: "150",
			numbers:  []DetectedNumber{{Number: 150, StartIndex: 0, EndIndex: 11}},
		},
		{
			name:     "twenty_thousand_thirty_two",
			input:    "بیست هزار و سی و دو تا تسک دارم",
			expected: "20032 تا تسک دارم",
			numbers:  []DetectedNumber{{Number: 20032, StartIndex: 0, EndIndex: 18}},
		},
		{
			name:     "two_thousand_thirty_two",
			input:    "دو هزار و سی و دو تا تسک دارم",
			expected: "2032 تا تسک دارم",
			numbers:  []DetectedNumber{{Number: 2032, StartIndex: 0, EndIndex: 16}},
		},
		{
			name:     "three_hundred_thousand",
			input:    "سیصد هزار و سی و دو تا تسک دارم",
			expected: "300032 تا تسک دارم",
			numbers:  []DetectedNumber{{Number: 300032, StartIndex: 0, EndIndex: 18}},
		},
		{
			name:     "separated_four_hundred_ten",
			input:    "چهار صد و ده",
			expected: "410",
			numbers:  []DetectedNumber{{Number: 410, StartIndex: 0, EndIndex: 11}},
		},
		{
			name:     "compound_ordinal_twenty_fifth",
			input:    "بیست و پنجم نمایشگاه آفرود",
			expected: "25 نمایشگاه آفرود",
			numbers:  []DetectedNumber{{Number: 25, StartIndex: 0, EndIndex: 10}},
		},
		{
			name:     "compound_ordinal_with_suffix",
			input:    "بیست و پنجمین نمایشگاه آفرود",
			expected: "25 نمایشگاه آفرود",
			numbers:  []DetectedNumber{{Number: 25, StartIndex: 0, EndIndex: 12}},
		},
		{
			name:     "multiple_ordinals",
			input:    "این اولین و دومین تست است",
			expected: "این 1 و 2 تست است",
			numbers:  []DetectedNumber{{Number: 1, StartIndex: 4, EndIndex: 8}, {Number: 2, StartIndex: 12, EndIndex: 16}},
		},
	}

	detector := &PersianNumberDetector{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := detector.DetectNumbers(tt.input)
			if !reflect.DeepEqual(results, tt.numbers) {
				t.Errorf("input: %v, detected numbers: %v, want: %v", tt.input, results, tt.numbers)
			}
		})
	}
}
