package files

import (
	"io/fs"
	"os"
	"reflect"
	"testing"
)

func TestNewJSONFile(t *testing.T) {
	t.Parallel()

	type args struct {
		path    string
		options []Option
	}

	tests := []struct {
		name    string
		args    args
		want    *JSON
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
				path: "testoutput/test.json.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure non-existing file returns appropriately",
			args: args{
				path: "testoutput/fake.json",
			},
			want: &JSON{
				Directory: &Directory{
					Path: "testoutput",
				},
				File:      "fake",
				Extension: ExtensionJSON,
			},
			wantErr: false,
		},
		{
			name: "ensure existing file returns appropriately",
			args: args{
				path: "testinput/existing.json",
			},
			want: &JSON{
				Directory: &Directory{
					Path: "testinput",
				},
				File:      "existing",
				Extension: ExtensionJSON,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewJSONFile(tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_Write(t *testing.T) {
	t.Parallel()

	type Test struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type fields struct {
		Directory *Directory
		File      string
		Extension string
	}

	type args struct {
		object      any
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
				Extension: ExtensionJSON,
			},
			args: args{
				object: &Test{
					ID:   "test",
					Name: "test",
				},
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
				Extension: ExtensionJSON,
			},
			args: args{
				object: &Test{
					ID:   "test",
					Name: "test",
				},
				permissions: 0600,
				options:     []Option{WithOverwrite},
			},
			wantErr: false,
			wantData: `{
	"id": "test",
	"name": "test"
}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file := &JSON{
				Directory: tt.fields.Directory,
				File:      tt.fields.File,
				Extension: tt.fields.Extension,
			}

			if err := file.Write(tt.args.object, tt.args.permissions, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("JSON.Write() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			data, _ := os.ReadFile(file.Path())
			got := string(data)
			if !reflect.DeepEqual(got, tt.wantData) {
				t.Errorf("JSON.Write() = %s, wantData %s", got, tt.wantData)
			}
		})
	}
}
