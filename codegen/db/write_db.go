// The writer assumes a valid schema. The generation involves generating a few
// pieces: node queries, mutators, and deleters, and edge queries, mutators, and
// deleters.

package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/format"
	c "splits-go-api/db/models/constraints"
	i "splits-go-api/db/models/indices"
	cg "splits-go-schema-codegen/codegen"
	"strings"
)

// WriteSchemaNode generates the string that represents a schema node.
func WriteSchemaNode(s cg.Schema, packageName string) string {

	// Use templates to generate the node
	sections := make([]string, 0, 11)
	sections = append(sections, GetNodeFileHeaderCommentStr(s))
	sections = append(sections, GetNodePackageStr(s, packageName))
	sections = append(sections, GetNodeImportStr(s))
	sections = append(sections, GetNodeStr(s))
	sections = append(sections, GetNodeQueryStructStr(s))
	sections = append(sections, GetNodeQueryConstructorStr(s))
	sections = append(sections, GetNodeQueryWhereStr(s))
	sections = append(sections, GetNodeQueryReturnStr(s))
	sections = append(sections, GetNodeQueryOrderStr(s))
	sections = append(sections, GetNodeQueryEdgesStr(s))
	sections = append(sections, GetNodeMutatorStr(s))
	sections = append(sections, GetNodeDeleterStr(s))
	result := strings.Join(sections, "\n")
	res, err := format.Source([]byte(result))
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	// Generate the MD5 signature
	sum := md5.Sum([]byte(res))
	signature := hex.EncodeToString([]byte(sum[:]))

	// Add the signature to the top of the file
	return "// @SignedSource (" + signature + ")\n" + string(res)
}

// WriteSchemaEdge generates the string that represents a schema edge.
func WriteSchemaEdge(s cg.Schema, e cg.EdgeStruct, packageName string) string {

	// Use templates to generate the edge
	sections := make([]string, 0, 11)
	sections = append(sections, GetEdgeFileHeaderCommentStr(e))
	sections = append(sections, GetEdgePackageStr(s, packageName))
	sections = append(sections, GetEdgeImportStr(e))
	sections = append(sections, GetEdgeStr(e))
	sections = append(sections, GetEdgeQueryStructStr(e))
	sections = append(sections, GetEdgeQueryConstructorStr(e))
	sections = append(sections, GetEdgeQueryWhereStr(e))
	sections = append(sections, GetEdgeQueryReturnStr(e))
	sections = append(sections, GetEdgeQueryOrderStr(e))
	sections = append(sections, GetEdgeQueryNodesStr(e))
	sections = append(sections, GetEdgeMutatorStr(e))
	sections = append(sections, GetEdgeDeleterStr(e))
	result := strings.Join(sections, "\n")
	res, err := format.Source([]byte(result))
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	// Generate the MD5 signature
	sum := md5.Sum([]byte(res))
	signature := hex.EncodeToString([]byte(sum[:]))

	// Add the signature to the top of the file
	return "// @SignedSource (" + signature + ")\n" + string(res)
}

// WriteConstraints generates the string that represents the constraints.
func WriteConstraints(schemas []cg.Schema) string {
	cd := c.ConstraintData{
		Nodes: []c.ConstraintNode{},
		Edges: []c.ConstraintEdge{},
	}
	for _, s := range schemas {
		cn := new(c.ConstraintNode)
		cn.Type = s.GetName()
		cn.Properties = []string{}
		for _, f := range s.GetFields() {
			if f.Unique {
				cn.Properties = append(cn.Properties, f.Name)
			}
		}
		for _, e := range s.GetEdges() {
			ce := new(c.ConstraintEdge)
			ce.Type = e.Name
			ce.Properties = []string{}
			for _, f := range e.Fields {
				if f.Unique {
					ce.Properties = append(ce.Properties, f.Name)
				}
			}
			cd.Edges = append(cd.Edges, *ce)
		}
		cd.Nodes = append(cd.Nodes, *cn)
	}

	res, err := json.MarshalIndent(cd, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(res)
}

// WriteIndices generates the string that represents the indices.
func WriteIndices(schemas []cg.Schema) string {
	id := i.IndexData{
		Nodes: []i.IndexNode{},
	}
	for _, s := range schemas {
		in := new(i.IndexNode)
		in.Type = s.GetName()
		in.Properties = []string{}
		for _, f := range s.GetFields() {
			if f.Indexed {
				in.Properties = append(in.Properties, f.Name)
			}
		}
		id.Nodes = append(id.Nodes, *in)
	}

	res, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(res)
}

// WriteConstants helps write some constants.
func WriteConstants(schemas []cg.Schema) string {
	constants := map[string]string{}
	for _, s := range schemas {
		constants[s.GetName()+"Label"] = s.GetName()
		for _, e := range s.GetEdges() {
			constants[e.CodeName+"Label"] = e.Name
		}
	}
	data := struct {
		Constants map[string]string
	}{
		Constants: constants,
	}
	template :=
		"package models\n\n" +
			"var constants = struct {\n" +
			"{{ range $var, $value := .Constants }}" +
			"\t{{$var}} string\n" +
			"{{ end }}" +
			"} {\n" +
			"{{ range $var, $value := .Constants }}" +
			"\t{{$var}}: \"{{$value}}\",\n" +
			"{{ end }}" +
			"}"

	result := cg.ExecTemplate(template, "constants", data, nil)
	res, err := format.Source([]byte(result))
	if err != nil {
		fmt.Println(result)
		panic(err)
	}
	return string(res)
}

// =============================================================================
// Node generation
// =============================================================================

// GetNodeFileHeaderCommentStr generates an autogenerated tag.
func GetNodeFileHeaderCommentStr(s cg.Schema) string {
	data := struct {
		Name string
	}{
		Name: s.GetName(),
	}
	template := "// Autogenerated {{.Name}} - regenerate with splits-go-schema-" +
		"codegen\n// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "node_file_header_comment", data, nil)
}

// GetNodePackageStr generates the package tag.
func GetNodePackageStr(s cg.Schema, packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "node_package", data, nil)
}

// GetNodeImportStr generates the import statements.
func GetNodeImportStr(s cg.Schema) string {
	data := struct {
		Imports []string
	}{
		Imports: []string{
			"\"splits-go-api/db/models/base\"",
			"p \"splits-go-api/db/models/predicates\"",
		},
	}
	template :=
		"import (\n" +
			"{{range .Imports}} \t{{.}}\n {{end}}" +
			")\n"
	return cg.ExecTemplate(template, "node_import", data, nil)
}

// GetNodeStr generates the base node definition.
func GetNodeStr(s cg.Schema) string {
	edgePointers := make([]cg.EdgeStruct, 0, len(s.GetEdgePointers()))
	for _, v := range s.GetEdgePointers() {
		edgePointers = append(edgePointers, v)
	}
	data := struct {
		Name         string
		Fields       []cg.FieldStruct
		Edges        []cg.EdgeStruct
		EdgePointers []cg.EdgeStruct
	}{
		Name:         s.GetName(),
		Fields:       s.GetFields(),
		Edges:        s.GetEdges(),
		EdgePointers: edgePointers,
	}
	template := "// {{.Name}}Node is the base {{.Name}} definition.\n" +
		"type {{.Name}}Node struct {\n" +
		"\t// Node fields\n" +
		"{{range .Fields}} \t{{.CodeName}} {{.Type}}\n{{end}}\n" +
		"\t// Edges\n" +
		"{{range .Edges}} \t{{.CodeName}} *{{.CodeName}}Edge\n{{end}}" +
		"{{range .EdgePointers}} \t{{.CodeName}} *{{.CodeName}}Edge\n{{end}}\n" +
		"}\n"
	return cg.ExecTemplate(template, "node", data, nil)
}

// GetNodeQueryStructStr generates the base node query struct.
func GetNodeQueryStructStr(s cg.Schema) string {
	data := struct {
		Name string
	}{
		Name: s.GetName(),
	}
	template := "// {{.Name}}Q is the base {{.Name}} query struct.\n" +
		"type {{.Name}}Q struct {\n" +
		"\tbase.Query\n" +
		"}\n"
	return cg.ExecTemplate(template, "node_query", data, nil)
}

// GetNodeQueryConstructorStr generates the base node query constructor.
func GetNodeQueryConstructorStr(s cg.Schema) string {
	data := struct {
		Name    string
		VarName string
	}{
		Name:    s.GetName(),
		VarName: strings.ToLower(string(s.GetName()[0])),
	}
	template := "// {{.Name}}Query is the {{.Name}} query constructor.\n" +
		"func {{.Name}}Query() *{{.Name}}Q {\n" +
		"\t{{.VarName}} := new({{.Name}}Q)\n" +
		"\t{{.VarName}}.Fields = []p.WhereClauseStruct{}\n" +
		"\t{{.VarName}}.Return = []p.ReturnClauseStruct{}\n" +
		"\t{{.VarName}}.IsNode = true\n" +
		"\t{{.VarName}}.Prefix = 'a'\n " +
		"\t{{.VarName}}.Label = constants.{{.Name}}Label\n" +
		"\treturn {{.VarName}}\n" +
		"}"
	return cg.ExecTemplate(template, "node_query_constructor", data, nil)
}

// GetNodeQueryWhereStr generates all the WhereClause functions for a node.
func GetNodeQueryWhereStr(s cg.Schema) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.FieldStruct
	}{
		Name:    s.GetName(),
		VarName: strings.ToLower(string(s.GetName()[0])) + "q",
		Fields:  s.GetFields(),
	}
	template := "{{range .Fields}}" +
		"// Where{{.CodeName}} is the query where clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Where{{.CodeName}}(pred p.Predicate) " +
		"*{{$.Name}}Q {\n" +
		"{{$.VarName}}.Fields = " +
		"append({{$.VarName}}.Fields, p.WhereClause(\"{{.Name}}\", pred))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "node_query_where", data, nil)
}

// GetNodeQueryReturnStr generates all the Return clause functions for a node.
func GetNodeQueryReturnStr(s cg.Schema) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.FieldStruct
	}{
		Name:    s.GetName(),
		VarName: strings.ToLower(string(s.GetName()[0])) + "q",
		Fields:  s.GetFields(),
	}

	template := "{{range .Fields}}" +
		"// Return{{.CodeName}} is the return clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Return{{.CodeName}}() *{{$.Name}}Q {\n" +
		"{{$.VarName}}.Return = append({{$.VarName}}.Return, p.ReturnClause" +
		"(\"{{.Name}}\"))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "node_query_return", data, nil)
}

// GetNodeQueryOrderStr generates all the Order clause functions for a node.
func GetNodeQueryOrderStr(s cg.Schema) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.FieldStruct
	}{
		Name:    s.GetName(),
		VarName: strings.ToLower(string(s.GetName()[0])) + "q",
		Fields:  s.GetFields(),
	}

	template := "{{range .Fields}}" +
		"// OrderBy{{.CodeName}} is the order clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) OrderBy{{.CodeName}}(desc bool) " +
		"*{{$.Name}}Q {\n" +
		"{{$.VarName}}.Order = append({{$.VarName}}.Order, p.OrderClause" +
		"(\"{{.Name}}\", desc))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "node_query_order", data, nil)
}

// GetNodeQueryEdgesStr generates the Query functions for traversing the graph.
func GetNodeQueryEdgesStr(s cg.Schema) string {
	edgePointers := make([]cg.EdgeStruct, 0, len(s.GetEdgePointers()))
	for _, v := range s.GetEdgePointers() {
		edgePointers = append(edgePointers, v)
	}
	data := struct {
		Name         string
		VarName      string
		Edges        []cg.EdgeStruct
		EdgePointers []cg.EdgeStruct
	}{
		Name:         s.GetName(),
		VarName:      strings.ToLower(string(s.GetName()[0])) + "q",
		Edges:        s.GetEdges(),
		EdgePointers: edgePointers,
	}
	template := "{{range .Edges}}" +
		"// Query{{.CodeName}} traverses the graph to the {{.CodeName}} edge.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Query{{.CodeName}}() *{{.CodeName}}Q " +
		"{\n" +
		"\tquery := {{.CodeName}}Query()\n" +
		"\tquery.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tquery.Prev = &{{$.VarName}}.Query\n" +
		"\t{{$.VarName}}.Next = &query.Query\n" +
		"\treturn query\n" +
		"}\n\n" +
		"{{end}}" +

		"{{range .EdgePointers}}" +
		"// Query{{.CodeName}} traverses the graph to the {{.CodeName}} edge.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Query{{.CodeName}}() *{{.CodeName}}Q " +
		"{\n" +
		"\tquery := {{.CodeName}}Query()\n" +
		"\tquery.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tquery.Prev = &{{$.VarName}}.Query\n" +
		"\t{{$.VarName}}.Next = &query.Query\n" +
		"\treturn query\n" +
		"}\n\n" +
		"{{end}}"

	return cg.ExecTemplate(template, "node_query_edges", data, nil)
}

// GetNodeMutatorStr generates the mutator helper functions.
func GetNodeMutatorStr(s cg.Schema) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.FieldStruct
	}{
		Name:    s.GetName(),
		VarName: strings.ToLower(string(s.GetName()[0])) + "m",
		Fields:  s.GetFields(),
	}
	// Base mutator
	template := "// {{.Name}}M is the base {{.Name}} mutator struct.\n" +
		"type {{.Name}}M struct {\n" +
		"\tbase.NodeMutator\n" +
		"}\n\n" +

		// Mutator constructor
		"// {{.Name}}Mutator is the {{.Name}} mutator constructor.\n" +
		"func {{.Name}}Mutator(id string) *{{.Name}}M " +
		"{\n" +
		"\t{{.VarName}} := new({{.Name}}M)\n" +
		"\t{{.VarName}}.ID = id\n" +
		"\t{{.VarName}}.Fields = map[string]interface{}{}\n" +
		"\t{{.VarName}}.DefaultFields = map[string]interface{}{}\n" +
		"\t{{.VarName}}.Label = constants.{{.Name}}Label\n" +

		// Default fields
		"{{range .Fields}}" +
		"\t{{$.VarName}}.DefaultFields[\"{{.Name}}\"] = {{.DefaultValue}}\n" +
		"{{end}}" +

		"\treturn {{.VarName}}\n" +
		"}\n\n" +

		// Setters for the mutator
		"{{range .Fields}}" +
		"// Set{{.CodeName}} is the mutator setter for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}M) Set{{.CodeName}}(v {{.Type}}) " +
		"*{{$.Name}}M {\n" +
		"{{$.VarName}}.Fields[\"{{.Name}}\"] = v\n" +
		"{{$.VarName}}.DefaultFields[\"{{.Name}}\"] = v\n" +
		"\treturn {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"

	return cg.ExecTemplate(template, "node_mutator", data, nil)
}

// GetNodeDeleterStr generates the deleter helper functions.
func GetNodeDeleterStr(s cg.Schema) string {
	data := struct {
		Name         string
		VarName      string
		Fields       []cg.FieldStruct
		Edges        []cg.EdgeStruct
		EdgePointers map[string]cg.EdgeStruct
	}{
		Name:         s.GetName(),
		VarName:      strings.ToLower(string(s.GetName()[0])) + "d",
		Fields:       s.GetFields(),
		Edges:        s.GetEdges(),
		EdgePointers: s.GetEdgePointers(),
	}
	// Base deleter
	template := "// {{.Name}}D is the base {{.Name}} deleter struct.\n" +
		"type {{.Name}}D struct {\n" +
		"\nbase.Deleter\n" +
		"}\n\n" +

		// Deleter constructor
		"// {{.Name}}Deleter is the {{.Name}} deleter constructor.\n" +
		"func {{.Name}}Deleter() *{{.Name}}D {\n" +
		"\t{{.VarName}} := new({{.Name}}D)\n" +
		"\t{{.VarName}}.Prefix = 'a'\n" +
		"\t{{.VarName}}.IsNode = true\n" +
		"\t{{.VarName}}.Fields = []p.WhereClauseStruct{}\n" +
		"\t{{.VarName}}.Label = constants.{{.Name}}Label\n" +
		"\treturn {{.VarName}}\n" +
		"}\n\n" +

		// Where clauses for the deleter
		"{{range .Fields}}" +
		"// Where{{.CodeName}} is the deleter where clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Where{{.CodeName}}(pred p.Predicate) " +
		"*{{$.Name}}D {\n" +
		"{{$.VarName}}.Fields = " +
		"append({{$.VarName}}.Fields, p.WhereClause(\"{{.Name}}\", pred))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}" +

		// Delete clause for the deleter
		"// Delete the actual node\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete() *{{$.Name}}D {\n" +
		"{{$.VarName}}.WillDelete = true\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +

		// Traverse the graph
		"{{range .Edges}}" +
		"// Delete{{.CodeName}} traverses the deleter to the {{.CodeName}} " +
		"edge.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete{{.CodeName}}() *{{.CodeName}}D " +
		"{\n" +
		"\tdeleter := {{.CodeName}}Deleter()\n" +
		"\tdeleter.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tdeleter.Prev = &{{$.VarName}}.Deleter\n" +
		"\t{{$.VarName}}.Next = &deleter.Deleter\n" +
		"\treturn deleter\n" +
		"}\n\n" +
		"{{end}}" +

		"{{range .EdgePointers}}" +
		"// Delete{{.CodeName}} traverses the deleter to the {{.CodeName}} " +
		"edge.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete{{.CodeName}}() *{{.CodeName}}D " +
		"{\n" +
		"\tdeleter:= {{.CodeName}}Deleter()\n" +
		"\tdeleter.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tdeleter.Prev = &{{$.VarName}}.Deleter\n" +
		"\t{{$.VarName}}.Next = &deleter.Deleter\n" +
		"\treturn deleter\n" +
		"}\n\n" +
		"{{end}}"

	return cg.ExecTemplate(template, "node_deleter", data, nil)
}

// =============================================================================
// Edge generation
// =============================================================================

// GetEdgeFileHeaderCommentStr generates an autogenerated tag.
func GetEdgeFileHeaderCommentStr(e cg.EdgeStruct) string {
	data := struct {
		Name string
	}{
		Name: e.CodeName,
	}
	template := "// Autogenerated {{.Name}} - regenerate with " +
		"splits-go-schema-codegen\n" +
		"// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "edge_file_header_comment", data, nil)
}

// GetEdgePackageStr generates the package tag.
func GetEdgePackageStr(s cg.Schema, packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "edge_package", data, nil)
}

// GetEdgeImportStr generates the import statements.
func GetEdgeImportStr(e cg.EdgeStruct) string {
	data := struct {
		Imports []string
	}{
		Imports: []string{
			"\"splits-go-api/db/models/base\"",
			"p \"splits-go-api/db/models/predicates\"",
		},
	}
	template :=
		"import (\n" +
			"{{range .Imports}} \t{{.}}\n {{end}}" +
			")\n"
	return cg.ExecTemplate(template, "edge_import", data, nil)
}

// GetEdgeStr generates the base edge definition.
func GetEdgeStr(e cg.EdgeStruct) string {
	data := struct {
		Name     string
		CodeName string
		Fields   []cg.EdgeFieldStruct
		FromNode cg.Schema
		ToNode   cg.Schema
	}{
		Name:     e.Name,
		CodeName: e.CodeName,
		Fields:   e.Fields,
		FromNode: e.FromNode,
		ToNode:   e.ToNode,
	}
	template := "// {{.CodeName}}Edge is the base {{.CodeName}} definition.\n" +
		"type {{.CodeName}}Edge struct {\n" +
		"\t// Edge fields\n" +
		"{{range .Fields}} \t{{.CodeName}} {{.Type}}\n{{end}}\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge", data, nil)
}

// GetEdgeQueryStructStr generates the base edge query struct.
func GetEdgeQueryStructStr(e cg.EdgeStruct) string {
	data := struct {
		Name string
	}{
		Name: e.CodeName,
	}
	template := "// {{.Name}}Q is the base {{.Name}} query struct.\n" +
		"type {{.Name}}Q struct {\n" +
		"\tbase.Query\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_query", data, nil)
}

// GetEdgeQueryConstructorStr generates the base edge query constructor.
func GetEdgeQueryConstructorStr(e cg.EdgeStruct) string {
	data := struct {
		Name     string
		CodeName string
		VarName  string
	}{
		Name:     e.Name,
		CodeName: e.CodeName,
		VarName:  strings.ToLower(string(e.Name[0])),
	}
	t := "// {{.CodeName}}Query is the {{.CodeName}} query constructor.\n" +
		"func {{.CodeName}}Query() *{{.CodeName}}Q {\n" +
		"\t{{.VarName}} := new({{.CodeName}}Q)\n" +
		"\t{{.VarName}}.Fields = []p.WhereClauseStruct{}\n" +
		"\t{{.VarName}}.Return = []p.ReturnClauseStruct{}\n" +
		"\t{{.VarName}}.IsNode = false\n" +
		"\t{{.VarName}}.Prefix = 'a'\n" +
		"\t{{.VarName}}.Label = constants.{{.CodeName}}Label\n" +
		"\treturn {{.VarName}}\n" +
		"}\n"
	return cg.ExecTemplate(t, "node_query_constructor", data, nil)
}

// GetEdgeQueryWhereStr generates all the WhereClause functions for an edge.
func GetEdgeQueryWhereStr(e cg.EdgeStruct) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.EdgeFieldStruct
	}{
		Name:    e.CodeName,
		VarName: strings.ToLower(string(e.CodeName[0])) + "q",
		Fields:  e.Fields,
	}
	t := "{{range .Fields}}" +
		"// Where{{.CodeName}} is the where clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Where{{.CodeName}}(pred p.Predicate) " +
		"*{{$.Name}}Q {\n" +
		"{{$.VarName}}.Fields = " +
		"append({{$.VarName}}.Fields, p.WhereClause(\"{{.Name}}\", pred))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(t, "edge_query_where", data, nil)
}

// GetEdgeQueryReturnStr generates all the Return clause functions for an edge.
func GetEdgeQueryReturnStr(e cg.EdgeStruct) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.EdgeFieldStruct
	}{
		Name:    e.CodeName,
		VarName: strings.ToLower(string(e.Name[0])) + "q",
		Fields:  e.Fields,
	}

	template := "{{range .Fields}}" +
		"// Return{{.CodeName}} is the return clause for {{.CodeName}}\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) Return{{.CodeName}}() *{{$.Name}}Q {\n" +
		"{{$.VarName}}.Return = append({{$.VarName}}.Return, p.ReturnClause" +
		"(\"{{.Name}}\"))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "edge_query_return", data, nil)
}

// GetEdgeQueryOrderStr generates all the Order clause functions for an edge.
func GetEdgeQueryOrderStr(e cg.EdgeStruct) string {
	data := struct {
		Name    string
		VarName string
		Fields  []cg.EdgeFieldStruct
	}{
		Name:    e.CodeName,
		VarName: strings.ToLower(string(e.Name[0])) + "q",
		Fields:  e.Fields,
	}

	template := "{{range .Fields}}" +
		"// OrderBy{{.CodeName}} is the return clause for {{.CodeName}}\n" +
		"func ({{$.VarName}} *{{$.Name}}Q) OrderBy{{.CodeName}}(desc bool) *{{$.Name}}Q {\n" +
		"{{$.VarName}}.Order = append({{$.VarName}}.Order, p.OrderClause" +
		"(\"{{.Name}}\", desc))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "edge_query_order", data, nil)
}

// GetEdgeQueryNodesStr generates the Query functions for traversing the graph.
func GetEdgeQueryNodesStr(e cg.EdgeStruct) string {
	data := struct {
		Name           string
		CodeName       string
		VarName        string
		FromNode       string
		ToNode         string
		DifferentNodes bool
	}{
		Name:           e.Name,
		CodeName:       e.CodeName,
		VarName:        strings.ToLower(string(e.Name[0])) + "q",
		FromNode:       e.FromNode.GetName(),
		ToNode:         e.ToNode.GetName(),
		DifferentNodes: e.FromNode.GetName() != e.ToNode.GetName(),
	}
	t := "// Query{{.FromNode}} traverses the graph to the {{.FromNode}} node." +
		"\nfunc ({{$.VarName}} *{{$.CodeName}}Q) Query{{.FromNode}}() " +
		"*{{.FromNode}}Q {\n\tquery := {{.FromNode}}Query()\n" +
		"\tquery.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tquery.Prev = &{{$.VarName}}.Query\n" +
		"\t{{$.VarName}}.Next = &query.Query\n" +
		"\treturn query\n" +
		"}\n\n" +
		"{{if .DifferentNodes}}" +
		"// Query{{.ToNode}} traverses the graph to the {{.ToNode}} node.\n" +
		"func ({{$.VarName}} *{{$.CodeName}}Q) Query{{.ToNode}}" +
		"() *{{.ToNode}}Q {" +
		"\n\tquery := {{.ToNode}}Query()\n" +
		"\tquery.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tquery.Prev = &{{$.VarName}}.Query\n" +
		"\t{{$.VarName}}.Next = &query.Query\n" +
		"\treturn query\n" +
		"}\n" +
		"{{end}}"
	return cg.ExecTemplate(t, "edge_query_nodes", data, nil)
}

// GetEdgeMutatorStr generates the mutator helper functions.
func GetEdgeMutatorStr(e cg.EdgeStruct) string {
	data := struct {
		Name     string
		VarName  string
		FromNode string
		ToNode   string
		Fields   []cg.EdgeFieldStruct
	}{
		Name:     e.CodeName,
		VarName:  strings.ToLower(string(e.Name[0])) + "m",
		FromNode: e.FromNode.GetName(),
		ToNode:   e.ToNode.GetName(),
		Fields:   e.Fields,
	}
	// Base mutator
	template := "// {{.Name}}M is the base {{.Name}} mutator struct.\n" +
		"type {{.Name}}M struct {\n" +
		"\tbase.EdgeMutator\n" +
		"}\n\n" +

		// Mutator constructor
		"// {{.Name}}Mutator is the {{.Name}} mutator constructor.\n" +
		"func {{.Name}}Mutator(id string, " +
		"fromID string, toID string) *{{.Name}}M {\n" +
		"\t{{.VarName}} := new({{.Name}}M)\n" +
		"\t{{.VarName}}.ID = id\n" +
		"\t{{.VarName}}.Fields = map[string]interface{}{}\n" +
		"\t{{.VarName}}.DefaultFields = map[string]interface{}{}\n" +
		"\t{{.VarName}}.IsNode = true\n" +
		"\t{{.VarName}}.FromNode = constants.{{.FromNode}}Label\n" +
		"\t{{.VarName}}.ToNode = constants.{{.ToNode}}Label\n" +
		"\t{{.VarName}}.FromID = fromID\n" +
		"\t{{.VarName}}.ToID = toID\n" +
		"\t{{.VarName}}.Label = constants.{{.Name}}Label\n" +

		// Default fields
		"{{range .Fields}}" +
		"\t{{$.VarName}}.DefaultFields[\"{{.Name}}\"] = {{.DefaultValue}}\n" +
		"{{end}}" +

		"\treturn {{.VarName}}\n" +
		"}\n\n" +

		// Setters for the mutator
		"{{range .Fields}}" +
		"// Set{{.CodeName}} is the mutator setter for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}M) Set{{.CodeName}}(v {{.Type}}) " +
		"*{{$.Name}}M {\n" +
		"{{$.VarName}}.Fields[\"{{.Name}}\"] = v\n" +
		"\treturn {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}"

	return cg.ExecTemplate(template, "edge_mutator", data, nil)
}

// GetEdgeDeleterStr generates the deleter helper functions.
func GetEdgeDeleterStr(e cg.EdgeStruct) string {
	data := struct {
		Name     string
		VarName  string
		FromNode string
		ToNode   string
		Fields   []cg.EdgeFieldStruct
	}{
		Name:     e.CodeName,
		VarName:  strings.ToLower(string(e.Name[0])) + "m",
		FromNode: e.FromNode.GetName(),
		ToNode:   e.ToNode.GetName(),
		Fields:   e.Fields,
	}
	// Base deleter
	template := "// {{.Name}}D is the base {{.Name}} deleter struct.\n" +
		"type {{.Name}}D struct {\n" +
		"\nbase.Deleter\n" +
		"}\n\n" +

		// Deleter constructor
		"// {{.Name}}Deleter is the {{.Name}} deleter constructor.\n" +
		"func {{.Name}}Deleter() *{{.Name}}D {\n" +
		"\t{{.VarName}} := new({{.Name}}D)\n" +
		"\t{{.VarName}}.Prefix = 'a'\n" +
		"\t{{.VarName}}.IsNode = false\n" +
		"\t{{.VarName}}.Fields = []p.WhereClauseStruct{}\n" +
		"\t{{.VarName}}.Label = constants.{{.Name}}Label\n" +
		"\treturn {{.VarName}}\n" +
		"}\n\n" +

		// Where clauses for the deleter
		"{{range .Fields}}" +
		"// Where{{.CodeName}} is the deleter where clause for {{.CodeName}}.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Where{{.CodeName}}(pred p.Predicate) " +
		"*{{$.Name}}D {\n" +
		"{{$.VarName}}.Fields = " +
		"append({{$.VarName}}.Fields, p.WhereClause(\"{{.Name}}\", pred))\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +
		"{{end}}" +

		// Delete clause for the deleter
		"// Delete the actual node\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete() *{{$.Name}}D {\n" +
		"{{$.VarName}}.WillDelete = true\n" +
		"return {{$.VarName}}\n" +
		"}\n\n" +

		// Traverse the graph
		"// Delete{{.FromNode}} traverses the deleter to the {{.FromNode}} " +
		"node.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete{{.FromNode}}() *{{.FromNode}}D " +
		"{\n" +
		"\tdeleter := {{.FromNode}}Deleter()\n" +
		"\tdeleter.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tdeleter.Prev = &{{$.VarName}}.Deleter\n" +
		"\t{{$.VarName}}.Next = &deleter.Deleter\n" +
		"\treturn deleter\n" +
		"}\n\n" +

		"// Delete{{.ToNode}} traverses the deleter to the {{.ToNode}} node.\n" +
		"func ({{$.VarName}} *{{$.Name}}D) Delete{{.ToNode}}() *{{.ToNode}}D " +
		"{\n" +
		"\tdeleter := {{.ToNode}}Deleter()\n" +
		"\tdeleter.Prefix = {{$.VarName}}.Prefix + 1\n" +
		"\tdeleter.Prev = &{{$.VarName}}.Deleter\n" +
		"\t{{$.VarName}}.Next = &deleter.Deleter\n" +
		"\treturn deleter\n" +
		"}\n"

	return cg.ExecTemplate(template, "edge_deleter", data, nil)
}
