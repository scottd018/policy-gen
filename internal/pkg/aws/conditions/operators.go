package conditions

// see https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_String
// for a complete list of operators.
//
//nolint:revive,stylecheck
const (
	// string condition operators.
	StringEqualsOperator              = "StringEquals"
	StringNotEqualsOperator           = "StringNotEquals"
	StringEqualsIgnoreCaseOperator    = "StringEqualsIgnoreCase"
	StringNotEqualsIgnoreCaseOperator = "StringNotEqualsIgnoreCase"
	StringLikeOperator                = "StringLike"
	StringNotLikeOperator             = "StringNotLike"

	// numeric condition operators.
	NumericEqualsOperator            = "NumericEquals"
	NumericNotEqualsOperator         = "NumericNotEquals"
	NumericLessThanOperator          = "NumericLessThan"
	NumericLessThanEqualsOperator    = "NumericLessThanEquals"
	NumericGreaterThanOperator       = "NumericGreaterThan"
	NumericGreaterThanEqualsOperator = "NumericGreaterThanEquals"

	// date condition operators.
	DateEqualsOperator            = "DateEquals"
	DateNotEqualsOperator         = "DateNotEquals"
	DateLessThanOperator          = "DateLessThan"
	DateLessThanEqualsOperator    = "DateLessThanEquals"
	DateGreaterThanOperator       = "DateGreaterThan"
	DateGreaterThanEqualsOperator = "DateGreaterThanEquals"

	// boolean condition operators.
	BoolOperator = "Bool"

	// binary condition operators.
	BinaryEqualsOperator = "BinaryEquals"

	// ip condition operators.
	IpAddressOperator    = "IpAddress"
	NotIpAddressOperator = "NotIpAddress"

	// arn condition operators.
	ArnEqualsOperator    = "ArnEquals"
	ArnNotEqualsOperator = "ArnNotEquals"
	ArnLikeOperator      = "ArnLike"
	ArnNotLikeOperator   = "ArnNotLike"
)

// Operator is an operator which represents an operator condition.
type Operator map[string]string

// ToOperatorString returns the OperatorString value of a string-typed operator.
func ToOperatorString(operator string) string {
	return map[string]string{
		// string condition operators
		StringEqualsOperator:              StringEqualsOperator,
		StringNotEqualsOperator:           StringNotEqualsOperator,
		StringEqualsIgnoreCaseOperator:    StringEqualsIgnoreCaseOperator,
		StringNotEqualsIgnoreCaseOperator: StringNotEqualsIgnoreCaseOperator,
		StringLikeOperator:                StringLikeOperator,
		StringNotLikeOperator:             StringNotLikeOperator,

		// numeric condition operators
		NumericEqualsOperator:            NumericEqualsOperator,
		NumericNotEqualsOperator:         NumericNotEqualsOperator,
		NumericLessThanOperator:          NumericLessThanOperator,
		NumericLessThanEqualsOperator:    NumericLessThanEqualsOperator,
		NumericGreaterThanOperator:       NumericGreaterThanOperator,
		NumericGreaterThanEqualsOperator: NumericGreaterThanEqualsOperator,

		// date condition operators
		DateEqualsOperator:            DateEqualsOperator,
		DateNotEqualsOperator:         DateNotEqualsOperator,
		DateLessThanOperator:          DateLessThanOperator,
		DateLessThanEqualsOperator:    DateLessThanEqualsOperator,
		DateGreaterThanOperator:       DateGreaterThanOperator,
		DateGreaterThanEqualsOperator: DateGreaterThanEqualsOperator,

		// boolean condition operators
		BoolOperator: BoolOperator,

		// binary condition operators
		BinaryEqualsOperator: BinaryEqualsOperator,

		// ip condition operators
		IpAddressOperator:    IpAddressOperator,
		NotIpAddressOperator: NotIpAddressOperator,

		// arn condition operators
		ArnEqualsOperator:    ArnEqualsOperator,
		ArnNotEqualsOperator: ArnNotEqualsOperator,
		ArnLikeOperator:      ArnLikeOperator,
		ArnNotLikeOperator:   ArnNotLikeOperator,
	}[operator]
}
