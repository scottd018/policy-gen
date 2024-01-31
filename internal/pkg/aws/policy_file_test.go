package aws

import "testing"

func TestFilename(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
		name string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ensure filename with trailing / returns appropriately",
			args: args{
				path: "/home/test/",
				name: "test",
			},
			want: "/home/test/test.json",
		},
		{
			name: "ensure filename without trailing / returns appropriately",
			args: args{
				path: "/home/test",
				name: "test",
			},
			want: "/home/test/test.json",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Filename(tt.args.path, tt.args.name); got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}
