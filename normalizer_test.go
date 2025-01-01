package seperno

import (
	"testing"

	"github.com/snapp-incubator/seperno/internal"
)

func TestNormalize_BasicNormalizer(t *testing.T) {
	type args struct {
		input string
		ops   []internal.Options
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
				ops:   []internal.Options{WithConvertHalfSpaceToSpace()},
			},
			want: "اسمان ابی",
		},
		{
			name: "should remove url and outer spaces",
			args: args{
				input: "تست https://example.com",
				ops:   []internal.Options{WithURLRemover(), WithOuterSpaceRemover()},
			},
			want: "تست",
		},
		{
			name: "should remove url",
			args: args{
				input: "تست https://example.com",
				ops:   []internal.Options{WithURLRemover()},
			},
			want: "تست ",
		},
		{
			name: "should combine spaces",
			args: args{
				input: "تست   تست",
				ops:   []internal.Options{WithSpaceCombiner()},
			},
			want: "تست تست",
		},
		{
			name: "should remove end of line character",
			args: args{
				input: "تست.",
				ops:   []internal.Options{WithEndsWithEndOfLineChar()},
			},
			want: "تست",
		},
		{
			name: "should remove punctuations and replace with space",
			args: args{
				input: "سلام,خوبی؟چه خبرا.",
				ops:   []internal.Options{WithNormalizePunctuations(), WithOuterSpaceRemover()},
			},
			want: "سلام خوبی چه خبرا",
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
