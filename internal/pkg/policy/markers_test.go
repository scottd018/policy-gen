package policy

import (
	"reflect"
	"testing"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

func TestToPolicyFiles(t *testing.T) {
	t.Parallel()

	type args struct {
		markers   []Marker
		generator FileGenerator
	}

	tests := []struct {
		name    string
		args    args
		want    []*files.File
		wantErr bool
	}{
		// 		{
		// 			name: "ensure markers without defaulted fields set return appropriately",
		// 			markers: []Marker{
		// 				{
		// 					Name:   pointers.String("test"),
		// 					Action: pointers.String("ec2:DescribeVpcs"),
		// 				},
		// 			},
		// 			want: PolicyFiles{
		// 				"test": &PolicyDocument{
		// 					Version: defaultVersion,
		// 					Statements: []Statement{
		// 						{
		// 							SID:       defaultStatementID,
		// 							Effect:    defaultStatementEffect,
		// 							Resources: []string{defaultStatementResource},
		// 							Action:    []string{"ec2:DescribeVpcs"},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ToPolicyFiles(tt.args.markers, tt.args.generator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToPolicyFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPolicyFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
