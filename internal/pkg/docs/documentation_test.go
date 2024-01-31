package docs

// import "testing"

// func TestFilename(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		path string
// 		name string
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "ensure file with extra slash returns appropriately",
// 			args: args{
// 				path: "test/",
// 				name: "test",
// 			},
// 			want: "test/test.md",
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			if got := Filename(tt.args.path, tt.args.name); got != tt.want {
// 				t.Errorf("Filename() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
