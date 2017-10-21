// Helper structs and functions for representing an edge in a schema.

package codegen

import (
	"splits-go-api/privacy"
)

// EdgeStruct holds the internal representation of a schame edge.
type EdgeStruct struct {
	Name            string            // Label of the edge in neo4j (UPPER_CASE)
	CodeName        string            // Name to be used in generated code CmC
	Fields          []EdgeFieldStruct // Fields that belong to the edge
	FromNode        Schema            // Schema of the from node
	ToNode          Schema            // Schema of the to node
	ForwardsName    string
	BackwardsName   string
	Privacy         privacy.PrivacyPolicy
	ReversePrivacy  privacy.PrivacyPolicy
	WritePrivacy    privacy.PrivacyPolicy
	DeletionPrivacy privacy.PrivacyPolicy
}

// Edge constructor.
func Edge() *EdgeStruct {
	return &EdgeStruct{
		Name:            "",
		CodeName:        "",
		Fields:          []EdgeFieldStruct{},
		FromNode:        nil,
		ToNode:          nil,
		ForwardsName:    "",
		BackwardsName:   "",
		Privacy:         privacy.PrivacyPolicyStruct{},
		ReversePrivacy:  privacy.PrivacyPolicyStruct{},
		DeletionPrivacy: privacy.PrivacyPolicyStruct{},
	}
}

// SetName is the name setter for an edge.
func (es *EdgeStruct) SetName(name string) *EdgeStruct {
	es.Name = name
	return es
}

// SetCodeName is the codename setter for an edge.
func (es *EdgeStruct) SetCodeName(name string) *EdgeStruct {
	es.CodeName = name
	return es
}

// SetFields is the fields setter for an edge.
func (es *EdgeStruct) SetFields(edgeFields []EdgeFieldStruct) *EdgeStruct {
	es.Fields = edgeFields
	return es
}

// SetFromNode is the from node setter for an edge.
func (es *EdgeStruct) SetFromNode(n Schema) *EdgeStruct {
	es.FromNode = n
	return es
}

// SetToNode is the to node setter for an edge.
func (es *EdgeStruct) SetToNode(n Schema) *EdgeStruct {
	es.ToNode = n
	return es
}

// SetForwardsName is the forwards name setter for an edge.
func (es *EdgeStruct) SetForwardsName(n string) *EdgeStruct {
	es.ForwardsName = n
	return es
}

// SetBackwardsName is the backwards name setter for an edge.
func (es *EdgeStruct) SetBackwardsName(n string) *EdgeStruct {
	es.BackwardsName = n
	return es
}

// SetPrivacy is the to privacy setter for an edge.
func (es *EdgeStruct) SetPrivacy(pp privacy.PrivacyPolicy) *EdgeStruct {
	es.Privacy = pp
	return es
}

// SetReversePrivacy is the to reverse privacy setter for an edge.
func (es *EdgeStruct) SetReversePrivacy(pp privacy.PrivacyPolicy) *EdgeStruct {
	es.ReversePrivacy = pp
	return es
}

// SetDeletionPrivacy is the to deletion privacy setter for an edge.
func (es *EdgeStruct) SetDeletionPrivacy(pp privacy.PrivacyPolicy) *EdgeStruct {
	es.DeletionPrivacy = pp
	return es
}

// EdgeFieldStruct holds the internal representation of a schema edge field.
type EdgeFieldStruct struct {
	Name          string    // Name of the property in neo4j (under_scored)
	CodeName      string    // Name to be used in generated code (CamelCase)
	Type          FieldType // The type for the field (valid ones in types.go)
	DefaultValue  string    // Default value for the field (in string form)
	ExampleValue  string    // Example value for the field (in string form)
	Unique        bool      // Whether the field should be unique
	Indexed       bool      // Whether the field should be indexed
	Privacy       privacy.PrivacyPolicy
	WritePrivacy  privacy.PrivacyPolicy
	RWritePrivacy privacy.PrivacyPolicy
}

// EdgeField constructor.
func EdgeField() *EdgeFieldStruct {
	return &EdgeFieldStruct{
		Name:          "",
		CodeName:      "",
		Type:          "",
		DefaultValue:  "",
		ExampleValue:  "",
		Unique:        false,
		Indexed:       false,
		Privacy:       privacy.PrivacyPolicyStruct{},
		WritePrivacy:  privacy.PrivacyPolicyStruct{},
		RWritePrivacy: privacy.PrivacyPolicyStruct{},
	}
}

// SetName is the name setter for an edge field.
func (es *EdgeFieldStruct) SetName(name string) *EdgeFieldStruct {
	es.Name = name
	return es
}

// SetCodeName is the codename setter for an edge field.
func (es *EdgeFieldStruct) SetCodeName(name string) *EdgeFieldStruct {
	es.CodeName = name
	return es
}

// SetDefaultValue is the default value setter for an edge field.
func (es *EdgeFieldStruct) SetDefaultValue(v string) *EdgeFieldStruct {
	es.DefaultValue = v
	return es
}

// SetExampleValue is the test value setter for a node field.
func (es *EdgeFieldStruct) SetExampleValue(v string) *EdgeFieldStruct {
	es.ExampleValue = v
	return es
}

// SetType is the type setter for an edge field.
func (es *EdgeFieldStruct) SetType(t FieldType) *EdgeFieldStruct {
	es.Type = t
	return es
}

// SetUnique is the unique setter for an edge field.
func (es *EdgeFieldStruct) SetUnique(unique bool) *EdgeFieldStruct {
	es.Unique = unique
	return es
}

// SetIndexed is the indexed setter for an edge field.
func (es *EdgeFieldStruct) SetIndexed(indexed bool) *EdgeFieldStruct {
	es.Indexed = indexed
	return es
}

// SetPrivacy is the privacy setter for an edge field.
func (es *EdgeFieldStruct) SetPrivacy(
	pp privacy.PrivacyPolicy,
) *EdgeFieldStruct {
	es.Privacy = pp
	return es
}

// SetWritePrivacy is the privacy setter for writing an edge field.
func (es *EdgeFieldStruct) SetWritePrivacy(
	pp privacy.PrivacyPolicy,
) *EdgeFieldStruct {
	es.WritePrivacy = pp
	return es
}
