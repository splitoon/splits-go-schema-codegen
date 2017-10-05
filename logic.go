// // Writer for generating logic related code.
//
package main

import (
	"fmt"
	"log"
	"os"
	"splits-go-schema-codegen/codegen/logic"
	"strings"
)

func generateLogicCode(mergeFlag bool, forceFlag bool) {
	fmt.Println("\nGENERATING LOGIC...")

	destination := os.Args[1] + "/"
	packageName := "logic"

	// Validate the schemas
	manualParts, err := logic.ValidateLogicSchemas(schemas, packageName,
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

	// Generate the node and edge definition code
	for _, s := range schemas {
		f := destination + packageName + "/" + strings.ToLower(s.GetName()) + ".go"
		manualPart := manualParts[f]
		schemaCode := logic.WriteSchemaLogicNode(s, manualPart, packageName)
		err = os.MkdirAll(destination+packageName+"/", os.ModePerm)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		filePath := destination + packageName + "/" + strings.ToLower(s.GetName()) +
			".go"
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

		// Generate edge definitions
		for _, e := range s.GetEdges() {
			f := destination + packageName + "/" + strings.ToLower(e.Name) + ".go"
			manualPart := manualParts[f]
			edgeCode := logic.WriteSchemaLogicEdge(s, e, manualPart, packageName)
			edgeFilePath := destination + packageName + "/" + strings.ToLower(e.Name) +
				".go"
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

	for _, file := range filesGenerated {
		fmt.Printf("Generated %s\n", file)
	}
}
