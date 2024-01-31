package aws

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/scottd018/go-utils/pkg/pointers"
)

func TestMarker_WithDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *Marker
	}{
		{
			name: "ensure default fields are set appropriately",
			want: &Marker{
				Effect:   pointers.String(defaultStatementEffect),
				Resource: pointers.String(defaultStatementResource),
				Id:       pointers.String(defaultStatementID),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			marker := &Marker{}
			marker.WithDefault()
			if !reflect.DeepEqual(marker, tt.want) {
				t.Errorf("WithDefault() = %v, want %v", marker, tt.want)
			}
		})
	}
}

func TestMarkers_PolicyFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		markers Markers
		want    PolicyFiles
	}{
		{
			name: "ensure markers without defaulted fields set return appropriately",
			markers: Markers{
				{
					Name:   pointers.String("test"),
					Action: pointers.String("ec2:DescribeVpcs"),
				},
			},
			want: PolicyFiles{
				"test": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       defaultStatementID,
							Effect:    defaultStatementEffect,
							Resources: []string{defaultStatementResource},
							Action:    []string{"ec2:DescribeVpcs"},
						},
					},
				},
			},
		},
		{
			name: "ensure markers with defaulted fields set return appropriately",
			markers: Markers{
				{
					Name:     pointers.String("test"),
					Action:   pointers.String("ec2:DescribeVpcs"),
					Effect:   pointers.String("Allow"),
					Resource: pointers.String("thisisfake"),
					Id:       pointers.String("test"),
				},
			},
			want: PolicyFiles{
				"test": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{"thisisfake"},
							Action:    []string{"ec2:DescribeVpcs"},
						},
					},
				},
			},
		},
		{
			name: "ensure multiple markers in the same file with same id return appropriately",
			markers: Markers{
				{
					Name:   pointers.String("test"),
					Action: pointers.String("ec2:DescribeVpcs"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test"),
					Action: pointers.String("iam:DescribePolicy"),
					Id:     pointers.String("test"),
				},
			},
			want: PolicyFiles{
				"test": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"ec2:DescribeVpcs",
								"iam:DescribePolicy",
							},
						},
					},
				},
			},
		},
		{
			name: "ensure multiple markers in different files with same id return appropriately",
			markers: Markers{
				{
					Name:   pointers.String("test1"),
					Action: pointers.String("ec2:DescribeVpcs"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test2"),
					Action: pointers.String("iam:DescribePolicy"),
					Id:     pointers.String("test"),
				},
			},
			want: PolicyFiles{
				"test1": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"ec2:DescribeVpcs",
							},
						},
					},
				},
				"test2": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"iam:DescribePolicy",
							},
						},
					},
				},
			},
		},
		{
			name: "ensure multiple markers in different files with multiple different ids return appropriately",
			markers: Markers{
				{
					Name:   pointers.String("test1"),
					Action: pointers.String("ec2:DescribeVpcs"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test1"),
					Action: pointers.String("ec2:Describe*"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test1"),
					Action: pointers.String("iam:*"),
				},
				{
					Name:   pointers.String("test1"),
					Action: pointers.String("route53:*"),
				},
				{
					Name:   pointers.String("test2"),
					Action: pointers.String("ec2:*"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test2"),
					Action: pointers.String("s3:*"),
					Id:     pointers.String("test"),
				},
				{
					Name:   pointers.String("test2"),
					Action: pointers.String("sts:*"),
				},
				{
					Name:     pointers.String("test2"),
					Action:   pointers.String("rds:*"),
					Resource: pointers.String("thisisfake"),
				},
				{
					Name:     pointers.String("test2"),
					Action:   pointers.String("rds:*"),
					Resource: pointers.String("thisisfake2"),
				},
				{
					Name:     pointers.String("test2"),
					Action:   pointers.String("elasticache:*"),
					Resource: pointers.String("thisisfake2"),
					Effect:   pointers.String(ValidEffectDeny),
				},
			},
			want: PolicyFiles{
				"test1": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"ec2:DescribeVpcs",
								"ec2:Describe*",
							},
						},
						{
							SID:       defaultStatementID,
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"iam:*",
								"route53:*",
							},
						},
					},
				},
				"test2": &PolicyDocument{
					Version: defaultVersion,
					Statements: []Statement{
						{
							SID:       "test",
							Effect:    ValidEffectAllow,
							Resources: []string{defaultStatementResource},
							Action: []string{
								"ec2:*",
								"s3:*",
							},
						},
						{
							SID:    defaultStatementID,
							Effect: ValidEffectAllow,
							Resources: []string{
								"thisisfake",
								"thisisfake2",
							},
							Action: []string{
								"sts:*",
								"rds:*",
							},
						},
						{
							SID:    fmt.Sprintf("%s%s", defaultStatementID, ValidEffectDeny),
							Effect: ValidEffectDeny,
							Resources: []string{
								"thisisfake2",
							},
							Action: []string{
								"elasticache:*",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.markers.PolicyFiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Markers.PolicyFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
