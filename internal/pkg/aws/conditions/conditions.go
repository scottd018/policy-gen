package conditions

import (
	"encoding/json"
)

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

	switch operator {
	case StringEqualsOperator:
		return &Condition{StringEquals: operatorMap}
	case StringNotEqualsOperator:
		return &Condition{StringNotEquals: operatorMap}
	case StringEqualsIgnoreCaseOperator:
		return &Condition{StringEqualsIgnoreCase: operatorMap}
	case StringNotEqualsIgnoreCaseOperator:
		return &Condition{StringNotEqualsIgnoreCase: operatorMap}
	case StringLikeOperator:
		return &Condition{StringLike: operatorMap}
	case StringNotLikeOperator:
		return &Condition{StringNotLike: operatorMap}
	case NumericEqualsOperator:
		return &Condition{NumericEquals: operatorMap}
	case NumericNotEqualsOperator:
		return &Condition{NumericNotEquals: operatorMap}
	case NumericLessThanOperator:
		return &Condition{NumericLessThan: operatorMap}
	case NumericLessThanEqualsOperator:
		return &Condition{NumericLessThanEquals: operatorMap}
	case NumericGreaterThanOperator:
		return &Condition{NumericGreaterThan: operatorMap}
	case NumericGreaterThanEqualsOperator:
		return &Condition{NumericGreaterThanEquals: operatorMap}
	case DateEqualsOperator:
		return &Condition{DateEquals: operatorMap}
	case DateNotEqualsOperator:
		return &Condition{DateNotEquals: operatorMap}
	case DateLessThanOperator:
		return &Condition{DateLessThan: operatorMap}
	case DateLessThanEqualsOperator:
		return &Condition{DateLessThanEquals: operatorMap}
	case DateGreaterThanOperator:
		return &Condition{DateGreaterThan: operatorMap}
	case DateGreaterThanEqualsOperator:
		return &Condition{DateGreaterThanEquals: operatorMap}
	case BoolOperator:
		return &Condition{Bool: operatorMap}
	case BinaryEqualsOperator:
		return &Condition{BinaryEquals: operatorMap}
	case IpAddressOperator:
		return &Condition{IpAddress: operatorMap}
	case NotIpAddressOperator:
		return &Condition{NotIpAddress: operatorMap}
	case ArnEqualsOperator:
		return &Condition{ArnEquals: operatorMap}
	case ArnNotEqualsOperator:
		return &Condition{ArnNotEquals: operatorMap}
	case ArnLikeOperator:
		return &Condition{ArnLike: operatorMap}
	default:
		return &Condition{ArnNotLike: operatorMap}
	}
}

// String returns the string value of a tag condition.  It is used to satisfy the condition interface.  This method
// should return a unique value across various different types of conditions.  For now, all we care about are tag
// conditions, but we need to account for this if more conditions are added to the functionality in the future.
func (condition *Condition) String() string {
	jsonData, _ := json.Marshal(condition)

	return string(jsonData)
}
