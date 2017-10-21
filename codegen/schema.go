// Representation of a schema for code generation.

package codegen

import "splits-go-api/privacy"

// Schema interface for code generation.
type Schema interface {
	GetName() string
	GetFields() []FieldStruct
	GetEdges() []EdgeStruct
	GetEdgePointers() map[string]EdgeStruct
	AddEdgePointer(e EdgeStruct)
	GetDeletionPrivacy() privacy.Policy
}
