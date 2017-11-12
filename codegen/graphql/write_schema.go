// The writer assumes a valid schema. The generation involves creating manual
// sections of code, as well as implementing getters that respect privacy for
// all the fields.

package graphql

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/format"
	cg "splits-go-schema-codegen/codegen"
	"strings"
	t "text/template"
)

// WriteGraphQLSchema writes the graphql base schema.
func WriteGraphQLSchema(
	schemas []cg.Schema,
	schema cg.GraphQLSchema,
	manualParts []string,
	packageName string,
) string {

	getManualPart := initManualPart(manualParts)

	// Use templates to generate the node
	sections := []string{}
	sections = append(sections, GetSchemaFileHeaderCommentStr())
	sections = append(sections, GetSchemaPackageStr(packageName))
	sections = append(sections, GetSchemaImportStr(getManualPart()))
	sections = append(sections, GetSchemaExtraFunctionsStr(getManualPart()))
	sections = append(sections, GetSchemaGeneratedFunctionsTagStr())
	sections = append(sections, GetSchemaParseSchemaStr())
	sections = append(sections, GetSchemaStringStr(schema))
	result := strings.Join(sections, "\n")
	res, err := format.Source([]byte(result))
	if err != nil {
		fmt.Println(result)
		panic(err)
	}

	signatureRes := []byte(cg.ReplaceAllStringSubmatchFunc(
		cg.ManualExtractor,
		string(res),
		func(groups []string) string {
			return cg.StartManual + groups[2] + cg.EndManual
		},
	))

	// Generate the MD5 signature
	sum := md5.Sum([]byte(signatureRes))
	signature := hex.EncodeToString([]byte(sum[:]))

	// Add the signature to the top of the file
	return "// @SignedSource (" + signature + ")\n" + string(res)
}

// GetSchemaFileHeaderCommentStr generates an autogenerated tag.
func GetSchemaFileHeaderCommentStr() string {
	data := struct{}{}
	template := "// Autogenerated schema - regenerate with splits-go-schema-" +
		"codegen\n// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "schema_file_header_comment", data, nil)
}

// GetSchemaPackageStr generates the package string.
func GetSchemaPackageStr(packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "schema_package_string", data, nil)
}

// GetSchemaImportStr generates the import block.
func GetSchemaImportStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template := "import (\n" +
		"\t\"splits-go-api/api/graphql/resolvers\"\n" +
		"\n" +
		"\tgraphql \"github.com/neelance/graphql-go\"\n" +
		"\n" +
		cg.StartManual + "\n" +
		"{{.ManualPart}}\n" +
		cg.EndManual + "\n" +
		")\n"
	return cg.ExecTemplate(template, "schema_import_string", data, nil)
}

// GetSchemaExtraFunctionsStr adds a manual sections for user defined functions.
func GetSchemaExtraFunctionsStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template :=
		cg.StartManual + "\n" +
			"{{.ManualPart}}\n" +
			cg.EndManual + "\n"
	return cg.ExecTemplate(template, "schema_extra_functions", data, nil)
}

// GetSchemaGeneratedFunctionsTagStr writes a generated functions tagline.
func GetSchemaGeneratedFunctionsTagStr() string {
	return "// === GENERATED FUNCTIONS === \n"
}

// GetSchemaParseSchemaStr writes the parse schema function.
func GetSchemaParseSchemaStr() string {
	template := "var schema *graphql.Schema\n" +
		"\n" +
		"// ParseSchema at startup to check for schema issues.\n" +
		"func ParseSchema() {\n" +
		"\tschema = graphql.MustParseSchema(Schema, &resolvers.Resolver{})\n" +
		"}\n"
	return cg.ExecTemplate(template, "schema_parse_schema", nil, nil)
}

// GetSchemaStringStr returns the string form of the graphql schema.
func GetSchemaStringStr(s cg.GraphQLSchema) string {

	edges := []cg.GraphQLEdge{}
	edgeMap := map[string]bool{}
	for _, e := range s.Edges {
		name := e.From + e.TotalName + e.To
		if _, ok := edgeMap[name]; !ok {
			edges = append(edges, e)
			edgeMap[name] = true
		}
	}

	data := struct {
		Nodes []cg.GraphQLNode
		Edges []cg.GraphQLEdge
	}{
		Nodes: s.Nodes,
		Edges: edges,
	}
	funcMap := t.FuncMap{
		"ToLower": strings.ToLower,
	}
	template := "// Schema of the graphql api.\n" +
		"var Schema = `\n" +
		"\n" +
		"scalar Time\n" +
		"\n" +
		"schema {\n" +
		"\tquery: Query\n" +
		"\tmutation: Mutation\n" +
		"}\n" +
		"\n" +
		"# The Query type represents all the entry points into the graph.\n" +
		"type Query {\n" +
		"\tnode(id: ID!): Node\n" +
		"\tviewer: User\n" +
		"{{range .Nodes}}" +
		"\t{{.Name | ToLower}}(id: ID!): {{.Name}}\n" +
		"{{end}}" +
		"}\n" +
		"\n" +
		"# The Node represents a generic node in the graph.\n" +
		"interface Node {\n" +
		"\t# The ID of the node.\n" +
		"\tid: ID!\n" +
		"}\n" +
		"{{range .Nodes}}" +
		"\n" +
		"# {{.Description}}\n" +
		"type {{.Name}} implements Node {\n" +
		"{{range .Fields}}" +
		"\t# {{.Description}}\n" +
		"\t{{.Name}}: {{.Type}}\n" +
		"{{end}}" +
		"{{range .Edges}}" +
		"\n\t# {{.Description}}\n" +
		"\t{{.FieldName}}(first: Int, after: ID, orderBy: [OrderBy!]): " +
		"{{.From}}To{{.To}}Connection!\n" +
		"{{end}}" +
		"}\n" +
		"{{end}}" +
		"{{range .Edges}}" +
		"\n" +
		"type {{.From}}To{{.To}}Connection {\n" +
		"\ttotalCount: Int!\n" +
		"\tedges: [{{.From}}To{{.To}}Edge]\n" +
		"\tnodes: [{{.To}}]\n" +
		"\tpageInfo: PageInfo!\n" +
		"}\n" +
		"\n" +
		"type {{.From}}To{{.To}}Edge {\n" +
		"\tcursor: ID!\n" +
		"\tnode: {{.To}}\n" +
		"\t\n" +
		"{{range .Fields}}" +
		"\t# {{.Description}}\n" +
		"\t{{.Name}}: {{.Type}}\n" +
		"{{end}}" +
		"}\n" +
		"{{end}}" +
		"\n" +
		"input OrderBy {\n" +
		"\tfield: String!\n" +
		"\tdesc: Boolean!\n" +
		"}\n" +
		"\n" +
		"# Information for paginating connections.\n" +
		"type PageInfo {\n" +
		"\tstartCursor: ID\n" +
		"\tendCursor: ID\n" +
		"\thasNextPage: Boolean!\n" +
		"\thasPreviousPage: Boolean!\n" +
		"}\n" +
		"` + resolvers.MutationSchema\n"
	return cg.ExecTemplate(template, "schema_string", data, funcMap)
}
