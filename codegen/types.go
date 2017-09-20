// Valid types for fields.

package codegen

// FieldType is a string wrapper for types of fields.
type FieldType string

// Types of valid field values.
const (
	StringType = FieldType("string")
	FloatType  = FieldType("float64")
	IntType    = FieldType("int64")
)
