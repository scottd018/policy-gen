package aws

import (
	"reflect"
	"testing"

	"github.com/scottd018/go-utils/pkg/pointers"
)

func TestMarker_Definition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "ensure marker returns appropriately",
			want: "+policy-gen:aws:iam:policy",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{}
			if got := marker.Definition(); got != tt.want {
				t.Errorf("Marker.Definition() = %v, want %v", got, tt.want)
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

func TestMarker_GetName(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ensure marker with nil name returns appropriately",
			fields: fields{
				Name: nil,
			},
			want: "",
		},
		{
			name: "ensure marker with non-nil name returns appropriately",
			fields: fields{
				Name: pointers.String("test"),
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Name: tt.fields.Name}
			if got := marker.GetName(); got != tt.want {
				t.Errorf("Marker.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarker_EffectColumn(t *testing.T) {
	t.Parallel()

	type fields struct {
		Effect *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ensure marker with nil effect returns appropriately",
			fields: fields{
				Effect: nil,
			},
			want: defaultStatementEffect,
		},
		{
			name: "ensure marker with non-nil effect returns appropriately",
			fields: fields{
				Effect: pointers.String(defaultStatementEffect),
			},
			want: defaultStatementEffect,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Effect: tt.fields.Effect}
			if got := marker.EffectColumn(); got != tt.want {
				t.Errorf("Marker.EffectColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarker_PermissionColumn(t *testing.T) {
	t.Parallel()

	type fields struct {
		Action *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ensure marker with nil action returns appropriately",
			fields: fields{
				Action: nil,
			},
			want: "",
		},
		{
			name: "ensure marker with non-nil action returns appropriately",
			fields: fields{
				Action: pointers.String("ec2:DescribeVpcs"),
			},
			want: "ec2:DescribeVpcs",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Action: tt.fields.Action}
			if got := marker.PermissionColumn(); got != tt.want {
				t.Errorf("Marker.PermissionColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarker_ResourceColumn(t *testing.T) {
	t.Parallel()

	type fields struct {
		Resource *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ensure marker with nil resource returns appropriately",
			fields: fields{
				Resource: nil,
			},
			want: defaultStatementResource,
		},
		{
			name: "ensure marker with non-nil resource returns appropriately",
			fields: fields{
				Resource: pointers.String(defaultStatementResource),
			},
			want: defaultStatementResource,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Resource: tt.fields.Resource}
			if got := marker.ResourceColumn(); got != tt.want {
				t.Errorf("Marker.ResourceColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarker_ReasonColumn(t *testing.T) {
	t.Parallel()

	type fields struct {
		Reason *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ensure marker with nil reason returns appropriately",
			fields: fields{
				Reason: nil,
			},
			want: "",
		},
		{
			name: "ensure marker with non-nil reason returns appropriately",
			fields: fields{
				Reason: pointers.String("because i said so"),
			},
			want: "because i said so",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Reason: tt.fields.Reason}
			if got := marker.ReasonColumn(); got != tt.want {
				t.Errorf("Marker.ReasonColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarker_AdjustID(t *testing.T) {
	t.Parallel()

	//nolint:revive,stylecheck
	type fields struct {
		Id *string
	}

	tests := []struct {
		name   string
		fields fields
		want   *Marker
	}{
		{
			name: "ensure marker without value returns 1",
			fields: fields{
				Id: pointers.String("Default"),
			},
			want: &Marker{
				Id: pointers.String("Default1"),
			},
		},
		{
			name: "ensure marker with existing integer adds 1",
			fields: fields{
				Id: pointers.String("Default1"),
			},
			want: &Marker{
				Id: pointers.String("Default2"),
			},
		},
		{
			name: "ensure marker with trailing integers adds 1",
			fields: fields{
				Id: pointers.String("Default123"),
			},
			want: &Marker{
				Id: pointers.String("Default124"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{Id: tt.fields.Id}
			marker.AdjustID()
		})
	}
}

func TestMarker_Validate(t *testing.T) {
	t.Parallel()

	//nolint:revive,stylecheck
	type fields struct {
		Name     *string
		Id       *string
		Action   *string
		Effect   *string
		Resource *string
		Reason   *string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ensure marker with nil name returns an error",
			fields: fields{
				Name:     nil,
				Id:       pointers.String(defaultStatementID),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure marker with empty name returns an error",
			fields: fields{
				Name:     pointers.String(""),
				Id:       pointers.String(defaultStatementID),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure marker with nil action returns an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String(defaultStatementID),
				Action:   nil,
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure marker with empty action returns an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String(defaultStatementID),
				Action:   pointers.String(""),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure name with invalid characters returns an error",
			fields: fields{
				Name:     pointers.String("test-name"),
				Id:       pointers.String(defaultStatementID),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure name with too many characters returns an error",
			fields: fields{
				Name:     pointers.String("testnametestnametestnametestnametestnametestnametestnametestnamet"),
				Id:       pointers.String(defaultStatementID),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure marker with invalid id returns an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String("Test-Id"),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure marker with invalid effect returns an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String("TestId"),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String("Fake"),
				Resource: pointers.String(defaultStatementResource),
				Reason:   pointers.String("test"),
			},
			wantErr: true,
		},
		{
			name: "ensure valid marker without effect returns without an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String("TestId"),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   nil,
				Resource: pointers.String("*"),
				Reason:   pointers.String("test"),
			},
			wantErr: false,
		},
		{
			name: "ensure valid marker returns without an error",
			fields: fields{
				Name:     pointers.String("test"),
				Id:       pointers.String("TestId"),
				Action:   pointers.String("ec2:DescribeVpcs"),
				Effect:   pointers.String(ValidEffectAllow),
				Resource: pointers.String("*"),
				Reason:   pointers.String("test"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			marker := &Marker{
				Name:     tt.fields.Name,
				Id:       tt.fields.Id,
				Action:   tt.fields.Action,
				Effect:   tt.fields.Effect,
				Resource: tt.fields.Resource,
				Reason:   tt.fields.Reason,
			}
			if err := marker.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Marker.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
