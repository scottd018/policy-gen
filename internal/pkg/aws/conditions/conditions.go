package conditions

import (
	"encoding/json"
)

// Condition represent a condition statement.
//
//nolint:revive,stylecheck
type Condition struct {
	// string condition operators
	StringEquals              Operator `json:"StringEquals,omitempty"`
	StringNotEquals           Operator `json:"StringNotEquals,omitempty"`
	StringEqualsIgnoreCase    Operator `json:"StringEqualsIgnoreCase,omitempty"`
	StringNotEqualsIgnoreCase Operator `json:"StringNotEqualsIgnoreCase,omitempty"`
	StringLike                Operator `json:"StringLike,omitempty"`
	StringNotLike             Operator `json:"StringNotLike,omitempty"`

	// numeric condition operators
	NumericEquals            Operator `json:"NumericEquals,omitempty"`
	NumericNotEquals         Operator `json:"NumericNotEquals,omitempty"`
	NumericLessThan          Operator `json:"NumericLessThan,omitempty"`
	NumericLessThanEquals    Operator `json:"NumericLessThanEquals,omitempty"`
	NumericGreaterThan       Operator `json:"NumericGreaterThan,omitempty"`
	NumericGreaterThanEquals Operator `json:"NumericGreaterThanEquals,omitempty"`

	// date condition operators
	DateEquals            Operator `json:"DateEquals,omitempty"`
	DateNotEquals         Operator `json:"DateNotEquals,omitempty"`
	DateLessThan          Operator `json:"DateLessThan,omitempty"`
	DateLessThanEquals    Operator `json:"DateLessThanEquals,omitempty"`
	DateGreaterThan       Operator `json:"DateGreaterThan,omitempty"`
	DateGreaterThanEquals Operator `json:"DateGreaterThanEquals,omitempty"`

	// bool condition operators
	Bool Operator `json:"Bool,omitempty"`

	// binary condition operators
	BinaryEquals Operator `json:"BinaryEquals,omitempty"`

	// ip address condition operators
	IpAddress    Operator `json:"IpAddress,omitempty"`
	NotIpAddress Operator `json:"NotIpAddress,omitempty"`

	// arn condition operators
	ArnEquals    Operator `json:"ArnEquals,omitempty"`
	ArnNotEquals Operator `json:"ArnNotEquals,omitempty"`
	ArnLike      Operator `json:"ArnLike,omitempty"`
	ArnNotLike   Operator `json:"ArnNotLike,omitempty"`
}

// NewCondition returns a new instance of a condition type.
func NewCondition(key, value, operator string) *Condition {
	operatorMap := Operator{key: value}

	return map[string]*Condition{
		// string condition operators
		StringEqualsOperator:              {StringEquals: operatorMap},
		StringNotEqualsOperator:           {StringNotEquals: operatorMap},
		StringEqualsIgnoreCaseOperator:    {StringEqualsIgnoreCase: operatorMap},
		StringNotEqualsIgnoreCaseOperator: {StringNotEqualsIgnoreCase: operatorMap},
		StringLikeOperator:                {StringLike: operatorMap},
		StringNotLikeOperator:             {StringNotLike: operatorMap},

		// numeric condition operators
		NumericEqualsOperator:            {NumericEquals: operatorMap},
		NumericNotEqualsOperator:         {NumericNotEquals: operatorMap},
		NumericLessThanOperator:          {NumericLessThan: operatorMap},
		NumericLessThanEqualsOperator:    {NumericLessThanEquals: operatorMap},
		NumericGreaterThanOperator:       {NumericGreaterThan: operatorMap},
		NumericGreaterThanEqualsOperator: {NumericGreaterThanEquals: operatorMap},

		// date condition operators
		DateEqualsOperator:            {DateEquals: operatorMap},
		DateNotEqualsOperator:         {DateNotEquals: operatorMap},
		DateLessThanOperator:          {DateLessThan: operatorMap},
		DateLessThanEqualsOperator:    {DateLessThanEquals: operatorMap},
		DateGreaterThanOperator:       {DateGreaterThan: operatorMap},
		DateGreaterThanEqualsOperator: {DateGreaterThanEquals: operatorMap},

		// boolean condition operators
		BoolOperator: {Bool: operatorMap},

		// binary condition operators
		BinaryEqualsOperator: {BinaryEquals: operatorMap},

		// ip condition operators
		IpAddressOperator:    {IpAddress: operatorMap},
		NotIpAddressOperator: {NotIpAddress: operatorMap},

		// arn condition operators
		ArnEqualsOperator:    {ArnEquals: operatorMap},
		ArnNotEqualsOperator: {ArnNotEquals: operatorMap},
		ArnLikeOperator:      {ArnLike: operatorMap},
		ArnNotLikeOperator:   {ArnNotLike: operatorMap},
	}[operator]
}

// String returns the string value of a tag condition.  It is used to satisfy the condition interface.  This method
// should return a unique value across various different types of conditions.  For now, all we care about are tag
// conditions, but we need to account for this if more conditions are added to the functionality in the future.
func (condition *Condition) String() string {
	jsonData, _ := json.Marshal(condition)

	return string(jsonData)
}
