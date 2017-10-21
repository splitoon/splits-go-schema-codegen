// Writer for generating db related code.

package main

import (
	"fmt"
	"log"
	"os"
	"splits-go-schema-codegen/codegen/db"
	"strings"
)

func generateDBCode(mergeFlag bool, forceFlag bool) {
	fmt.Println("\nGENERATING DB...")

	destination := os.Args[1] + "/db/"
	packageName := "models"

	// Validate the schemas
	err := db.ValidateDBSchemas(schemas, packageName)
	if err != nil {
		log.Printf("Error in validating schemas")
		log.Println(err)
		os.Exit(1)
	}
	filesGenerated := []string{}

	// Generate the constants
	constantsContent := db.WriteConstants(schemas)
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
		schemaCode := db.WriteSchemaNode(s, packageName)
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
			edgeCode := db.WriteSchemaEdge(s, e, packageName)
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
	constraintContent := db.WriteConstraints(schemas)
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
	indicesContent := db.WriteIndices(schemas)
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

	// Generate the db tests
	autogenTestContent := db.WriteAutogenTests(schemas, packageName)
	err = os.MkdirAll(destination+packageName, os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	autogenFilePath := destination + packageName + "/autogen_test.go"
	autogenFile, err := os.Create(autogenFilePath)
	if err != nil {
		log.Printf("Error in writing to file: %s", autogenFilePath)
		log.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(autogenFile, autogenTestContent)
	filesGenerated = append(filesGenerated, autogenFilePath)

	// List the generated files
	for _, file := range filesGenerated {
		fmt.Printf("Generated %s\n", file)
	}
}
