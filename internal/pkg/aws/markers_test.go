package aws

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/nukleros/markers/parser"
	"github.com/scottd018/go-utils/pkg/pointers"
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

func TestMarkerResults(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		want    []*parser.Result
		wantErr bool
	}{
		{
			name: "ensure missing path returns an error",
			args: args{
				path: "thisisfake",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure path with no results returns empty",
			args: args{
				path: fmt.Sprintf("%s/testinput/empty.txt", thisFilePath()),
			},
			want:    []*parser.Result{},
			wantErr: false,
		},
		{
			name: "ensure results are parsed properly",
			args: args{
				path: fmt.Sprintf("%s/testinput", thisFilePath()),
			},
			want: []*parser.Result{
				{
					MarkerText: "+policygen:aws:iam:policy:name=test,action=`ec2:DescribeVpcs`\n",
					Object: Marker{
						Name:   pointers.String("test"),
						Action: pointers.String("ec2:DescribeVpcs"),
					},
				},
				{
					MarkerText: "+policygen:aws:iam:policy:name=test,action=`ec2:Describe*`,effect=Deny\n",
					Object: Marker{
						Name:   pointers.String("test"),
						Action: pointers.String("ec2:Describe*"),
						Effect: pointers.String(ValidEffectDeny),
					},
				},
				{
					MarkerText: "+policygen:aws:iam:policy:name=test,action=`iam:Describe*`,effect=Allow,resource=`arn:aws:iam::aws:policy/aws-service-role/*`\n",
					Object: Marker{
						Name:     pointers.String("test"),
						Action:   pointers.String("iam:Describe*"),
						Effect:   pointers.String(ValidEffectAllow),
						Resource: pointers.String("arn:aws:iam::aws:policy/aws-service-role/*"),
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
			got, err := MarkerResults(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkerResults() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkerResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMarkers(t *testing.T) {
	t.Parallel()

	type args struct {
		results []*parser.Result
	}

	tests := []struct {
		name    string
		args    args
		want    Markers
		wantErr bool
	}{
		{
			name: "ensure marker with invalid object type returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: PolicyDocument{},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with missing name returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with blank name returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String(""),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with nil name returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   nil,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with missing action returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:   pointers.String("test"),
							Name: pointers.String("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with blank action returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String(""),
							Name:   pointers.String("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with nil action returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: nil,
							Name:   pointers.String("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with invalid id returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test-123"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with invalid id length returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttestt"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure marker with invalid effect returns error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String("test"),
							Effect: pointers.String("Fail"),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure valid marker with no effect returns without error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String("test"),
							Effect: nil,
						},
					},
				},
			},
			want: Markers{
				{
					Id:     pointers.String("test"),
					Action: pointers.String("ec2:DescribeVpcs"),
					Name:   pointers.String("test"),
					Effect: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "ensure valid marker with returns without error",
			args: args{
				results: []*parser.Result{
					{
						Object: Marker{
							Id:     pointers.String("test"),
							Action: pointers.String("ec2:DescribeVpcs"),
							Name:   pointers.String("test"),
							Effect: pointers.String(defaultStatementEffect),
						},
					},
				},
			},
			want: Markers{
				{
					Id:     pointers.String("test"),
					Action: pointers.String("ec2:DescribeVpcs"),
					Name:   pointers.String("test"),
					Effect: pointers.String(defaultStatementEffect),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindMarkers(tt.args.results)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMarkers() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMarkers() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				Id:       pointers.String(defaultStatementId),
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

func TestMarkers_Process(t *testing.T) {
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
							SID:       defaultStatementId,
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
							SID:       defaultStatementId,
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
							SID:    defaultStatementId,
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
							SID:    fmt.Sprintf("%s%s", defaultStatementId, ValidEffectDeny),
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
			if got := tt.markers.Process(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Markers.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
