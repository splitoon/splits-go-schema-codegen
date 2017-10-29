package graphql

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/format"
	cg "splits-go-schema-codegen/codegen"
	"strings"
)

// WriteDataloaderBatcher writes the dataloader batcher muxing.
func WriteDataloaderBatcher(
	schemas []cg.Schema,
	schema cg.GraphQLSchema,
	manualParts []string,
	packageName string,
) string {

	getManualPart := initManualPart(manualParts)

	// Use templates to generate the node
	sections := []string{}
	sections = append(sections, GetDLBatcherFileHeaderCommentStr())
	sections = append(sections, GetDLBatcherPackageStr(packageName))
	sections = append(sections, GetDLBatcherImportStr(getManualPart()))
	sections = append(sections, GetDLBatcherExtraFunctionsStr(getManualPart()))
	sections = append(sections, GetDLBatcherGeneratedFunctionsTagStr())
	sections = append(sections, GetDLBatcherBatcherStr(schema, getManualPart()))
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

// GetDLBatcherFileHeaderCommentStr generates an autogenerated tag.
func GetDLBatcherFileHeaderCommentStr() string {
	data := struct{}{}
	template := "// Autogenerated dataloader batcher - regenerate with " +
		"splits-go-schema-" +
		"codegen\n// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "dl_batcher_file_header_comment", data, nil)
}

// GetDLBatcherPackageStr generates the package string.
func GetDLBatcherPackageStr(packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "dl_batcher_package_string", data, nil)
}

// GetDLBatcherImportStr generates the import block.
func GetDLBatcherImportStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template := "import (\n" +
		"\t\"context\"\n" +
		"\t\"splits-go-api/auth/contexts\"\n" +
		"\t\"splits-go-api/constants\"\n" +
		"\t\"splits-go-api/db\"\n" +
		"\t\"splits-go-api/logic\"\n" +
		"\t\"splits-go-api/logic/util\"\n" +
		"\t\"strings\"\n" +
		"\n" +
		cg.StartManual + "\n" +
		"{{.ManualPart}}\n" +
		cg.EndManual + "\n" +
		")\n"
	return cg.ExecTemplate(template, "dl_batcher_import_string", data, nil)
}

// GetDLBatcherExtraFunctionsStr adds a manual sections for user defined
// functions.
func GetDLBatcherExtraFunctionsStr(manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template :=
		cg.StartManual + "\n" +
			"{{.ManualPart}}\n" +
			cg.EndManual + "\n"
	return cg.ExecTemplate(template, "dl_batcher_extra_functions", data, nil)
}

// GetDLBatcherGeneratedFunctionsTagStr writes a generated functions tagline.
func GetDLBatcherGeneratedFunctionsTagStr() string {
	return "// === GENERATED FUNCTIONS === \n"
}

// GetDLBatcherBatcherStr writes the batcher function.
func GetDLBatcherBatcherStr(s cg.GraphQLSchema, manualPart string) string {
	edgeFields := []cg.GraphQLEdge{}
	edgeFieldMap := map[string]bool{}
	for _, e := range s.Edges {
		if _, ok := edgeFieldMap[e.TotalName]; !ok {
			edgeFields = append(edgeFields, e)
			edgeFieldMap[e.TotalName] = true
		}
	}
	edges := []cg.GraphQLEdge{}
	edgeMap := map[string]bool{}
	for _, e := range s.Edges {
		if _, ok := edgeMap[e.FromCodeName+e.ToCodeName]; !ok {
			edges = append(edges, e)
			edgeMap[e.FromCodeName+e.ToCodeName] = true
		}
	}
	data := struct {
		Nodes      []cg.GraphQLNode
		EdgeFields []cg.GraphQLEdge
		Edges      []cg.GraphQLEdge
		ManualPart string
	}{
		Nodes:      s.Nodes,
		EdgeFields: edgeFields,
		Edges:      edges,
		ManualPart: manualPart,
	}
	template := "// Generates the batchers to be batched later.\n" +
		"func genBatchedQueries(\n" +
		"\tctx context.Context,\n" +
		"\tqueries map[string]map[string][]string,\n" +
		"\tfetchedData map[string]interface{},\n" +
		"\tfetchedErrs map[string]error,\n" +
		") ([]*util.LogicGetWrapper, map[int]map[int]string) {\n" +
		"\t\n" +
		"\tvc := ctx.Value(constants.VCKey).(contexts.ViewerContext)\n" +
		"\tconn := ctx.Value(constants.ConnKey).(*db.Conn)\n" +
		"\t\n" +
		"\t// Initialize\n" +
		"\tindex := 0 // Refers to the index of the query in the pipeline\n" +
		"\tbatchedMapper := map[int]map[int]string{}\n" +
		"\tbatchedQueries := []*util.LogicGetWrapper{}\n" +
		"\t\n" +
		"\tfor kind, idToFieldMap := range queries {\n" +
		"\t\tfor id, fields := range idToFieldMap {\n" +
		"\t\t\t\n" +
		"\t\t\tswitch kind {\n" +

		// Node fields

		"{{range .Nodes}}" +
		"\t\t\tcase \"{{.CodeName}}\":\n" +
		"\t\t\t\t{\n" +
		"\t\t\t\t\tb, err := logic.Get{{.CodeName}}ByIDBatcher(conn, vc, " +
		"ctx, id, fields)\n" +
		"\t\t\t\t\tif err != nil {\n" +
		"\t\t\t\t\t\taddFetchedErrors(kind, id, fields, err, " +
		"fetchedData, fetchedErrs)\n" +
		"\t\t\t\t\t} else {\n" +
		"\t\t\t\t\t\tbatchedQueries = addBatchers(b, batchedMapper, index, " +
		"kind, id,\n\t\t\t\t\t\t\t\tfields, batchedQueries)\n" +
		"\t\t\t\t\t\tindex++\n" +
		"\t\t\t\t\t}\n" +
		"\t\t\t\t}\n" +
		"{{end}}" +

		// Node to node

		"{{range .Edges}}" +
		"\t\t\tcase \"{{.FromCodeName}}{{.ToCodeName}}\":\n" +
		"\t\t\t\t{\n" +
		"\t\t\t\t\tb, err := logic.Get{{.FromCodeName}}{{.FieldResolveName}}" +
		"Batcher(" +
		"conn, vc, ctx, id)\n" +
		"\t\t\t\t\tif err != nil {\n" +
		"\t\t\t\t\t\taddFetchedErrors(kind, id, fields, err, fetchedData, " +
		"fetchedErrs)\n" +
		"\t\t\t\t\t} else {\n" +
		"\t\t\t\t\t\tbatchedQueries = addBatchers(b, batchedMapper, index, " +
		"kind, id,\n\t\t\t\t\t\t\tfields, batchedQueries)\n" +
		"\t\t\t\t\t\tindex++\n" +
		"\t\t\t\t\t}\n" +
		"\t\t\t\t}\n" +
		"{{end}}" +

		// Edge fields

		"{{range .EdgeFields}}" +
		"\t\t\tcase \"{{.TotalName}}\":\n" +
		"\t\t\t\t{\n" +
		"\t\t\t\t\tids := strings.Split(id, \"|\")\n" +
		"\t\t\t\t\tb, err := logic.Get{{.EdgeCodeName}}ByIDsBatcher(conn, vc, " +
		"\n\t\t\t\t\t\tctx, ids[0], ids[1], fields)\n" +
		"\t\t\t\t\tif err != nil {\n" +
		"\t\t\t\t\t\taddFetchedErrors(kind, id, fields, err, " +
		"fetchedData, fetchedErrs)\n" +
		"\t\t\t\t\t} else {\n" +
		"\t\t\t\t\t\tbatchedQueries = addBatchers(b, batchedMapper, index, " +
		"kind, id,\n\t\t\t\t\t\t\tfields, batchedQueries)\n" +
		"\t\t\t\t\t\tindex++\n" +
		"\t\t\t\t\t}\n" +
		"\t\t\t\t}\n" +
		"{{end}}" +
		"\t\t\t" + cg.StartManual + "\n" +
		"\t\t\t{{.ManualPart}}\n" +
		"\t\t\t" + cg.EndManual + "\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\treturn batchedQueries, batchedMapper\n" +
		"}\n"
	return cg.ExecTemplate(template, "node_type_parse_schema", data, nil)
}