package graphql

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/format"
	cg "splits-go-schema-codegen/codegen"
	"strings"
)

// WriteGraphQLNodeType writes the graphql base node type.
func WriteGraphQLNodeType(
	schemas []cg.Schema,
	schema cg.GraphQLSchema,
	manualParts []string,
	packageName string,
) string {

	getManualPart := initManualPart(manualParts)

	// Use templates to generate the node
	sections := []string{}
	sections = append(sections, GetGQLNodeFileHeaderCommentStr())
	sections = append(sections, GetGQLNodePackageStr(packageName))
	sections = append(sections, GetGQLNodeImportStr(getManualPart()))
	sections = append(sections, GetGQLNodeExtraFunctionsStr(getManualPart()))
	sections = append(sections, GetGQLNodeGeneratedFunctionsTagStr())
	sections = append(sections, GetGQLNodeInterfaceAndResolverStr(schema))
	sections = append(sections, GetGQLNodeRootQueryStr(schema))
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

// GetGQLNodeFileHeaderCommentStr generates an autogenerated tag.
func GetGQLNodeFileHeaderCommentStr() string {
	data := struct{}{}
	template := "// Autogenerated node type - regenerate with splits-go-schema-" +
		"codegen\n// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "node_type_file_header_comment", data, nil)
}

// GetGQLNodePackageStr generates the package string.
func GetGQLNodePackageStr(packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "node_type_package_string", data, nil)
}

// GetGQLNodeImportStr generates the import block.
func GetGQLNodeImportStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template := "import (\n" +
		"\t\"context\"\n" +
		"\t\"errors\"\n" +
		"\n" +
		"\tgraphql \"github.com/neelance/graphql-go\"\n" +
		"\n" +
		cg.StartManual + "\n" +
		"{{.ManualPart}}\n" +
		cg.EndManual + "\n" +
		")\n"
	return cg.ExecTemplate(template, "node_type_import_string", data, nil)
}

// GetGQLNodeExtraFunctionsStr adds a manual sections for user defined
// functions.
func GetGQLNodeExtraFunctionsStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template :=
		cg.StartManual + "\n" +
			"{{.ManualPart}}\n" +
			cg.EndManual + "\n"
	return cg.ExecTemplate(template, "node_type_extra_functions", data, nil)
}

// GetGQLNodeGeneratedFunctionsTagStr writes a generated functions tagline.
func GetGQLNodeGeneratedFunctionsTagStr() string {
	return "// === GENERATED FUNCTIONS === \n"
}

// GetGQLNodeInterfaceAndResolverStr writes the node interface and resolver
// types.
func GetGQLNodeInterfaceAndResolverStr(s cg.GraphQLSchema) string {
	data := struct {
		Nodes []cg.GraphQLNode
		Edges []cg.GraphQLEdge
	}{
		Nodes: s.Nodes,
		Edges: s.Edges,
	}
	template := "// Node interface represents a generic node in the graph.\n" +
		"type Node interface {\n" +
		"\tID(context.Context) (graphql.ID, error)\n" +
		"}\n" +
		"\n" +
		"// NodeResolver is the graphql resolver for a node.\n" +
		"type NodeResolver struct {\n" +
		"\tNode\n" +
		"}\n" +
		"\n" +
		"{{range .Nodes}}" +
		"// To{{.CodeName}} converts the generic node resolver to more specific " +
		"one.\n" +
		"func (n *NodeResolver) To{{.CodeName}}() (*{{.CodeName}}Resolver, " +
		"bool) {\n" +
		"\tres, ok := n.Node.(*{{.CodeName}}Resolver)\n" +
		"\treturn res, ok\n" +
		"}\n" +
		"\n" +
		"{{end}}"
	return cg.ExecTemplate(template, "node_type_parse_schema", data, nil)
}

// GetGQLNodeRootQueryStr generates the function that is the node root query.
func GetGQLNodeRootQueryStr(s cg.GraphQLSchema) string {
	data := struct {
		Nodes []cg.GraphQLNode
		Edges []cg.GraphQLEdge
	}{
		Nodes: s.Nodes,
		Edges: s.Edges,
	}
	template := "// Node is the root query resolver for a specific node.\n" +
		"func (r *Resolver) Node(ctx context.Context, args idArg) " +
		"(*NodeResolver, error) {\n" +
		"\t\n" +
		"\t// Demux the id\n" +
		"\tkind, id, err := demuxKindAndID(args.ID)\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\t\n" +
		"\tswitch kind {\n" +
		"{{range .Nodes}}" +
		"\tcase \"{{.CodeName}}\":\n" +
		"\t\treturn &NodeResolver{&{{.CodeName}}Resolver{id}}, nil\n" +
		"{{end}}" +
		"\t}\n" +
		"\treturn nil, errors.New(\"invalid node type: \" + kind)" +
		"}\n"
	return cg.ExecTemplate(template, "node_type_root_query", data, nil)
}
