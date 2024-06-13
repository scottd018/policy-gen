package aws

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/scottd018/go-utils/pkg/pointers"

	"github.com/scottd018/policy-gen/internal/pkg/aws/conditions"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
)

var (
	ErrMarkerMissingName                     = errors.New("marker missing name field")
	ErrMarkerMissingAction                   = errors.New("marker missing action field")
	ErrMarkerInvalidEffect                   = errors.New("invalid marker effect")
	ErrMarkerInvalidStatementID              = errors.New("invalid statement id - must contain a-z, A-Z, 0-9 and limited to 64 characters")
	ErrMarkerInvalidConditionMissingKey      = errors.New("condition key is missing")
	ErrMarkerInvalidConditionMissingValue    = errors.New("condition value is missing")
	ErrMarkerInvalidConditionMissingOperator = errors.New("condition operator is missing")
	ErrMarkerInvalidConditionOperator        = errors.New("invalid condition operator")
	ErrMarkerInvalidName                     = errors.New(
		"invalid name - must contain only lowercase alphanumeric characters with underscores or dashes and is limited to 64 characters",
	)
)

const (
	awsMarkerDefinition = "aws:iam:policy"
	nameRegex           = "^[a-z0-9_-]{1,64}$"

	ValidEffectAllow = "Allow"
	ValidEffectDeny  = "Deny"
)

// we must not lint Id for ID here as the markers package incorrectly parses a
// capitalized ID.
//
//nolint:revive,stylecheck
type Marker struct {
	Name     *string
	Id       *string
	Action   *string
	Effect   *string
	Resource *string
	Reason   *string

	// conditions
	ConditionOperator *string
	ConditionKey      *string
	ConditionValue    *string
}

// MarkerDefinition returns the marker definition for an AWS IAM policy marker.
func MarkerDefinition() string {
	return fmt.Sprintf("%s%s:%s", policy.MarkerPrefixStart, policy.MarkerPrefixString, awsMarkerDefinition)
}

// Definition returns the marker definition for an AWS IAM policy marker.  It is used
// as a way to return the definition as part of the policymarkers.Marker interface.
func (marker *Marker) Definition() string {
	return MarkerDefinition()
}

// Validate validates that a marker is valid.  It is used to satisfy the policymarkers.Marker
// interface.
func (marker *Marker) Validate() error {
	// ensure required markers are set
	for err, markerField := range map[error]*string{
		ErrMarkerMissingName:   marker.Name,
		ErrMarkerMissingAction: marker.Action,
	} {
		if markerField == nil {
			return err
		}

		if *markerField == "" {
			return err
		}
	}

	// ensure the name only contains lowercase characters with underscores/dashes and is limited
	// to 64 characters in length.  this is because we are generating file names based upon
	// the policy name and grouping those with like names together into separate policy
	// files.
	nameCheck := regexp.MustCompile(nameRegex)
	if !nameCheck.MatchString(*marker.Name) {
		return fmt.Errorf("%w - [%s]", ErrMarkerInvalidName, *marker.Name)
	}

	// ensure the sid is valid if specified
	if marker.Id != nil {
		statementIDCheck := regexp.MustCompile(statementIDRegex)
		if !statementIDCheck.MatchString(*marker.Id) {
			return fmt.Errorf("%w - [%s]", ErrMarkerInvalidStatementID, *marker.Id)
		}
	}

	// ensure effect is valid
	if marker.Effect != nil {
		if *marker.Effect != ValidEffectAllow && *marker.Effect != ValidEffectDeny {
			return fmt.Errorf("%w [%s]", ErrMarkerInvalidEffect, *marker.Effect)
		}
	}

	// ensure the condition is valid
	if err := marker.ValidateCondition(); err != nil {
		return fmt.Errorf("invalid condition specified - %w", err)
	}

	return nil
}

// WithDefault sets a marker with its default values.  It is used to satisfy the policymarkers.Marker
// interface.
func (marker *Marker) WithDefault() {
	// add the effect if we specified one otherwise default to allow
	if marker.Effect == nil {
		marker.Effect = &defaultStatementEffect
	}

	// add the resource if we specified one otherwise default to all
	if marker.Resource == nil {
		marker.Resource = &defaultStatementResource
	}

	// add the id if we specified one otherwise use the default statement id
	if marker.Id == nil {
		marker.Id = &defaultStatementID
	}
}

// GetName returns the name of the marker.  It is used to satisfy the policymarkers.Marker
// interface.
func (marker *Marker) GetName() string {
	if marker.Name == nil {
		return ""
	}

	return *marker.Name
}

// ToStatement converts a marker to an AWS IAM policy statement.
func (marker Marker) ToStatement() Statement {
	return Statement{
		Action:    []string{*marker.Action},
		Effect:    *marker.Effect,
		Resources: []string{*marker.Resource},
		SID:       *marker.Id,
		Condition: marker.Condition(),
	}
}

// Condition returns the condition for a given marker.
func (marker *Marker) Condition() *conditions.Condition {
	if marker.HasConditionKey() && marker.HasConditionValue() && marker.HasConditionOperator() {
		return conditions.NewCondition(*marker.ConditionKey, *marker.ConditionValue, *marker.ConditionOperator)
	}

	return nil
}

// EffectColumn returns the effect for the marker.  It is used to satisfy
// the docs.Row interface.
func (marker *Marker) EffectColumn() string {
	if marker.Effect == nil {
		return defaultStatementEffect
	}

	return *marker.Effect
}

// PermissionColumn returns the permission (action) for the marker.  It is used to satisfy
// the docs.Row interface.
func (marker *Marker) PermissionColumn() string {
	if marker.Action == nil {
		return ""
	}

	return *marker.Action
}

// ResourceColumn returns the applicable resource that this permission is valid for.  It
// is used to satisfy the docs.Row interface.
func (marker *Marker) ResourceColumn() string {
	if marker.Resource == nil {
		return defaultStatementResource
	}

	return *marker.Resource
}

// ReasonColumn returns the reason for the permission.  It is used to satisfy the docs.Row
// interface.
func (marker *Marker) ReasonColumn() string {
	if marker.Reason == nil {
		return ""
	}

	return *marker.Reason
}

// ConditionColumn returns the conditions for the permission, represented in a comma-separated list.
// It is used to satisfy the docs.Row interface.  For now, we only care about tag conditions, but this may expand
// if more conditions are added in the future.
func (marker *Marker) ConditionColumn() string {
	condition := marker.Condition()

	if condition != nil {
		return condition.String()
	}

	return ""
}

// AdjustID adjusts an ID for situations where a conflict arises.
func (marker *Marker) AdjustID() {
	// this collects the suffix integers on the current id
	var suffix string

	// the prefix is considered to be the non-integer prefix
	prefix := *marker.Id

	// loop until we do not find a trailing integer
	for i := len(prefix) - 1; i >= 0; i-- {
		id := *marker.Id

		// break the loop if we found a non-integer
		if !unicode.IsDigit(rune(id[i])) {
			prefix = prefix[:(i + 1)]

			break
		}

		// collect the integer and store it
		suffix = fmt.Sprintf("%s%s", string(id[i]), suffix)
	}

	// add to the collected value
	value, _ := strconv.Atoi(suffix)
	value++

	// set the id
	marker.Id = pointers.String(fmt.Sprintf("%s%d", prefix, value))
}

// HasConditionOperator returns whether or not a marker has a condition operator.
func (marker *Marker) HasConditionOperator() bool {
	return hasStringValue(marker.ConditionOperator)
}

// HasConditionKey returns whether or not a marker has a condition key.
func (marker *Marker) HasConditionKey() bool {
	return hasStringValue(marker.ConditionKey)
}

// HasConditionValue returns whether or not a marker has a condition value.
func (marker *Marker) HasConditionValue() bool {
	return hasStringValue(marker.ConditionValue)
}

// ValidateCondition returns whether or not a marker has a valid condition.
func (marker *Marker) ValidateCondition() error {
	hasConditionKey := marker.HasConditionKey()
	hasConditionValue := marker.HasConditionValue()
	hasConditionOperator := marker.HasConditionOperator()

	if (hasConditionKey == hasConditionValue) && (hasConditionValue == hasConditionOperator) {
		if hasConditionOperator && conditions.ToOperatorString(*marker.ConditionOperator) == "" {
			return fmt.Errorf("found operator [%s] - %w", *marker.ConditionOperator, ErrMarkerInvalidConditionOperator)
		}

		return nil
	}

	messages := []string{"conditionKey, conditionValue, and conditionOperator are mutually inclusive options"}

	if !hasConditionKey {
		messages = append(messages, ErrMarkerInvalidConditionMissingKey.Error())
	}

	if !hasConditionValue {
		messages = append(messages, ErrMarkerInvalidConditionMissingValue.Error())
	}

	if !hasConditionOperator {
		messages = append(messages, ErrMarkerInvalidConditionMissingOperator.Error())
	} else if hasConditionOperator && conditions.ToOperatorString(*marker.ConditionOperator) == "" {
		messages = append(messages, fmt.Sprintf("found operator [%s] - %s", *marker.ConditionOperator, ErrMarkerInvalidConditionOperator.Error()))
	}

	return fmt.Errorf(strings.Join(messages, " : "))
}

// hasStringValue returns whether or not a marker has a string value given a pointer to a string.
func hasStringValue(value *string) bool {
	if value == nil {
		return false
	}

	return *value != ""
}
