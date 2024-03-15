package files

import (
	"fmt"
	"reflect"
	"sort"
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
				path: "test/input/",
			},
			want:    &Directory{Path: "test/input"},
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

func TestDirectory_ListFilePaths(t *testing.T) {
	t.Parallel()

	inputPath := thisFilePathFor("test/input")

	type fields struct {
		Path string
	}

	type args struct {
		recursive bool
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ensure file paths return correctly when recursive is requested",
			fields: fields{
				Path: inputPath,
			},
			args: args{
				recursive: true,
			},
			want: []string{
				fmt.Sprintf("%s/%s", inputPath, "existing.json"),
				fmt.Sprintf("%s/%s", inputPath, "existing.md"),
				fmt.Sprintf("%s/%s", inputPath, "existing.txt"),
				fmt.Sprintf("%s/%s/%s", inputPath, "directory", "test.txt"),
				fmt.Sprintf("%s/%s/%s/%s", inputPath, "directory", "recursive", "test.txt"),
			},
			wantErr: false,
		},
		{
			name: "ensure file paths return correctly when recursive is not requested",
			fields: fields{
				Path: inputPath,
			},
			args: args{
				recursive: false,
			},
			want: []string{
				fmt.Sprintf("%s/%s", inputPath, "existing.json"),
				fmt.Sprintf("%s/%s", inputPath, "existing.md"),
				fmt.Sprintf("%s/%s", inputPath, "existing.txt"),
			},
			wantErr: false,
		},
		{
			name: "ensure missing directory returns an error",
			fields: fields{
				Path: "missing",
			},
			want:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := &Directory{
				Path: tt.fields.Path,
			}

			got, err := dir.ListFilePaths(tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("Directory.ListFilePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Directory.ListFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
