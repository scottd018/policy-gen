package aws

import (
	"reflect"
	"testing"

	"github.com/scottd018/go-utils/pkg/pointers"

	"github.com/scottd018/policy-gen/internal/pkg/aws/conditions"
	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
)

func TestPolicyDocumentGenerator_ToPolicyMarkerMap(t *testing.T) {
	t.Parallel()

	type fields struct {
		Directory *files.Directory
	}

	type args struct {
		markers []policy.Marker
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    policy.MarkerMap
		wantErr bool
	}{
		{
			name: "ensure incompatible marker returns an error",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					policy.NewFakeMarker(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure markers without defaulted fields set return appropriately",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					&Marker{
						Name:   pointers.String("test"),
						Action: pointers.String("ec2:DescribeVpcs"),
					},
				},
			},
			want: policy.MarkerMap{
				"test/test.json": []policy.Marker{
					&Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:DescribeVpcs"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
						Id:       pointers.String(defaultStatementID),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ensure markers with defaulted fields set return appropriately",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					&Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:DescribeVpcs"),
						Effect:   pointers.String("Deny"),
						Resource: pointers.String("fake"),
						Id:       pointers.String("StatementID"),
					},
				},
			},
			want: policy.MarkerMap{
				"test/test.json": []policy.Marker{
					&Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:DescribeVpcs"),
						Effect:   pointers.String("Deny"),
						Resource: pointers.String("fake"),
						Id:       pointers.String("StatementID"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ensure multiple markers in the same file with same id return appropriately",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					&Marker{
						Name:   pointers.String("test"),
						Action: pointers.String("ec2:DescribeVpcs"),
						Id:     pointers.String("test"),
					},
					&Marker{
						Name:   pointers.String("test"),
						Action: pointers.String("iam:DescribePolicy"),
						Id:     pointers.String("test"),
					},
				},
			},
			want: policy.MarkerMap{
				"test/test.json": []policy.Marker{
					&Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:DescribeVpcs"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
						Id:       pointers.String("test"),
					},
					&Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("iam:DescribePolicy"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
						Id:       pointers.String("test"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ensure multiple markers in different files with same id return appropriately",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					&Marker{
						Name:   pointers.String("test1"),
						Action: pointers.String("ec2:DescribeVpcs"),
						Id:     pointers.String("test"),
					},
					&Marker{
						Name:   pointers.String("test2"),
						Action: pointers.String("iam:DescribePolicy"),
						Id:     pointers.String("test"),
					},
				},
			},
			want: policy.MarkerMap{
				"test/test1.json": []policy.Marker{
					&Marker{
						Name:     pointers.String("test1"),
						Action:   pointers.String("ec2:DescribeVpcs"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
						Id:       pointers.String("test"),
					},
				},
				"test/test2.json": []policy.Marker{
					&Marker{
						Name:     pointers.String("test2"),
						Action:   pointers.String("iam:DescribePolicy"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
						Id:       pointers.String("test"),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			generator := &PolicyDocumentGenerator{
				Directory: tt.fields.Directory,
			}

			got, err := generator.ToPolicyMarkerMap(tt.args.markers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyDocumentGenerator.ToPolicyMarkerMap() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyDocumentGenerator.ToPolicyMarkerMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyDocumentGenerator_ToDocument(t *testing.T) {
	t.Parallel()

	type fields struct {
		Directory *files.Directory
	}

	type args struct {
		markers []policy.Marker
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    policy.Document
		wantErr bool
	}{
		{
			name: "ensure mismatched marker names return an error",
			args: args{
				markers: []policy.Marker{
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test2"),
						Action:   pointers.String("s3:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure multiple with multiple different ids return appropriately",
			fields: fields{
				Directory: &files.Directory{
					Path: "test",
				},
			},
			args: args{
				markers: []policy.Marker{
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("ec2:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("s3:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("sts:*"),
						Effect:   pointers.String(ValidEffectDeny),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("sts:*"),
						Effect:   pointers.String(ValidEffectDeny),
						Resource: pointers.String("thisisfake"),
					},
					&Marker{
						Id:       pointers.String("test"),
						Name:     pointers.String("test"),
						Action:   pointers.String("iam:*"),
						Effect:   pointers.String(ValidEffectDeny),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test1"),
						Name:     pointers.String("test"),
						Action:   pointers.String("kms:*"),
						Effect:   pointers.String(ValidEffectDeny),
						Resource: pointers.String(defaultStatementResource),
					},
					&Marker{
						Id:       pointers.String("test1"),
						Name:     pointers.String("test"),
						Action:   pointers.String("kms:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String("thisisfake"),
					},
					&Marker{
						Id:       pointers.String("test1"),
						Name:     pointers.String("test"),
						Action:   pointers.String("kms:*"),
						Effect:   pointers.String(defaultStatementEffect),
						Resource: pointers.String("thisisfake2"),
					},
					&Marker{
						Id:                pointers.String("test1"),
						Name:              pointers.String("test"),
						Action:            pointers.String("route53:*"),
						Effect:            pointers.String(defaultStatementEffect),
						Resource:          pointers.String(defaultStatementResource),
						ConditionOperator: pointers.String(conditions.StringEqualsOperator),
						ConditionKey:      pointers.String("test"),
						ConditionValue:    pointers.String("test"),
					},
					&Marker{
						Id:                pointers.String("test1"),
						Name:              pointers.String("test"),
						Action:            pointers.String("rds:*"),
						Effect:            pointers.String(defaultStatementEffect),
						Resource:          pointers.String(defaultStatementResource),
						ConditionOperator: pointers.String(conditions.StringEqualsOperator),
						ConditionKey:      pointers.String("test"),
						ConditionValue:    pointers.String("test"),
					},
				},
			},
			want: &PolicyDocument{
				Version: defaultVersion,
				Statements: []Statement{
					{
						SID:    "test",
						Effect: defaultStatementEffect,
						Action: []string{
							"ec2:*",
							"s3:*",
						},
						Resources: []string{defaultStatementResource},
					},
					{
						SID:    "test1",
						Effect: ValidEffectDeny,
						Action: []string{
							"sts:*",
							"iam:*",
							"kms:*",
						},
						Resources: []string{defaultStatementResource},
					},
					{
						SID:    "test2",
						Effect: ValidEffectDeny,
						Action: []string{
							"sts:*",
						},
						Resources: []string{"thisisfake"},
					},
					{
						SID:    "test3",
						Effect: defaultStatementEffect,
						Action: []string{
							"kms:*",
						},
						Resources: []string{"thisisfake"},
					},
					{
						SID:    "test4",
						Effect: defaultStatementEffect,
						Action: []string{
							"kms:*",
						},
						Resources: []string{"thisisfake2"},
					},
					{
						SID:    "test5",
						Effect: defaultStatementEffect,
						Action: []string{
							"route53:*",
							"rds:*",
						},
						Resources: []string{defaultStatementResource},
						Condition: &conditions.Condition{
							StringEquals: conditions.Operator{"test": "test"},
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

			generator := &PolicyDocumentGenerator{
				Directory: tt.fields.Directory,
			}

			got, err := generator.ToDocument(tt.args.markers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyDocumentGenerator.ToDocument() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyDocumentGenerator.ToDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
