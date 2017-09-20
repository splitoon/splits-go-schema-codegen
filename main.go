// Package for generating code for a schema, this includes queries, mutations,
// deleters, as well as constraints and indices.

package main

import (
	"fmt"
	"log"
	"os"
	s "splits-go-api/db/models/schemas"
	cg "splits-go-schema-codegen/codegen"
	"strings"
)

var schemas = []cg.Schema{
	s.UserSchema,
	s.GroupSchema,
	s.TransactionSchema,
}
var packageName = "models"
var destination string

func main() {
	destination = os.Args[1]
	schemaNames := make([]string, 0, len(schemas))
	for _, x := range schemas {
		schemaNames = append(schemaNames, x.GetName())
	}
	log.Printf("Imported schemas: %s\n", strings.Join(schemaNames, ", "))

	// Validate the schemas
	err := cg.ValidateSchemas(schemas, packageName)
	if err != nil {
		log.Printf("Error in validating schemas")
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated := []string{}

	// Generate the constants
	constantsContent := cg.WriteConstants(schemas)
	constantsFilePath := destination + packageName + "/constants.go"
	constantsFile, err := os.Create(constantsFilePath)
	if err != nil {
		log.Printf("Error in writing to file: %s", constantsFilePath)
		log.Println(err)
		os.Exit(1)
	}
	defer constantsFile.Close()
	constantsFile.Truncate(0)
	_, err = fmt.Fprintf(constantsFile, constantsContent)
	if err != nil {
		log.Printf("Error in writing to file: %s", constantsFilePath)
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated = append(filesGenerated, constantsFilePath)

	// Generate the node and edge definition code
	for _, s := range schemas {
		schemaCode := cg.WriteSchemaNode(s, packageName)
		err = os.MkdirAll(destination+packageName+"/", os.ModePerm)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		filePath := destination + packageName + "/" + strings.ToLower(s.GetName()) +
			"_node.go"
		file, err2 := os.Create(filePath)
		if err != nil {
			log.Printf("Error in writing to file: %s", filePath)
			log.Println(err2)
			os.Exit(1)
		}
		defer file.Close()
		file.Truncate(0)
		_, err = fmt.Fprintf(file, schemaCode)
		if err != nil {
			log.Printf("Error in writing to file: %s", filePath)
			log.Println(err)
			os.Exit(1)
		}
		filesGenerated = append(filesGenerated, filePath)
		for _, e := range s.GetEdges() {
			edgeCode := cg.WriteSchemaEdge(s, e, packageName)
			edgeFilePath := destination + packageName + "/" + strings.ToLower(e.Name) +
				"_edge.go"
			err = os.MkdirAll(destination+packageName+"/", os.ModePerm)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			edgeFile, err2 := os.Create(edgeFilePath)
			if err != nil {
				log.Printf("Error in writing to file: %s", filePath)
				log.Println(err2)
				os.Exit(1)
			}
			defer file.Close()
			_, err = fmt.Fprintf(edgeFile, edgeCode)
			if err != nil {
				log.Printf("Error in writing to file: %s", filePath)
				log.Println(err)
				os.Exit(1)
			}
			filesGenerated = append(filesGenerated, edgeFilePath)
		}
	}

	// Generate the constraints
	constraintContent := cg.WriteConstraints(schemas)
	err = os.MkdirAll(destination+packageName+"/constraints/data", os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	constraintFilePath := destination + packageName + "/constraints/data/constraints.json"
	constraintFile, err := os.Create(constraintFilePath)
	if err != nil {
		log.Printf("Error in writing to file: %s", constraintFilePath)
		log.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(constraintFile, constraintContent)
	filesGenerated = append(filesGenerated, constraintFilePath)

	// Generate the indices
	indicesContent := cg.WriteIndices(schemas)
	err = os.MkdirAll(destination+packageName+"/indices/data", os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	indicesFilePath := destination + packageName + "/indices/data/indices.json"
	indicesFile, err := os.Create(indicesFilePath)
	if err != nil {
		log.Printf("Error in writing to file: %s", indicesFilePath)
		log.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(indicesFile, indicesContent)
	filesGenerated = append(filesGenerated, indicesFilePath)

	for _, file := range filesGenerated {
		fmt.Printf("Generated %s\n", file)
	}
}
