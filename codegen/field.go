// Helper structs and functions for representing a node in a schema.

package codegen

import (
	"splits-go-api/privacy"
)

// FieldStruct holds the internal representation of a schame node.
type FieldStruct struct {
	Name         string    // Name of the field (under_scored)
	CodeName     string    // Name to be used in generated code (CamelCase)
	Type         FieldType // The type for the field (valid ones in types.go)
	DefaultValue string    // Default value for the field (in string form)
	ExampleValue string    // Example value for the field (in string form)
	Unique       bool      // Whether the field be unique
	Indexed      bool      // Whether the field should have an index on it
	Privacy      privacy.Policy
	WritePrivacy privacy.Policy
	GQLField     *GraphQLField
	CanOrderBy   bool
}

// Field constructor.
func Field() *FieldStruct {
	return &FieldStruct{
		Name:         "",
		CodeName:     "",
		Type:         "",
		DefaultValue: "",
		ExampleValue: "",
		Unique:       false,
		Indexed:      false,
		Privacy:      privacy.PolicyStruct{},
		WritePrivacy: privacy.PolicyStruct{},
		GQLField:     nil,
		CanOrderBy:   false,
	}
}

// SetName is the name setter for a node field.
func (fs *FieldStruct) SetName(name string) *FieldStruct {
	fs.Name = name
	return fs
}

// SetCodeName is the codename setter for a node field.
func (fs *FieldStruct) SetCodeName(name string) *FieldStruct {
	fs.CodeName = name
	return fs
}

// SetType is the type setter for a node field.
func (fs *FieldStruct) SetType(t FieldType) *FieldStruct {
	fs.Type = t
	return fs
}

// SetDefaultValue is the default value setter for a node field.
func (fs *FieldStruct) SetDefaultValue(v string) *FieldStruct {
	fs.DefaultValue = v
	return fs
}

// SetExampleValue is the example value setter for a node field.
func (fs *FieldStruct) SetExampleValue(v string) *FieldStruct {
	fs.ExampleValue = v
	return fs
}

// SetUnique is the unique setter for a node field.
func (fs *FieldStruct) SetUnique(unique bool) *FieldStruct {
	fs.Unique = unique
	return fs
}

// SetIndexed is the indexed setter for a node field.
func (fs *FieldStruct) SetIndexed(indexed bool) *FieldStruct {
	fs.Indexed = indexed
	return fs
}

// SetPrivacy is the privacy setter for a node field.
func (fs *FieldStruct) SetPrivacy(pp privacy.Policy) *FieldStruct {
	fs.Privacy = pp
	return fs
}

// SetWritePrivacy is the privacy setter for writing a node field.
func (fs *FieldStruct) SetWritePrivacy(pp privacy.Policy) *FieldStruct {
	fs.WritePrivacy = pp
	return fs
}

// SetGQLField is the graphql field setter for hte node field.
func (fs *FieldStruct) SetGQLField(gqlField *GraphQLField) *FieldStruct {
	fs.GQLField = gqlField
	return fs
}

// SetCanOrderBy is the can order by setter for ordering edges.
func (fs *FieldStruct) SetCanOrderBy(canOrderBy bool) *FieldStruct {
	fs.CanOrderBy = canOrderBy
	return fs
}
