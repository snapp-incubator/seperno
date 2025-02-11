package seperno

import (
	"testing"

	"github.com/snapp-incubator/seperno/pkg/options"
)

func TestNormalize_BasicNormalizer(t *testing.T) {
	type args struct {
		input string
		ops   []options.Options
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should convert half space into space",
			args: args{
				input: "آسمان‌آبی",
				ops:   []options.Options{WithConvertHalfSpaceToSpace()},
			},
			want: "اسمان ابی",
		},
		{
			name: "should remove url and outer spaces",
			args: args{
				input: "تست https://example.com",
				ops:   []options.Options{WithURLRemover(), WithOuterSpaceRemover()},
			},
			want: "تست",
		},
		{
			name: "should remove url",
			args: args{
				input: "تست https://example.com",
				ops:   []options.Options{WithURLRemover()},
			},
			want: "تست ",
		},
		{
			name: "should combine spaces",
			args: args{
				input: "تست   تست",
				ops:   []options.Options{WithSpaceCombiner()},
			},
			want: "تست تست",
		},
		{
			name: "should remove end of line character",
			args: args{
				input: "تست.",
				ops:   []options.Options{WithEndsWithEndOfLineChar()},
			},
			want: "تست",
		},
		{
			name: "should remove punctuations and replace with space",
			args: args{
				input: "سلام,خوبی؟چه خبرا.",
				ops:   []options.Options{WithNormalizePunctuations(), WithOuterSpaceRemover()},
			},
			want: "سلام خوبی چه خبرا",
		},
		{
			name: "Should replace number with words",
			args: args{
				input: "خیابان ۱۵ خرداد",
				ops: []options.Options{
					WithIntToWord(),
				},
			},
			want: "خیابان پانزده خرداد",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNormalize(tt.args.ops...).BasicNormalizer(tt.args.input); got != tt.want {
				t.Errorf("spaceNormalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
