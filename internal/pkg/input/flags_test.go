package input

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/processor"
)

func TestFlags_ToProcessorConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		flags        Flags
		want         *processor.Config
		wantErr      bool
		overrideFunc func(flags *Flags)
	}{
		{
			name:    "ensure missing input path input returns an error",
			flags:   NewFlags(),
			want:    nil,
			wantErr: true,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagInputPath].StringValue = ""
			},
		},
		{
			name:    "ensure missing output path input returns an error",
			flags:   NewFlags(),
			want:    nil,
			wantErr: true,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagOutputPath].StringValue = ""
			},
		},
		{
			name:    "ensure missing input path returns an error",
			flags:   NewFlags(),
			want:    nil,
			wantErr: true,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagInputPath].StringValue = "this/path/is/fake/input"
			},
		},
		{
			name:    "ensure missing output path returns an error",
			flags:   NewFlags(),
			want:    nil,
			wantErr: true,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagOutputPath].StringValue = "this/path/is/fake/output"
			},
		},
		{
			name:    "ensure missing documentation file path returns an error",
			flags:   NewFlags(),
			want:    nil,
			wantErr: true,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagDocumentation].StringValue = "this/path/is/fake/README.md"
			},
		},
		{
			name:  "ensure processor config returns correctly",
			flags: NewFlags(),
			want: &processor.Config{
				InputDirectory:  &files.Directory{Path: "."},
				OutputDirectory: &files.Directory{Path: "."},
				DocumentationFile: &files.File{
					Directory: &files.Directory{Path: "."},
					File:      "README.md",
				},
				Force: false,
				Debug: false,
			},
			wantErr: false,
			overrideFunc: func(flags *Flags) {
				f := *flags
				f[FlagInputPath].StringValue = "."
				f[FlagOutputPath].StringValue = "."
				f[FlagDocumentation].StringValue = "README.md"
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.flags.Initialize(&cobra.Command{})
			tt.overrideFunc(&tt.flags)

			got, err := tt.flags.ToProcessorConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Flags.ToProcessorConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flags.ToProcessorConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
