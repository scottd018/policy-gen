package files

import "testing"

func Test_hasOption(t *testing.T) {
	t.Parallel()

	type args struct {
		option  Option
		options []Option
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ensure option list that has an option returns true",
			args: args{
				option:  WithOverwrite,
				options: []Option{WithPreExistingDirectory, WithOverwrite},
			},
			want: true,
		},
		{
			name: "ensure option list that does not have an option returns false",
			args: args{
				option:  WithOverwrite,
				options: []Option{WithPreExistingDirectory},
			},
			want: false,
		},
		{
			name: "ensure empty option list returns false",
			args: args{
				option:  WithOverwrite,
				options: []Option{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := hasOption(tt.args.option, tt.args.options...); got != tt.want {
				t.Errorf("hasOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
