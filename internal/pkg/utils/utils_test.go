package utils

import (
	"reflect"
	"testing"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
)

func TestConvertToMarker(t *testing.T) {
	t.Parallel()

	type args struct {
		marker interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    policy.Marker
		wantErr bool
	}{
		{
			name: "ensure aws.Marker returns correctly",
			args: args{
				marker: aws.Marker{},
			},
			want:    &aws.Marker{},
			wantErr: false,
		},
		{
			name: "ensure aws.Marker pointer returns correctly",
			args: args{
				marker: &aws.Marker{},
			},
			want:    &aws.Marker{},
			wantErr: false,
		},
		{
			name: "ensure invalid object returns with error",
			args: args{
				marker: policy.NewFakeMarker(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ConvertToMarker(tt.args.marker)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToMarker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToMarker() = %v, want %v", got, tt.want)
			}
		})
	}
}
