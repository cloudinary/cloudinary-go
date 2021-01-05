// Package metadata defines the structured metadata.
//
// https://cloudinary.com/documentation/metadata_api
package metadata

type Field struct {
	Type         FieldType   `json:"type"`
	ExternalID   string      `json:"external_id"`
	Label        string      `json:"label"`
	Mandatory    bool        `json:"mandatory"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Validation   interface{} `json:"validation,omitempty"`
	DataSource   DataSource  `json:"datasource,omitempty"`
}

type FieldType string

const (
	StringFieldType  FieldType = "string"
	IntegerFieldType FieldType = "integer"
	DateFieldType    FieldType = "date"
	EnumFieldType    FieldType = "enum"
	SetFieldType     FieldType = "set"
)

type Validation interface{}

type ValidationType string

const (
	greaterThanValidationType  ValidationType = "greater_than"
	lessThanValidationType     ValidationType = "less_than"
	stringLengthValidationType ValidationType = "strlen"
	andValidationType          ValidationType = "and"
)

type compositeValidation struct {
	Type  ValidationType `json:"type"`
	Rules []interface{}  `json:"rules"`
}

type comparisonValidationRule struct {
	Type   ValidationType `json:"type"`
	Value  interface{}    `json:"value"`
	Equals bool           `json:"equals"`
}

type valueLengthValidationRule struct {
	Type ValidationType `json:"type"`
	Min  int            `json:"min,omitempty"`
	Max  int            `json:"max,omitempty"`
}

func AndValidation(rules []interface{}) *compositeValidation {
	return &compositeValidation{Type: andValidationType, Rules: rules}
}

func GreaterThanValidation(value interface{}, equals bool) *comparisonValidationRule {
	return &comparisonValidationRule{Type: greaterThanValidationType, Value: value, Equals: equals}
}

func LessThanValidation(value interface{}, equals bool) *comparisonValidationRule {
	return &comparisonValidationRule{Type: lessThanValidationType, Value: value, Equals: equals}
}

func StringLengthValidation(min, max int) *valueLengthValidationRule {
	return &valueLengthValidationRule{Type: stringLengthValidationType, Min: min, Max: max}
}

type DataSource struct {
	Values []DataSourceValue `json:"values"`
}

type DataSourceValue struct {
	ExternalID string `json:"external_id"`
	Value      string `json:"value"`
	State      string `json:"state"`
}
