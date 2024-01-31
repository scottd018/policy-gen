package aws

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/nukleros/markers/parser"
	"github.com/rs/zerolog"
	"github.com/scottd018/go-utils/pkg/pointers"

	"github.com/scottd018/policy-gen/internal/pkg/input"
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

func Test_markerProcessor_Parse(t *testing.T) {
	t.Parallel()

	type fields struct {
		Input input.Input
	}

	tests := []struct {
		name    string
		fields  fields
		want    []*parser.Result
		wantErr bool
	}{
		{
			name: "ensure missing path returns an error",
			fields: fields{
				Input: input.Input{
					InputPath: "thisisfake",
					Debug:     true,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ensure path with no results returns empty",
			fields: fields{
				Input: input.Input{InputPath: fmt.Sprintf("%s/testinput/empty.txt", thisFilePath())},
			},
			want:    []*parser.Result{},
			wantErr: false,
		},
		{
			name: "ensure results are parsed properly",
			fields: fields{
				Input: input.Input{InputPath: fmt.Sprintf("%s/testinput", thisFilePath())},
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

			processor := NewMarkerProcessor(tt.fields.Input)

			got, err := processor.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("markerProcessor.Parse() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markerProcessor.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_markerProcessor_FindMarkers(t *testing.T) {
	t.Parallel()

	type fields struct {
		Input input.Input
		Log   zerolog.Logger
	}

	type args struct {
		results []*parser.Result
	}

	tests := []struct {
		name    string
		fields  fields
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
			processor := &MarkerProcessor{}

			got, err := processor.FindMarkers(tt.args.results)
			if (err != nil) != tt.wantErr {
				t.Errorf("markerProcessor.FindMarkers() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markerProcessor.FindMarkers() = %v, want %v", got, tt.want)
			}
		})
	}
}
