package files

import (
	"reflect"
	"testing"
)

func TestNewDirectory(t *testing.T) {
	t.Parallel()

	type args struct {
		path    string
		options []Option
	}

	tests := []struct {
		name    string
		args    args
		want    *Directory
		wantErr bool
	}{
		{
			name: "ensure directory with trailing slash is handled correctly",
			args: args{
				path: "testinput/",
			},
			want:    &Directory{Path: "testinput"},
			wantErr: false,
		},
		{
			name: "ensure root is handled correctly",
			args: args{
				path: "/",
			},
			want:    &Directory{Path: "/"},
			wantErr: false,
		},
		{
			name: "ensure missing directory with option to validate directory fails",
			args: args{
				path:    "fake",
				options: []Option{WithPreExistingDirectory},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure missing directory with option not to validate directory passes",
			args: args{
				path: "fake",
			},
			want:    &Directory{Path: "fake"},
			wantErr: false,
		},
		{
			name: "ensure file fails",
			args: args{
				path:    "directory.go",
				options: []Option{WithPreExistingDirectory},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewDirectory(tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDirectory() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestDirectory_CollectData(t *testing.T) {
// 	t.Parallel()

// 	type fields struct {
// 		Path string
// 	}

// 	tests := []struct {
// 		name       string
// 		fields     fields
// 		wantString string
// 		wantErr    bool
// 	}{
// 		{
// 			name: "ensure data is collected correctly",
// 			fields: fields{
// 				Path: "testinput",
// 			},
// 			wantString: "testdirectory",
// 			wantErr:    false,
// 		},
// 		{
// 			name: "ensure invalid path returns an error",
// 			fields: fields{
// 				Path: "fake",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			dir := &Directory{
// 				Path: tt.fields.Path,
// 			}

// 			got, err := dir.CollectData()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Directory.CollectData() error = %v, wantErr %v", err, tt.wantErr)

// 				return
// 			}

// 			if !reflect.DeepEqual(string(got), tt.wantString) {
// 				t.Errorf("Directory.CollectData() = %s, want %v", got, tt.wantString)
// 			}
// 		})
// 	}
// }
