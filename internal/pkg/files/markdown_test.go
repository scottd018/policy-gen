package files

import (
	"io/fs"
	"os"
	"reflect"
	"testing"
)

func TestNewMarkdownFile(t *testing.T) {
	t.Parallel()

	type args struct {
		path    string
		options []Option
	}

	tests := []struct {
		name    string
		args    args
		want    *Markdown
		wantErr bool
	}{
		{
			name: "ensure missing directory path fails",
			args: args{
				path:    "fake/file.txt",
				options: []Option{WithPreExistingDirectory},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure directory path without file fails",
			args: args{
				path: "/",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure directory path fails",
			args: args{
				path: "testoutput",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure strange file extension fails",
			args: args{
				path: "testoutput/test.md.md",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure non-existing file returns appropriately",
			args: args{
				path: "testoutput/fake.md",
			},
			want: &Markdown{
				Directory: &Directory{
					Path: "testoutput",
				},
				File:      "fake",
				Extension: ExtensionMarkdown,
			},
			wantErr: false,
		},
		{
			name: "ensure existing file returns appropriately",
			args: args{
				path: "testinput/existing.md",
			},
			want: &Markdown{
				Directory: &Directory{
					Path: "testinput",
				},
				File:      "existing",
				Extension: ExtensionMarkdown,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewMarkdownFile(tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMarkdownFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarkdownFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdown_Write(t *testing.T) {
	t.Parallel()

	type fields struct {
		Directory *Directory
		File      string
		Extension string
	}

	type args struct {
		object      []byte
		permissions fs.FileMode
		options     []Option
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantData string
	}{
		{
			name: "ensure existing file without overwrite option fails",
			fields: fields{
				Directory: &Directory{
					Path: "testinput",
				},
				File:      "existing",
				Extension: ExtensionMarkdown,
			},
			args: args{
				object:      []byte("# Test\n\ntest"),
				permissions: 0600,
			},
			wantErr: true,
		},
		{
			name: "ensure data is written appropriately",
			fields: fields{
				Directory: &Directory{
					Path: "testoutput",
				},
				File:      "test",
				Extension: ExtensionMarkdown,
			},
			args: args{
				object:      []byte("# Test\n\ntest"),
				permissions: 0600,
				options:     []Option{WithOverwrite},
			},
			wantErr:  false,
			wantData: "# Test\n\ntest",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file := &Markdown{
				Directory: tt.fields.Directory,
				File:      tt.fields.File,
				Extension: tt.fields.Extension,
			}

			if err := file.Write(tt.args.object, tt.args.permissions, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Markdown.Write() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			data, _ := os.ReadFile(file.Path())
			got := string(data)
			if !reflect.DeepEqual(got, tt.wantData) {
				t.Errorf("Markdown.Write() = %s, wantData %s", got, tt.wantData)
			}
		})
	}
}
