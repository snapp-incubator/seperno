package internal

import "testing"

func TestIntegerToPersian(t *testing.T) {
	type args struct {
		input int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				input: 3,
			},
			want: "سه",
		},
		{
			name: "test 2",
			args: args{
				input: 23,
			},
			want: "بیست و سه",
		},
		{
			name: "test 3",
			args: args{
				input: 123,
			},
			want: "صد و بیست و سه",
		},
		{
			name: "test 4",
			args: args{
				input: 1235,
			},
			want: "یک هزار دویست و سی و پنج",
		},
		{
			name: "test 5",
			args: args{
				input: 12356,
			},
			want: "دوازده هزار سیصد و پنجاه و شش",
		},
		{
			name: "test 5",
			args: args{
				input: 123567,
			},
			want: "صد و بیست و سه هزار پانصد و شصت و هفت",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntegerToPersian(tt.args.input); got != tt.want {
				t.Errorf("IntegerToPersian() = %v, want %v", got, tt.want)
			}
		})
	}
}
