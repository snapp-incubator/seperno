package internal

import (
	"reflect"
	"testing"
)

func TestConvertWordsToIntFa(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		numbers  []int64
	}{
		// Simple single digits
		{
			name:     "ordinal_third",
			input:    "سوم",
			expected: "3",
			numbers:  []int64{3},
		},
		{
			name:     "ordinal_first",
			input:    "اول",
			expected: "1",
			numbers:  []int64{1},
		},
		{
			name:     "ordinal_second_with_suffix",
			input:    "دومین",
			expected: "2",
			numbers:  []int64{2},
		},
		{
			name:     "ordinal_fourth_with_suffix",
			input:    "چهارمین",
			expected: "4",
			numbers:  []int64{4},
		},
		{
			name:     "street_one",
			input:    "خیابان یک",
			expected: "خیابان 1",
			numbers:  []int64{1},
		},
		{
			name:     "alley_two",
			input:    "کوچه دو",
			expected: "کوچه 2",
			numbers:  []int64{2},
		},
		{
			name:     "street_number_nine",
			input:    "خیابان شماره نه",
			expected: "خیابان شماره 9",
			numbers:  []int64{9},
		},

		// Teens
		{
			name:     "plate_eleven",
			input:    "پلاک یازده",
			expected: "پلاک 11",
			numbers:  []int64{11},
		},
		{
			name:     "plate_thirteen",
			input:    "پلاک سیزده",
			expected: "پلاک 13",
			numbers:  []int64{13},
		},

		// Tens
		{
			name:     "street_twenty",
			input:    "خیابان بیست",
			expected: "خیابان 20",
			numbers:  []int64{20},
		},
		{
			name:     "street_thirty",
			input:    "خیابان سی",
			expected: "خیابان 30",
			numbers:  []int64{30},
		},
		{
			name:     "street_forty",
			input:    "خیابان چهل",
			expected: "خیابان 40",
			numbers:  []int64{40},
		},

		// Compound numbers
		{
			name:     "street_twenty_three",
			input:    "خیابان بیست و سه",
			expected: "خیابان 23",
			numbers:  []int64{23},
		},
		{
			name:     "alley_thirty_five",
			input:    "کوچه سی و پنج",
			expected: "کوچه 35",
			numbers:  []int64{35},
		},
		{
			name:     "plate_forty_seven",
			input:    "پلاک چهل و هفت",
			expected: "پلاک 47",
			numbers:  []int64{47},
		},

		// Addresses with multiple numbers
		{
			name:     "address_multiple_numbers",
			input:    "خیابان بیست و چهار پلاک ده طبقه سوم",
			expected: "خیابان 24 پلاک 10 طبقه 3",
			numbers:  []int64{24, 10, 3},
		},
		{
			name:     "address_ordinal_street",
			input:    "خیابان پنجم پلاک دوازده واحد هفت",
			expected: "خیابان 5 پلاک 12 واحد 7",
			numbers:  []int64{5, 12, 7},
		},

		// Mixed digits and words
		{
			name:     "mixed_word_and_digit",
			input:    "پلاک بیست 22",
			expected: "پلاک 20 22",
			numbers:  []int64{20, 22},
		},
		{
			name:     "duplicate_numbers",
			input:    "واحد سه 3",
			expected: "واحد 3 3",
			numbers:  []int64{3, 3},
		},

		// Phone numbers
		{
			name:     "phone_number",
			input:    "تلفن صفر نهصد و بیست",
			expected: "تلفن 0 920",
			numbers:  []int64{0, 920},
		},

		// Edge cases
		{
			name:     "no_numbers",
			input:    "خیابان بدون شماره",
			expected: "خیابان بدون شماره",
			numbers:  nil,
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
			numbers:  nil,
		},
		// Larger numbers
		{
			name:     "hundred_twenty_three",
			input:    "صد و بیست و سه",
			expected: "123",
			numbers:  []int64{123},
		},
		{
			name:     "two_hundred_fifty_six",
			input:    "دویست و پنجاه و شش",
			expected: "256",
			numbers:  []int64{256},
		},
		{
			name:     "thousand_two_hundred",
			input:    "هزار و دویست",
			expected: "1200",
			numbers:  []int64{1200},
		},
		{
			name:     "nineteen_ninety_nine",
			input:    "هزار و نهصد و نود و نه",
			expected: "1999",
			numbers:  []int64{1999},
		},

		// Ordinals beyond 10
		{
			name:     "ordinal_twelfth",
			input:    "دوازدهم",
			expected: "12",
			numbers:  []int64{12},
		},
		{
			name:     "ordinal_thirteenth",
			input:    "سیزدهم",
			expected: "13",
			numbers:  []int64{13},
		},

		// Mix of thousands and smaller units
		{
			name:     "complex_address",
			input:    "خیابان هزار و بیست پلاک دویست و سه",
			expected: "خیابان 1020 پلاک 203",
			numbers:  []int64{1020, 203},
		},

		// Persian digits in bigger numbers
		{
			name:     "persian_digits",
			input:    "پلاک ۱۲۳۴",
			expected: "پلاک 1234",
			numbers:  []int64{1234},
		},
		{
			name:     "compound_thousand_number",
			input:    "نمایندگی هزارو هفتصدوهفده کرمان موتور",
			expected: "نمایندگی 1717 کرمان موتور",
			numbers:  []int64{1717},
		},
		{
			name:     "thousand_twenty_combination",
			input:    "تعداد هزار و بیست",
			expected: "تعداد 1020",
			numbers:  []int64{1020},
		},
		{
			name:     "separated_hundred_ten",
			input:    "یک صد و ده تا تکه",
			expected: "110 تا تکه",
			numbers:  []int64{110},
		},
		{
			name:     "joined_hundred_fifty",
			input:    "یکصد و پنجاه",
			expected: "150",
			numbers:  []int64{150},
		},
		{
			name:     "twenty_thousand_thirty_two",
			input:    "بیست هزار و سی و دو تا تسک دارم",
			expected: "20032 تا تسک دارم",
			numbers:  []int64{20032},
		},
		{
			name:     "two_thousand_thirty_two",
			input:    "دو هزار و سی و دو تا تسک دارم",
			expected: "2032 تا تسک دارم",
			numbers:  []int64{2032},
		},
		{
			name:     "three_hundred_thousand",
			input:    "سیصد هزار و سی و دو تا تسک دارم",
			expected: "300032 تا تسک دارم",
			numbers:  []int64{300032},
		},
		{
			name:     "separated_four_hundred_ten",
			input:    "چهار صد و ده",
			expected: "410",
			numbers:  []int64{410},
		},
		{
			name:     "compound_ordinal_twenty_fifth",
			input:    "بیست و پنجم نمایشگاه آفرود",
			expected: "25 نمایشگاه آفرود",
			numbers:  []int64{25},
		},
		{
			name:     "compound_ordinal_with_suffix",
			input:    "بیست و پنجمین نمایشگاه آفرود",
			expected: "25 نمایشگاه آفرود",
			numbers:  []int64{25},
		},
		{
			name:     "multiple_ordinals",
			input:    "این اولین و دومین تست است",
			expected: "این 1 و 2 تست است",
			numbers:  []int64{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, nums := ConvertWordsToIntFa(tt.input)
			if got != tt.expected || !reflect.DeepEqual(nums, tt.numbers) {
				t.Errorf("error for input: %s", tt.input)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
			if !reflect.DeepEqual(nums, tt.numbers) {
				t.Errorf("numbers %v, want %v", nums, tt.numbers)
			}
		})
	}
}
