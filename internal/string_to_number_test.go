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
			name:     "ordinal سوم",
			input:    "سوم",
			expected: "3",
			numbers:  []int64{3},
		},
		{
			name:     "ordinal اول",
			input:    "اول",
			expected: "1",
			numbers:  []int64{1},
		},
		{
			name:     "ordinal دومین",
			input:    "دومین",
			expected: "2",
			numbers:  []int64{2},
		},
		{
			name:     "ordinal چهارمین",
			input:    "چهارمین",
			expected: "4",
			numbers:  []int64{4},
		},
		{
			name:     "خیابان یک",
			input:    "خیابان یک",
			expected: "خیابان 1",
			numbers:  []int64{1},
		},
		{
			name:     "کوچه دو",
			input:    "کوچه دو",
			expected: "کوچه 2",
			numbers:  []int64{2},
		},
		{
			name:     "خیابان شماره نه",
			input:    "خیابان شماره نه",
			expected: "خیابان شماره 9",
			numbers:  []int64{9},
		},

		// Teens
		{
			name:     "پلاک یازده",
			input:    "پلاک یازده",
			expected: "پلاک 11",
			numbers:  []int64{11},
		},
		{
			name:     "پلاک سیزده",
			input:    "پلاک سیزده",
			expected: "پلاک 13",
			numbers:  []int64{13},
		},

		// Tens
		{
			name:     "خیابان بیست",
			input:    "خیابان بیست",
			expected: "خیابان 20",
			numbers:  []int64{20},
		},
		{
			name:     "خیابان سی",
			input:    "خیابان سی",
			expected: "خیابان 30",
			numbers:  []int64{30},
		},
		{
			name:     "خیابان چهل",
			input:    "خیابان چهل",
			expected: "خیابان 40",
			numbers:  []int64{40},
		},

		// Compound numbers
		{
			name:     "خیابان بیست و سه",
			input:    "خیابان بیست و سه",
			expected: "خیابان 23",
			numbers:  []int64{23},
		},
		{
			name:     "کوچه سی و پنج",
			input:    "کوچه سی و پنج",
			expected: "کوچه 35",
			numbers:  []int64{35},
		},
		{
			name:     "پلاک چهل و هفت",
			input:    "پلاک چهل و هفت",
			expected: "پلاک 47",
			numbers:  []int64{47},
		},

		// Addresses with multiple numbers
		{
			name:     "خیابان بیست و چهار پلاک ده طبقه سوم",
			input:    "خیابان بیست و چهار پلاک ده طبقه سوم",
			expected: "خیابان 24 پلاک 10 طبقه 3",
			numbers:  []int64{24, 10, 3},
		},
		{
			name:     "خیابان پنجم پلاک دوازده واحد هفت",
			input:    "خیابان پنجم پلاک دوازده واحد هفت",
			expected: "خیابان 5 پلاک 12 واحد 7",
			numbers:  []int64{5, 12, 7},
		},

		// Mixed digits and words
		{
			name:     "پلاک بیست 22",
			input:    "پلاک بیست 22",
			expected: "پلاک 20 22",
			numbers:  []int64{20, 22},
		},
		{
			name:     "واحد سه 3",
			input:    "واحد سه 3",
			expected: "واحد 3 3",
			numbers:  []int64{3, 3},
		},

		// Phone numbers
		{
			name:     "تلفن صفر نهصد و بیست",
			input:    "تلفن صفر نهصد و بیست",
			expected: "تلفن 0 920",
			numbers:  []int64{0, 920},
		},

		// Edge cases
		{
			name:     "خیابان بدون شماره",
			input:    "خیابان بدون شماره",
			expected: "خیابان بدون شماره",
			numbers:  nil,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			numbers:  nil,
		},
		// Larger numbers
		{
			name:     "صد و بیست و سه",
			input:    "صد و بیست و سه",
			expected: "123",
			numbers:  []int64{123},
		},
		{
			name:     "دویست و پنجاه و شش",
			input:    "دویست و پنجاه و شش",
			expected: "256",
			numbers:  []int64{256},
		},
		{
			name:     "هزار و دویست",
			input:    "هزار و دویست",
			expected: "1200",
			numbers:  []int64{1200},
		},
		{
			name:     "هزار و نهصد و نود و نه",
			input:    "هزار و نهصد و نود و نه",
			expected: "1999",
			numbers:  []int64{1999},
		},

		// Ordinals beyond 10
		{
			name:     "دوازدهم",
			input:    "دوازدهم",
			expected: "12",
			numbers:  []int64{12},
		},
		{
			name:     "سیزدهم",
			input:    "سیزدهم",
			expected: "13",
			numbers:  []int64{13},
		},

		// Mix of thousands and smaller units
		{
			name:     "خیابان هزار و بیست پلاک دویست و سه",
			input:    "خیابان هزار و بیست پلاک دویست و سه",
			expected: "خیابان 1020 پلاک 203",
			numbers:  []int64{1020, 203},
		},

		// Persian digits in bigger numbers
		{
			name:     "پلاک ۱۲۳۴",
			input:    "پلاک ۱۲۳۴",
			expected: "پلاک 1234",
			numbers:  []int64{1234},
		},
		{
			name:     "عدد هزاری",
			input:    "نمایندگی هزارو هفتصدوهفده کرمان موتور",
			expected: "نمایندگی 1717 کرمان موتور",
			numbers:  []int64{1717},
		},
		{
			name:     "ترکیب هزار",
			input:    "تعداد هزار و بیست",
			expected: "تعداد 1020",
			numbers:  []int64{1020},
		},
		{
			name:     "یک قبل از صد بیاد",
			input:    "یک صد و ده تا تکه",
			expected: "110 تا تکه",
			numbers:  []int64{110},
		},
		{
			name:     "یک بچسبه به صد",
			input:    "یکصد و پنجاه",
			expected: "150",
			numbers:  []int64{150},
		},
		{
			name:     "بیست بچسبه به هزار",
			input:    "بیست هزار و سی و دو تا تسک دارم",
			expected: "20032 تا تسک دارم",
			numbers:  []int64{20032},
		},
		{
			name:     "دو بچسبه به هزار",
			input:    "دو هزار و سی و دو تا تسک دارم",
			expected: "2032 تا تسک دارم",
			numbers:  []int64{2032},
		},
		{
			name:     "سیصد بچسبه به هزار",
			input:    "سیصد هزار و سی و دو تا تسک دارم",
			expected: "300032 تا تسک دارم",
			numbers:  []int64{300032},
		},
		{
			name:     "چهار صد با اسپیس",
			input:    "چهار صد و ده",
			expected: "410",
			numbers:  []int64{410},
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
