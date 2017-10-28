// // Writer for generating logic related code.
//
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	cg "splits-go-schema-codegen/codegen"
	"splits-go-schema-codegen/codegen/graphql"
	"strings"
)

func prepGraphQLSchema() (cg.GraphQLSchema, error) {
	schema := cg.GraphQLSchema{}
	nodes := []*cg.GraphQLNode{}
	oppositeNodes := map[string]*cg.GraphQLNode{}
	for _, s := range schemas {
		n := s.GetGraphQLNode()
		if n != nil {
			nodes = append(nodes, n)
			oppositeNodes[n.Name] = n
		}
	}
	for _, n := range nodes {
		schema.Nodes = append(schema.Nodes, *n)
		for _, e := range n.Edges {
			schema.Edges = append(schema.Edges, e)
			if e.IncludeReverse {
				edge := cg.GraphQLEdge{
					From:             e.To,
					To:               e.From,
					FieldName:        e.ReverseFieldName,
					FieldCodeName:    e.ReverseFieldCodeName,
					FieldResolveName: e.ReverseFieldResolveName,
					Description:      e.ReverseDescription,
					FromCodeName:     e.ToCodeName,
					ToCodeName:       e.FromCodeName,
					Fields:           e.Fields,
					TotalName:        e.TotalName,
					IsReverse:        true,
					EdgeCodeName:     e.EdgeCodeName,
				}
				oppositeNode, ok := oppositeNodes[e.To]
				if !ok {
					return schema, errors.New("graphql edge has no opposite node")
				}
				oppositeNode.Edges = append(oppositeNode.Edges, edge)
				schema.Edges = append(schema.Edges, edge)
			}
		}
	}
	return schema, nil
}

func generateGraphQLCode(mergeFlag bool, forceFlag bool) {
	fmt.Println("\nGENERATING GRAPHQL...")

	graphqlSchema, err := prepGraphQLSchema()
	if err != nil {
		fmt.Println("Could not prepare graphql schema")
		return
	}
	destination := os.Args[1] + "/"
	packageName := "graphql_auto"

	// Validate the schemas
	manualParts, err := graphql.ValidateGraphQLSchemas(schemas, packageName,
		mergeFlag, forceFlag)
	if err != nil {
		log.Printf("Error in validating schemas")
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated := []string{}

	if forceFlag {
		manualParts = map[string][]string{}
	}

	// ===========================================================================
	// Generate the schema code
	// ===========================================================================
	f := destination + "api/" + packageName + "/schema.go"
	manualPart := manualParts[f]
	schemaCode := graphql.WriteGraphQLSchema(schemas, graphqlSchema, manualPart,
		packageName)
	err = os.MkdirAll(destination+packageName+"/", os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	file, err := os.Create(f)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	file.Truncate(0)
	_, err = fmt.Fprintf(file, schemaCode)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated = append(filesGenerated, f)

	// ===========================================================================
	// Generate the node type
	// ===========================================================================
	nextPackageName := "resolvers"
	f = destination + "api/" + packageName + "/" + nextPackageName +
		"/type_node.go"
	manualPart = manualParts[f]
	schemaCode = graphql.WriteGraphQLNodeType(schemas, graphqlSchema, manualPart,
		nextPackageName)
	err = os.MkdirAll(destination+packageName+"/"+nextPackageName+"/",
		os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	file, err = os.Create(f)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	file.Truncate(0)
	_, err = fmt.Fprintf(file, schemaCode)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated = append(filesGenerated, f)

	// ===========================================================================
	// Generate the root query
	// ===========================================================================
	nextPackageName = "resolvers"
	f = destination + "api/" + packageName + "/" + nextPackageName +
		"/type_root_query.go"
	manualPart = manualParts[f]
	schemaCode = graphql.WriteRootQueryType(schemas, graphqlSchema, manualPart,
		nextPackageName)
	err = os.MkdirAll(destination+packageName+"/"+nextPackageName+"/",
		os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	file, err = os.Create(f)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	file.Truncate(0)
	_, err = fmt.Fprintf(file, schemaCode)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated = append(filesGenerated, f)

	// ===========================================================================
	// Generate the node resolvers
	// ===========================================================================
	nextPackageName = "resolvers"
	for _, n := range graphqlSchema.Nodes {
		f = destination + "api/" + packageName + "/" + nextPackageName +
			"/type_" + strings.ToLower(n.CodeName) + ".go"
		manualPart = manualParts[f]
		schemaCode = graphql.WriteGQLNodeResolverType(n, manualPart,
			nextPackageName)
		err = os.MkdirAll(destination+packageName+"/"+nextPackageName+"/",
			os.ModePerm)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		file, err = os.Create(f)
		if err != nil {
			log.Printf("Error in writing to file: %s", f)
			log.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		file.Truncate(0)
		_, err = fmt.Fprintf(file, schemaCode)
		if err != nil {
			log.Printf("Error in writing to file: %s", f)
			log.Println(err)
			os.Exit(1)
		}
		filesGenerated = append(filesGenerated, f)
	}

	// ===========================================================================
	// Generate the edge resolvers
	// ===========================================================================
	nextPackageName = "resolvers"
	for _, n := range graphqlSchema.Nodes {
		for _, e := range n.Edges {
			f = destination + "api/" + packageName + "/" + nextPackageName +
				"/type_edge_" + strings.ToLower(e.FromCodeName+"to"+e.ToCodeName) +
				".go"
			manualPart = manualParts[f]
			schemaCode = graphql.WriteGQLEdgeResolverType(e, manualPart,
				nextPackageName)
			err = os.MkdirAll(destination+packageName+"/"+nextPackageName+"/",
				os.ModePerm)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			file, err = os.Create(f)
			if err != nil {
				log.Printf("Error in writing to file: %s", f)
				log.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			file.Truncate(0)
			_, err = fmt.Fprintf(file, schemaCode)
			if err != nil {
				log.Printf("Error in writing to file: %s", f)
				log.Println(err)
				os.Exit(1)
			}
			filesGenerated = append(filesGenerated, f)
		}
	}

	// ===========================================================================
	// Generate the dataloader
	// ===========================================================================
	nextPackageName = "resolvers"
	f = destination + "api/" + packageName + "/" + nextPackageName +
		"/dataloader_batcher.go"
	manualPart = manualParts[f]
	schemaCode = graphql.WriteDataloaderBatcher(schemas, graphqlSchema,
		manualPart, nextPackageName)
	err = os.MkdirAll(destination+packageName+"/"+nextPackageName+"/",
		os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	file, err = os.Create(f)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	file.Truncate(0)
	_, err = fmt.Fprintf(file, schemaCode)
	if err != nil {
		log.Printf("Error in writing to file: %s", f)
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated = append(filesGenerated, f)

	for _, file := range filesGenerated {
		fmt.Printf("Generated %s\n", file)
	}
}
