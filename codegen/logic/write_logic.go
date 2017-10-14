// The writer assumes a valid schema. The generation involves creating manual
// sections of code, as well as implementing getters that respect privacy for
// all the fields.

package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/format"
	"splits-go-api/auth/policies"
	cg "splits-go-schema-codegen/codegen"
	"strings"
)

func initManualPart(manualParts []string) func() string {
	index := 0
	return func() string {
		if manualParts == nil {
			return ""
		}
		if index >= len(manualParts) {
			return ""
		}
		index++
		return strings.Replace(manualParts[index-1], "%", "%%", -1)
	}
}

// WriteSchemaLogicNode writes the logic for a node.
func WriteSchemaLogicNode(
	s cg.Schema,
	manualParts []string,
	packageName string,
) string {

	getManualPart := initManualPart(manualParts)

	// Use templates to generate the node
	sections := make([]string, 0, 11)
	sections = append(sections, GetFileHeaderCommentStr(s))
	sections = append(sections, GetPackageStr(s, packageName))
	sections = append(sections, GetNodeImportStr(s, getManualPart()))
	sections = append(sections, GetExtraFunctionsStr(s, getManualPart()))
	sections = append(sections, GetGeneratedFunctionsTagStr())
	sections = append(sections, GetNodeCheckAuthStr(s))
	sections = append(sections, GetNodeFieldQueryStr(s))
	sections = append(sections, GetNodeGetByIDStr(s))
	sections = append(sections, GetNodeGetByIDBatchStr(s))
	sections = append(sections, GetNodeConnectedNodesStr(s))
	// sections = append(sections, GetNodeStr(s))
	// sections = append(sections, GetNodeQueryStructStr(s))
	// sections = append(sections, GetNodeQueryConstructorStr(s))
	// sections = append(sections, GetNodeQueryWhereStr(s))
	// sections = append(sections, GetNodeQueryReturnStr(s))
	// sections = append(sections, GetNodeQueryOrderStr(s))
	// sections = append(sections, GetNodeQueryEdgesStr(s))
	// sections = append(sections, GetNodeMutatorStr(s))
	// sections = append(sections, GetNodeDeleterStr(s))
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

// WriteSchemaLogicEdge writes the logic for a n edge.
func WriteSchemaLogicEdge(
	s cg.Schema,
	e cg.EdgeStruct,
	manualParts []string,
	packageName string,
) string {
	getManualPart := initManualPart(manualParts)

	// Use templates to generate the node
	sections := make([]string, 0, 11)
	sections = append(sections, GetFileHeaderCommentStr(s))
	sections = append(sections, GetPackageStr(s, packageName))
	sections = append(sections, GetEdgeImportStr(s, getManualPart()))
	sections = append(sections, GetExtraFunctionsStr(s, getManualPart()))
	sections = append(sections, GetGeneratedFunctionsTagStr())
	sections = append(sections, GetEdgeCheckAuthStr(s, e))
	sections = append(sections, GetEdgeFieldQueryStr(s, e))
	sections = append(sections, GetEdgeGetByIDStr(s, e))
	sections = append(sections, GetEdgeGetByIDBatcherStr(s, e))
	sections = append(sections, GetEdgeGetByIDsStr(s, e))
	sections = append(sections, GetEdgeGetByIDsBatcherStr(s, e))
	// sections = append(sections, GetNodeStr(s))
	// sections = append(sections, GetNodeQueryStructStr(s))
	// sections = append(sections, GetNodeQueryConstructorStr(s))
	// sections = append(sections, GetNodeQueryWhereStr(s))
	// sections = append(sections, GetNodeQueryReturnStr(s))
	// sections = append(sections, GetNodeQueryOrderStr(s))
	// sections = append(sections, GetNodeQueryEdgesStr(s))
	// sections = append(sections, GetNodeMutatorStr(s))
	// sections = append(sections, GetNodeDeleterStr(s))
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

// GetFileHeaderCommentStr generates an autogenerated tag.
func GetFileHeaderCommentStr(s cg.Schema) string {
	data := struct {
		Name string
	}{
		Name: s.GetName(),
	}
	template := "// Autogenerated {{.Name}} - regenerate with splits-go-schema-" +
		"codegen\n// Force autogen by deleting the @SignedSource line.\n"
	return cg.ExecTemplate(template, "file_header_comment", data)
}

// GetPackageStr generates the package string.
func GetPackageStr(s cg.Schema, packageName string) string {
	data := struct {
		Package string
	}{
		Package: packageName,
	}
	template := "package {{.Package}}\n"
	return cg.ExecTemplate(template, "package_string", data)
}

// GetExtraFunctionsStr adds a manual sections for user defined functions.
func GetExtraFunctionsStr(s cg.Schema, manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template :=
		cg.StartManual + "\n" +
			"{{.ManualPart}}\n" +
			cg.EndManual + "\n"

	return cg.ExecTemplate(template, "extra_functions", data)
}

// GetGeneratedFunctionsTagStr writes a generated functions tagline.
func GetGeneratedFunctionsTagStr() string {
	return "// === GENERATED FUNCTIONS === \n"
}

// =============================================================================
// Node
// =============================================================================

// GetNodeImportStr generates the import block.
func GetNodeImportStr(s cg.Schema, manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template := "import (\n" +
		"\t\"errors\"\n" +
		"\t\"splits-go-api/auth/contexts\"\n" +
		"\t\"splits-go-api/auth/policies\"\n" +
		"\t\"splits-go-api/constants\"\n" +
		"\t\"splits-go-api/auth/rules\"\n" +
		"\t\"splits-go-api/db\"\n" +
		"\t\"splits-go-api/db/models\"\n" +
		"\tp \"splits-go-api/db/models/predicates\"\n" +
		"\t\"splits-go-api/logic/privacy\"\n" +
		"\t\"splits-go-api/logic/util\"\n" +
		"\n" +
		"\t\"context\"\n" +
		"\t\"sync\"\n" +
		"\t\"time\"\n" +
		"\n" +
		"\tbolt \"github.com/johnnadratowski/golang-neo4j-bolt-driver\"\n" +
		"\te \"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors\"\n" +
		cg.StartManual + "\n" +
		"{{.ManualPart}}\n" +
		cg.EndManual + "\n" +
		")\n"
	return cg.ExecTemplate(template, "node_import_string", data)
}

// GetNodeCheckAuthStr creates the function that checks authorization of a vc.
func GetNodeCheckAuthStr(s cg.Schema) string {
	data := struct {
		Name string
	}{
		Name: s.GetName(),
	}
	template := strings.Join([]string{
		"func check{{.Name}}Auth(",
		"\tconn *bolt.Conn,",
		"\tvc contexts.ViewerContext,",
		"\tpp policies.PrivacyPolicy,",
		"\tparams context.Context,",
		"\tid string,",
		") (bool, error) {",
		"",
		"\t// Check context cache",
		"\tpermMap := params.Value(constants.PermsKey). " +
			"(map[string]map[string]bool)",
		"\tmutex := params.Value(constants.PermsMutexKey).(*sync.Mutex)",
		"\tmutex.Lock()",
		"\tperms, ok := permMap[id]",
		"\tif ok {",
		"\t\tif hasPerm, ok2 := perms[pp.GetName()]; ok2 {",
		"\t\t\tmutex.Unlock()",
		"\t\t\treturn hasPerm, nil",
		"\t\t}",
		"\t} else {",
		"\t\tpermMap[id] = map[string]bool{}",
		"\t}",
		"\tmutex.Unlock()",
		"",
		"\tauthContext := rules.AuthContext{SrcID: id, Conn: conn}",
		"\thasAuth, _, err := pp.CheckAuth(&vc, " +
			"authContext)",
		"\tif err != nil {",
		"\treturn false, err",
		"\t}",
		"\tmutex.Lock()",
		"\tperms = permMap[id]",
		"\tperms[pp.GetName()] = hasAuth",
		"\tmutex.Unlock()",
		"",
		"\treturn hasAuth, nil",
		"}",
	}, "\n") + "\n"
	return cg.ExecTemplate(template, "node_check_auth", data)
}

// GetNodeFieldQueryStr creates a function that generates a query for the
// sppecified fields.
func GetNodeFieldQueryStr(s cg.Schema) string {
	fields := s.GetFields()
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.FieldStruct
	}{
		Name:     s.GetName(),
		Policies: pp,
		Fields:   fields,
	}
	template := "func create{{.Name}}FieldQuery(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\tfields []string,\n" +
		"\tq *models.{{.Name}}Q,\n" +
		") (*models.{{.Name}}Q, []bool, error) {\n" +
		"\n" +
		"\t// Check the auth for the fields\n" +
		"\tfieldCheck := make([]bool, len(fields))\n" +
		"\tfor i := range fields {\n" +
		"\tfieldCheck[i] = true\n" +
		"\t}" +
		"\n" +
		"\t// Add the fields to the query if appropriate auth\n" +
		"\tfor i, x := range fields {\n" +
		"\t\tswitch x {\n\n" +
		"{{range .Fields}}" +
		"\t\tcase \"{{.Name}}\":\n" +
		"\t\t\t{\n" +
		"\t\t\t\thasAuth, err := check{{$.Name}}Auth(conn, vc, " +
		"privacy.{{.Privacy.GetName}}, params, id)\n" +
		"\t\t\t\tif err != nil {\n" +
		"\t\t\t\t\treturn nil, nil, err\n" +
		"\t\t\t\t}\n" +
		"\t\t\t\tif hasAuth {\n" +
		"\t\t\t\t\tq = q.Return{{.CodeName}}()\n" +
		"\t\t\t\t} else {\n" +
		"\t\t\t\t\tfieldCheck[i] = false\n" +
		"\t\t\t\t}\n" +
		"\t\t\t}\n" +
		"{{end}}" +
		"\t\tdefault:\n" +
		"\t\t\t{\n" +
		"\t\t\t\tfieldCheck[i] = false\n" +
		"\t\t\t\tlog.Warnf(\"invalid requested field: %%s-%%s\", \"{{$.Name}}\", " +
		"x)\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\treturn q, fieldCheck, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "node_field_query", data)
}

// GetNodeGetByIDStr gets the function that retrieves fields by the id of the
// node.
func GetNodeGetByIDStr(s cg.Schema) string {

	fields := s.GetFields()
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}

	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.FieldStruct
	}{
		Name:     s.GetName(),
		Policies: pp,
		Fields:   fields,
	}
	template := "// Get{{.Name}}ByID retrives the fields of a specific " +
		"{{.Name}}.\n" +
		"// If there is insufficient authorization, the field will return null.\n" +
		"func Get{{.Name}}ByID(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\tfields []string,\n" +
		") ([]interface{}, error) {\n" +
		"\n" +
		"\t// Generate the query\n" +
		"\t q := models.{{.Name}}Query().\n" +
		"\t\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(conn, vc, params, id, " +
		"fields, q)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\t// Return nil if no fields to request\n" +
		"\tif len(q.Return) == 0 {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n" +
		"\n" +
		"\t// Execute the query\n" +
		"\tvar row []interface{}\n" +
		"\tnewConn := conn\n" +
		"\tfor i := 0; row == nil && i < constants.LogicRetryCount; i++ {\n" +
		"\t\trow, err = q.GenOne(newConn)\n" +
		"\t\tif _, isBoltErr := err.(*e.Error); isBoltErr {\n" +
		"\t\t\t// Try a new connection\n" +
		"\t\t\t(*newConn).Close()\n" +
		"\t\t\ttime.Sleep(time.Millisecond * constants.LogicRetryWait)\n" +
		"\t\t\tc, _ := db.GetDriverConn()\n" +
		"\t\t\tnewConn = c\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\t*conn = *newConn\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\n" +
		"\t// Check for the authed fields\n" +
		"\tresults := util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\n" +
		"\treturn results, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "node_by_id", data)
}

func GetNodeGetByIDBatchStr(s cg.Schema) string {
	fields := s.GetFields()
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}

	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.FieldStruct
	}{
		Name:     s.GetName(),
		Policies: pp,
		Fields:   fields,
	}
	template := "// Get{{.Name}}ByIDBatcher wraps the Get{{.Name}}ByID request " +
		"to be batched later.\n" +
		"func Get{{.Name}}ByIDBatcher(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\tfields []string,\n" +
		") (*util.LogicGetWrapper, error) {\n" +
		"\n" +
		"\t// Generate the query\n" +
		"\t q := models.{{.Name}}Query().\n" +
		"\t\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(conn, vc, params, id, " +
		"fields, q)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\t// Return nil if no fields to request\n" +
		"\tif len(q.Return) == 0 {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n" +
		"\n" +
		"\tbatcher := new(util.LogicGetWrapper)\n" +
		"\tbatcher.Query = &q.Query\n" +
		"\tbatcher.EvalAuth = func(row []interface{}) []interface{} {\n" +
		"\t\treturn util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\t}\n" +
		"\treturn batcher, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "node_by_id_batch", data)
}

// GetNodeConnectedNodesStr generates the function that gets connected
// corresponding node ids. This also writes the corresponding batch wrapper.
func GetNodeConnectedNodesStr(s cg.Schema) string {

	type NamePrivacyPair struct {
		Name    string
		Privacy policies.PrivacyPolicy
	}

	// Extract the edge name to the node name
	edges := map[string]NamePrivacyPair{}
	for _, e := range s.GetEdges() {
		if e.ToNode.GetName() == s.GetName() { // group->user
			edges[e.CodeName] = NamePrivacyPair{e.FromNode.GetName(),
				e.ReversePrivacy}
		} else if e.FromNode.GetName() == s.GetName() {
			edges[e.CodeName] = NamePrivacyPair{e.ToNode.GetName(), e.Privacy}
		}
	}
	for _, e := range s.GetEdgePointers() {
		if e.ToNode.GetName() == s.GetName() { // group->user
			edges[e.CodeName] = NamePrivacyPair{e.FromNode.GetName(),
				e.ReversePrivacy}
		} else if e.FromNode.GetName() == s.GetName() {
			edges[e.CodeName] = NamePrivacyPair{e.ToNode.GetName(), e.Privacy}
		}
	}

	data := struct {
		Name       string
		EdgeToNode map[string]NamePrivacyPair
	}{
		Name:       s.GetName(),
		EdgeToNode: edges,
	}
	template := "{{range $edgeName, $value := .EdgeToNode}}" +
		"// Get{{$.Name}}{{$value.Name}}s retrieves the ids of connected " +
		"{{$value.Name}}s.\n" +
		"func Get{{$.Name}}{{$value.Name}}s(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		") ([]interface{}, error) {\n" +
		"\n" +
		"\t// Check auth\n" +
		"\thasAuth, err := check{{$.Name}}Auth(conn, vc, " +
		"privacy.{{$value.Privacy.GetName}}, params, id)\n" +
		"\tif err != nil{\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif !hasAuth {\n" +
		"\t\treturn nil, errors.New(\"invalid auth for " +
		"Get{{$.Name}}{{$value.Name}}\")\n" +
		"\t}" +
		"\n" +
		"\t// Build the query and execute it\n" +
		"\trows, stmt, err := models.{{$.Name}}Query().\n" +
		"\t\tWhereID(p.Equals(id)).\n" +
		"\t\tQuery{{$edgeName}}().\n" +
		"\t\tQuery{{$value.Name}}().\n" +
		"\t\tReturnID().\n" +
		"\t\tGen(conn)\n" +
		"\n" +
		"\tif stmt != nil {\n" +
		"\t\tdefer stmt.Close()\n" +
		"\t}\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\n" +
		"\tids, err := util.ExtractFirstFromRows(rows)\n" +
		"\treturn ids, err\n" +
		"}\n\n" +
		"// Get{{$.Name}}{{$value.Name}}sBatcher wraps the Get{{$.Name}}" +
		"{{$value.Name}}s request to be batched later.\n" +
		"func Get{{$.Name}}{{$value.Name}}sBatcher(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		") (*util.LogicGetWrapper, error) {\n" +
		"\n" +
		"\t// Check auth\n" +
		"\thasAuth, err := check{{$.Name}}Auth(conn, vc, " +
		"privacy.{{$value.Privacy.GetName}}, params, id)\n" +
		"\tif err != nil{\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif !hasAuth {\n" +
		"\t\treturn nil, errors.New(\"invalid auth for " +
		"Get{{$.Name}}{{$value.Name}}\")\n" +
		"\t}" +
		"\n" +
		"\tq := models.{{$.Name}}Query().\n" +
		"\t\tWhereID(p.Equals(id)).\n" +
		"\t\tQuery{{$edgeName}}().\n" +
		"\t\tQuery{{$value.Name}}().\n" +
		"\t\tReturnID()\n" +
		"\n" +
		"\tbatcher := new(util.LogicGetWrapper)\n" +
		"\tbatcher.Query = &q.Query\n" +
		"\tbatcher.EvalAuth = func(row []interface{}) []interface{} {\n" +
		"\t\treturn row\n" +
		"\t}\n" +
		"\treturn batcher, nil\n" +
		"}\n\n" +
		"{{end}}"

	return cg.ExecTemplate(template, "node_connected_nodes", data)
}

// =============================================================================
// Edges
// =============================================================================

// GetEdgeImportStr generates the import block.
func GetEdgeImportStr(s cg.Schema, manualPart string) string {
	data := struct {
		ManualPart string
	}{
		ManualPart: manualPart,
	}
	template := "import (\n" +
		"\t\"splits-go-api/auth/contexts\"\n" +
		"\t\"splits-go-api/auth/policies\"\n" +
		"\t\"splits-go-api/constants\"\n" +
		"\t\"splits-go-api/auth/rules\"\n" +
		"\t\"splits-go-api/db\"\n" +
		"\t\"splits-go-api/db/models\"\n" +
		"\tp \"splits-go-api/db/models/predicates\"\n" +
		"\t\"splits-go-api/logic/privacy\"\n" +
		"\t\"splits-go-api/logic/util\"\n" +
		"\n" +
		"\t\"context\"\n" +
		"\t\"sync\"\n" +
		"\t\"time\"\n" +
		"\n" +
		"\tbolt \"github.com/johnnadratowski/golang-neo4j-bolt-driver\"\n" +
		"\te \"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors\"\n" +
		cg.StartManual + "\n" +
		"{{.ManualPart}}\n" +
		cg.EndManual + "\n" +
		")\n"
	return cg.ExecTemplate(template, "edge_import_string", data)
}

// GetEdgeCheckAuthStr generates the auth check for edge fields.
func GetEdgeCheckAuthStr(s cg.Schema, e cg.EdgeStruct) string {
	fromVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"
	data := struct {
		Name    string
		FromVar string
		ToVar   string
	}{
		Name:    e.CodeName,
		FromVar: fromVar,
		ToVar:   toVar,
	}
	template := strings.Join([]string{
		"func check{{.Name}}Auth(",
		"\tconn *bolt.Conn,",
		"\tvc contexts.ViewerContext,",
		"\tpp policies.PrivacyPolicy,",
		"\tparams context.Context,",
		"\t{{.FromVar}} string,",
		"\t{{.ToVar}} string,",
		") (bool, error) {",
		"",
		"\t// Check context cache",
		"\tcacheID := {{.FromVar}} + \"-\" +{{.ToVar}}",
		"\tpermMap := params.Value(constants.PermsKey). " +
			"(map[string]map[string]bool)",
		"\tmutex := params.Value(constants.PermsMutexKey).(*sync.Mutex)",
		"\tmutex.Lock()",
		"\tperms, ok := permMap[cacheID]",
		"\tif ok {",
		"\t\tif hasPerm, ok2 := perms[pp.GetName()]; ok2 {",
		"\t\t\tmutex.Unlock()",
		"\t\t\treturn hasPerm, nil",
		"\t\t}",
		"\t} else {",
		"\t\tpermMap[cacheID] = map[string]bool{}",
		"\t}",
		"\tmutex.Unlock()",
		"",
		"\tauthContext := rules.AuthContext{SrcID: {{.FromVar}}, DestID: " +
			"{{.ToVar}}, Conn: conn}",
		"\thasAuth, _, err := pp.CheckAuth(&vc, " +
			"authContext)",
		"\tif err != nil {",
		"\treturn false, err",
		"\t}",
		"\tmutex.Lock()",
		"\tperms = permMap[cacheID]",
		"\tperms[pp.GetName()] = hasAuth",
		"\tmutex.Unlock()",
		"",
		"\treturn hasAuth, nil",
		"}",
	}, "\n") + "\n"
	return cg.ExecTemplate(template, "edge_check_auth", data)
}

// GetEdgeFieldQueryStr creates a function that generates a query for the
// sppecified fields.
func GetEdgeFieldQueryStr(s cg.Schema, e cg.EdgeStruct) string {
	fields := e.Fields
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	fromVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"

	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.EdgeFieldStruct
		FromVar  string
		ToVar    string
	}{
		Name:     e.CodeName,
		Policies: pp,
		Fields:   fields,
		FromVar:  fromVar,
		ToVar:    toVar,
	}
	template := "func create{{.Name}}FieldQuery(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\t{{.FromVar}} string,\n" +
		"\t{{.ToVar}} string,\n" +
		"\tfields []string,\n" +
		"\tq *models.{{.Name}}Q,\n" +
		") (*models.{{.Name}}Q, []bool, error) {\n" +
		"\n" +
		"\t// Check the auth for the fields\n" +
		"\tfieldCheck := make([]bool, len(fields))\n" +
		"\tfor i := range fields {\n" +
		"\tfieldCheck[i] = true\n" +
		"\t}" +
		"\n" +
		"\t// Add the fields to the query if appropriate auth\n" +
		"\tfor i, x := range fields {\n" +
		"\t\tswitch x {\n\n" +
		"{{range .Fields}}" +
		"\t\tcase \"{{.Name}}\":\n" +
		"\t\t\t{\n" +
		"\t\t\t\thasAuth, err := check{{$.Name}}Auth(\n\t\tconn,\n\t\tvc, " +
		"\n\t\tprivacy.{{.Privacy.GetName}},\n\t\tparams,\n\t\t{{$.FromVar}}, " +
		"\n\t\t{{$.ToVar}},\n)\n" +
		"\t\t\t\tif err != nil {\n" +
		"\t\t\t\t\treturn nil, nil, err\n" +
		"\t\t\t\t}\n" +
		"\t\t\t\tif hasAuth {\n" +
		"\t\t\t\t\tq = q.Return{{.CodeName}}()\n" +
		"\t\t\t\t} else {\n" +
		"\t\t\t\t\tfieldCheck[i] = false\n" +
		"\t\t\t\t}\n" +
		"\t\t\t}\n" +
		"{{end}}" +
		"\t\tdefault:\n" +
		"\t\t\t{\n" +
		"\t\t\t\tfieldCheck[i] = false\n" +
		"\t\t\t\tlog.Warnf(\"invalid requested field: %%s-%%s\", \"{{$.Name}}\", " +
		"x)\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\treturn q, fieldCheck, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_field_query", data)
}

// GetEdgeGetByIDStr generates the the function that retrieves edge fields.
func GetEdgeGetByIDStr(s cg.Schema, e cg.EdgeStruct) string {
	fields := e.Fields
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	fromVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"

	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.EdgeFieldStruct
		From     string
		To       string
		FromVar  string
		ToVar    string
	}{
		Name:     e.CodeName,
		Policies: pp,
		Fields:   fields,
		From:     e.FromNode.GetName(),
		To:       e.ToNode.GetName(),
		FromVar:  fromVar,
		ToVar:    toVar,
	}
	template := "// Get{{.Name}}ByID retrives the fields of a specific " +
		"{{.Name}}.\n" +
		"// If there is insufficient authorization, the field will return null.\n" +
		"func Get{{.Name}}ByID(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\tfields []string,\n" +
		") ([]interface{}, error) {\n" +
		"\n" +
		"\t// Find the {{.FromVar}} and {{.ToVar}}\n" +
		"\trow, err := models.{{.From}}Query().\n" +
		"\t\tReturnID().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tWhereID(p.Equals(id)).\n" +
		"\t\tQuery{{.To}}().\n" +
		"\t\tReturnID().\n" +
		"\t\tGenOne(conn)\n" +
		"\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif row == nil || row[0] == nil || row[1] == nil {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n" +
		"\t{{.FromVar}} := row[0].(string)\n" +
		"\t{{.ToVar}} := row[1].(string)\n" +
		"\n" +
		"\t// Create the query\n" +
		"\tq := models.{{.From}}Query().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(\n\t\tconn,\n\t\tvc, " +
		"\n\t\tparams,\n\t\tid,\n\t\t{{.FromVar}},\n\t\t{{.ToVar}},\n\t\tfields, " +
		"\n\t\tq,\n)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\tif len(q.Return) == 0 {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n" +
		"\n" +
		"\t// Execute the query\n" +
		"\trow = nil\n" +
		"\tnewConn := conn\n" +
		"\tfor i := 0; row == nil && i < constants.LogicRetryCount; i++ {\n" +
		"\t\trow, err = q.GenOne(newConn)\n" +
		"\t\tif _, isBoltErr := err.(*e.Error); isBoltErr {\n" +
		"\t\t\t// Try a new connection\n" +
		"\t\t\t(*newConn).Close()\n" +
		"\t\t\ttime.Sleep(time.Millisecond * constants.LogicRetryWait)\n" +
		"\t\t\tc, _ := db.GetDriverConn()\n" +
		"\t\t\tnewConn = c\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\t*conn = *newConn\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\n" +
		"\t// Check for the authed fields\n" +
		"\tresults := util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\n" +
		"\treturn results, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_by_id", data)
}

// GetEdgeGetByIDBatcherStr creates the batcher function for GetEdgeByID
func GetEdgeGetByIDBatcherStr(s cg.Schema, e cg.EdgeStruct) string {
	fields := e.Fields
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	fromVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"

	data := struct {
		Name     string
		Policies map[string]policies.PrivacyPolicy
		Fields   []cg.EdgeFieldStruct
		From     string
		To       string
		FromVar  string
		ToVar    string
	}{
		Name:     e.CodeName,
		Policies: pp,
		Fields:   fields,
		From:     e.FromNode.GetName(),
		To:       e.ToNode.GetName(),
		FromVar:  fromVar,
		ToVar:    toVar,
	}
	template := "// Get{{.Name}}ByIDBatcher wraps the Get{{.Name}}ByID to be " +
		"batched later.\n" +
		"func Get{{.Name}}ByIDBatcher(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\tid string,\n" +
		"\tfields []string,\n" +
		") (*util.LogicGetWrapper, error) {\n" +
		"\n" +
		"\t// Find the {{.FromVar}} and {{.ToVar}}\n" +
		"\trow, err := models.{{.From}}Query().\n" +
		"\t\tReturnID().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tWhereID(p.Equals(id)).\n" +
		"\t\tQuery{{.To}}().\n" +
		"\t\tReturnID().\n" +
		"\t\tGenOne(conn)\n" +
		"\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif row == nil || row[0] == nil || row[1] == nil {\n" +
		"\t\treturn nil, errors.New(\"no such edges\")\n" +
		"\t}\n" +
		"\t{{.FromVar}} := row[0].(string)\n" +
		"\t{{.ToVar}} := row[1].(string)\n" +
		"\n" +
		"\t// Create the query\n" +
		"\tq := models.{{.From}}Query().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(\n\t\tconn,\n\t\tvc, " +
		"\n\t\tparams,\n\t\tid,\n\t\t{{.FromVar}},\n\t\t{{.ToVar}},\n\t\tfields, " +
		"\n\t\tq,\n)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\tbatcher := new(util.LogicGetWrapper)\n" +
		"\tbatcher.Query = &q.Query\n" +
		"\tbatcher.EvalAuth = func(row []interface{}) []interface{} {\n" +
		"\t\treturn util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\t}\n" +
		"\treturn batcher, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_by_id_batcher", data)
}

// GetEdgeGetByIDsStr generates the function that gets fields on an edge.
func GetEdgeGetByIDsStr(s cg.Schema, e cg.EdgeStruct) string {

	fields := e.Fields
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	fromIDVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toIDVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"
	fromNode := e.FromNode.GetName()
	toNode := e.ToNode.GetName()

	data := struct {
		Name      string
		Policies  map[string]policies.PrivacyPolicy
		Fields    []cg.EdgeFieldStruct
		FromIDVar string
		ToIDVar   string
		FromNode  string
		ToNode    string
	}{
		Name:      e.CodeName,
		Policies:  pp,
		Fields:    fields,
		FromIDVar: fromIDVar,
		ToIDVar:   toIDVar,
		FromNode:  fromNode,
		ToNode:    toNode,
	}
	template := "// Get{{.Name}}ByIDs retrives the fields of a specific " +
		"{{.Name}}.\n" +
		"// If there is insufficient authorization, the field will return null.\n" +
		"func Get{{.Name}}ByIDs(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\t{{.FromIDVar}} string,\n" +
		"\t{{.ToIDVar}} string,\n" +
		"\tfields []string,\n" +
		") ([]interface{}, error) {\n" +
		"\n" +
		"\t// Find the ID\n" +
		"\trow, err := models.{{.FromNode}}Query().\n" +
		"\t\tWhereID(p.Equals({{.FromIDVar}})).\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tReturnID().\n" +
		"\t\tQuery{{.ToNode}}().\n" +
		"\t\tWhereID(p.Equals({{.ToIDVar}})).\n" +
		"\t\tGenOne(conn)\n" +
		"\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif row == nil || row[0] == nil {\n" +
		"\t\treturn nil, errors.New(\"no such edge\")\n" +
		"\t}\n" +
		"\tid := row[0].(string)\n" +
		"\n" +
		"\t// Create the query\n" +
		"\tq := models.{{.FromNode}}Query().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(\n\t\tconn,\n\t\tvc, " +
		"\n\t\tparams,\n\t\tid,\n\t\t{{$.FromIDVar}},\n\t\t{{$.ToIDVar}}, " +
		"\n\t\tfields,\n\t\tq,\n)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\tif len(q.Return) == 0 {\n" +
		"\t\treturn nil, nil\n" +
		"\t}\n" +
		"\n" +
		"\t// Execute the query\n" +
		"\trow = nil\n" +
		"\tnewConn := conn\n" +
		"\tfor i := 0; row == nil && i < constants.LogicRetryCount; i++ {\n" +
		"\t\trow, err = q.GenOne(newConn)\n" +
		"\t\tif _, isBoltErr := err.(*e.Error); isBoltErr {\n" +
		"\t\t\t// Try a new connection\n" +
		"\t\t\t(*newConn).Close()\n" +
		"\t\t\ttime.Sleep(time.Millisecond * constants.LogicRetryWait)\n" +
		"\t\t\tc, _ := db.GetDriverConn()\n" +
		"\t\t\tnewConn = c\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\t*conn = *newConn\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\n" +
		"\t// Check for the authed fields\n" +
		"\tresults := util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\n" +
		"\treturn results, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_by_ids", data)
}

// GetEdgeGetByIDsBatcherStr creates the batcher function for GetEdgeByIDs
func GetEdgeGetByIDsBatcherStr(s cg.Schema, e cg.EdgeStruct) string {
	fields := e.Fields
	pp := map[string]policies.PrivacyPolicy{}
	for _, x := range fields {
		pp[x.Privacy.GetName()] = x.Privacy
	}
	fromIDVar := strings.ToLower(string(e.FromNode.GetName()[0])) + "id"
	toIDVar := strings.ToLower(string(e.ToNode.GetName()[0])) + "id"
	fromNode := e.FromNode.GetName()
	toNode := e.ToNode.GetName()

	data := struct {
		Name      string
		Policies  map[string]policies.PrivacyPolicy
		Fields    []cg.EdgeFieldStruct
		FromIDVar string
		ToIDVar   string
		FromNode  string
		ToNode    string
	}{
		Name:      e.CodeName,
		Policies:  pp,
		Fields:    fields,
		FromIDVar: fromIDVar,
		ToIDVar:   toIDVar,
		FromNode:  fromNode,
		ToNode:    toNode,
	}
	template := "// Get{{.Name}}ByIDsBatcher wraps the Get{{.Name}}ByIDs to be " +
		"batched later.\n" +
		"// If there is insufficient authorization, the field will return null.\n" +
		"func Get{{.Name}}ByIDsBatcher(\n" +
		"\tconn *bolt.Conn,\n" +
		"\tvc contexts.ViewerContext,\n" +
		"\tparams context.Context,\n" +
		"\t{{.FromIDVar}} string,\n" +
		"\t{{.ToIDVar}} string,\n" +
		"\tfields []string,\n" +
		") (*util.LogicGetWrapper, error) {\n" +
		"\n" +
		"\t// Find the ID\n" +
		"\trow, err := models.{{.FromNode}}Query().\n" +
		"\t\tWhereID(p.Equals({{.FromIDVar}})).\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tReturnID().\n" +
		"\t\tQuery{{.ToNode}}().\n" +
		"\t\tWhereID(p.Equals({{.ToIDVar}})).\n" +
		"\t\tGenOne(conn)\n" +
		"\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tif row == nil || row[0] == nil {\n" +
		"\t\treturn nil, errors.New(\"no such edge\")\n" +
		"\t}\n" +
		"\tid := row[0].(string)\n" +
		"\n" +
		"\t// Create the query\n" +
		"\tq := models.{{.FromNode}}Query().\n" +
		"\t\tQuery{{.Name}}().\n" +
		"\t\tWhereID(p.Equals(id))\n" +
		"\tq, fieldCheck, err := create{{.Name}}FieldQuery(\n\t\tconn,\n\t\tvc, " +
		"\n\t\tparams,\n\t\tid,\n\t\t{{$.FromIDVar}},\n\t\t{{$.ToIDVar}}, " +
		"\n\t\tfields,\n\t\tq,\n)\n" +
		"\tif err != nil {\n" +
		"\t return nil, err\n" +
		"\t}\n" +
		"\n" +
		"\tbatcher := new(util.LogicGetWrapper)\n" +
		"\tbatcher.Query = &q.Query\n" +
		"\tbatcher.EvalAuth = func(row []interface{}) []interface{} {\n" +
		"\t\treturn util.RemoveUnauthedFields(row, fieldCheck)\n" +
		"\t}\n" +
		"\treturn batcher, nil\n" +
		"}\n"
	return cg.ExecTemplate(template, "edge_by_ids_batcher", data)
}
