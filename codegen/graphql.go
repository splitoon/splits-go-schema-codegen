package codegen

// GraphQLNode wrapper around the exposed graphql nodes.
type GraphQLNode struct {
	Name        string
	Description string
	CodeName    string
	Fields      []GraphQLField
	Edges       []GraphQLEdge
}

// GraphQLEdge wrapper around the exposed graphql edges.
type GraphQLEdge struct {
	From                    string
	To                      string
	FieldName               string
	FieldCodeName           string
	FieldResolveName        string
	TotalName               string
	ReverseFieldName        string
	ReverseFieldCodeName    string
	ReverseFieldResolveName string
	Description             string
	ReverseDescription      string
	FromCodeName            string
	ToCodeName              string
	Fields                  []GraphQLField
	IncludeReverse          bool
	IsReverse               bool
	EdgeCodeName            string
	OrderBy                 string
	ReverseOrderBy          string
}

// GraphQLField wrapper around a graphql field.
type GraphQLField struct {
	Name        string
	Type        string
	Description string
	CodeName    string
	CodeType    string
}

// GraphQLSchema wrapper around the exposed graphql parts of the schema.
type GraphQLSchema struct {
	Nodes []GraphQLNode
	Edges []GraphQLEdge
}
