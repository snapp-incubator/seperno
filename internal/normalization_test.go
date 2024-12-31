package internal

import (
	"testing"
)

func TestNormalize_BasicNormalizer(t *testing.T) {
	type args struct {
		input                   string
		convertHalfSpaceToSpace bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should convert half space into space",
			args: args{
				input:                   "آ‌س",
				convertHalfSpaceToSpace: true,
			},
			want: "ا س",
		},
		{
			name: "should keep half space",
			args: args{
				input:                   "آ‌س",
				convertHalfSpaceToSpace: false,
			},
			want: "ا‌س",
		},
		{
			name: "should combine spaces",
			args: args{
				input:                   "تست    اسپیس",
				convertHalfSpaceToSpace: true,
			},
			want: "تست اسپیس",
		},
		{
			name: "should remove outer spaces",
			args: args{
				input:                   "  تست   ",
				convertHalfSpaceToSpace: true,
			},
			want: "تست",
		},
		{
			name: "should convert numbers to english",
			args: args{
				input:                   "۶ \u0660 \u0039 \u06F7",
				convertHalfSpaceToSpace: true,
			},
			want: "6 0 9 7",
		},
		{
			name: "should remove combining characters",
			args: args{
				input:                   "ب‍",
				convertHalfSpaceToSpace: true,
			},
			want: "ب",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Normalize{convertHalfSpaceToSpace: tt.args.convertHalfSpaceToSpace}
			if got := n.BasicNormalizer(tt.args.input); got != tt.want {
				t.Errorf("BasicNormalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceMultiSpace(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "space combine",
			args: args{input: "test   test"},
			want: "test test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceMultiSpace(tt.args.input); got != tt.want {
				t.Errorf("replaceMultiSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeURLs(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remove url",
			args: args{input: "https://www.google.com به دنبال"},
			want: " به دنبال",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeURLs(tt.args.input); got != tt.want {
				t.Errorf("removeURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
