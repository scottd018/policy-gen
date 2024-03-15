package utils

import (
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
)

// ConvertToMarker converts an interface to its underlying marker type.
func ConvertToMarker(marker interface{}) (policy.Marker, error) {
	switch t := marker.(type) {
	case aws.Marker:
		return &t, nil
	case *aws.Marker:
		return t, nil
	default:
		return nil, fmt.Errorf("invalid marker type: [%T]", t)
	}
}
