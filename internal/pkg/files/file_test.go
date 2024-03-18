package files

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func thisFilePath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "."
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		return "."
	}

	return filepath.Dir(absPath)
}

func thisFilePathFor(file string) string {
	dirPath := thisFilePath()

	return fmt.Sprintf("%s/%s", dirPath, file)
}

func TestNewFile(t *testing.T) {
	t.Parallel()

	type args struct {
		path    string
		options []Option
	}

	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
	}{
		{
			name: "ensure missing directory when existing directory option is requested returns an error",
			args: args{
				path:    thisFilePathFor("fake/path"),
				options: []Option{WithPreExistingDirectory},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure pointing at a directory for a new file returns an error",
			args: args{
				path:    thisFilePath(),
				options: []Option{WithPreExistingDirectory},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure file object is created correctly for a non-existing file",
			args: args{
				path: thisFilePathFor("test/output/test.txt"),
			},
			want: &File{
				Directory: &Directory{
					Path: thisFilePathFor("test/output"),
				},
				File:    thisFilePathFor("test/output/test.txt"),
				Content: nil,
			},
			wantErr: false,
		},
		{
			name: "ensure file object is created correctly for a pre-existing file",
			args: args{
				path: thisFilePathFor("test/input/existing.txt"),
			},
			want: &File{
				Directory: &Directory{
					Path: thisFilePathFor("test/input"),
				},
				File:    thisFilePathFor("test/input/existing.txt"),
				Content: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFile(tt.args.path, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Write(t *testing.T) {
	t.Parallel()

	type fields struct {
		Directory *Directory
		File      string
		Content   []byte
	}

	type args struct {
		permissions fs.FileMode
		options     []Option
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ensure file with nil content returns an error",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/output"),
				},
				File: thisFilePathFor("test/output/nil-content.txt"),
			},
			wantErr: true,
		},
		{
			name: "ensure file with empty content returns an error",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/output"),
				},
				File:    thisFilePathFor("test/output/nil-content.txt"),
				Content: []byte{},
			},
			wantErr: true,
		},
		{
			name: "ensure writing to a non-existent directory returns an error",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/output/fake"),
				},
				File:    thisFilePathFor(fmt.Sprintf("test/output/fake/%d.txt", time.Now().Unix())),
				Content: []byte("test"),
			},
			args: args{
				permissions: ModePolicyFile,
			},
			wantErr: true,
		},
		//
		// we will simply use the unix seconds to create the file to ensure it does not exist.  the
		// make target should clear this out.
		//
		{
			name: "write a file that does not exist",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/output"),
				},
				File:    thisFilePathFor(fmt.Sprintf("test/output/%d.txt", time.Now().Unix())),
				Content: []byte("write a file that does not exist"),
			},
			args: args{
				permissions: ModePolicyFile,
			},
			wantErr: false,
		},

		//
		// we will use an input file to test here since these are stored in git and our output files
		// are not.  this ensures we are testing against a file guaranteed to be existing.
		//
		{
			name: "ensure existing file without overwrite requested returns an error",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/input"),
				},
				File:    thisFilePathFor("test/input/existing.json"),
				Content: []byte("test"),
			},
			wantErr: true,
		},
		{
			name: "write a file that does exist with overwrite requested",
			fields: fields{
				Directory: &Directory{
					Path: thisFilePathFor("test/input"),
				},
				File:    thisFilePathFor("test/input/existing.txt"),
				Content: []byte("write a file that does exist"),
			},
			args: args{
				permissions: ModePolicyFile,
				options:     []Option{WithOverwrite},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file := &File{
				Directory: tt.fields.Directory,
				File:      tt.fields.File,
				Content:   tt.fields.Content,
			}

			if err := file.Write(tt.args.permissions, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("File.Write() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				content, _ := os.ReadFile(file.File)
				if !bytes.Equal(content, tt.fields.Content) {
					t.Errorf("File.Write() content = %s, wantContent %s", content, tt.fields.Content)
				}
			}
		})
	}
}
