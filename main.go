// Package for generating code for a schema, this includes queries, mutations,
// deleters, as well as constraints and indices.

package main

import (
	"fmt"
	"os"
	s "splits-go-api/db/models/schemas"
	cg "splits-go-schema-codegen/codegen"
)

var schemas = []cg.Schema{
	s.UserSchema,
	s.GroupSchema,
	s.ReceiptSchema,
	s.TransactionSchema,
}

func main() {
	schemaNames := make([]string, 0, len(schemas))
	for _, x := range schemas {
		schemaNames = append(schemaNames, x.GetName())
	}

	// Check for flags. Due to golang's flags package being weird. This is using
	// a really dumb manual way.
	var mergeFlag bool
	var forceFlag bool
	for _, v := range os.Args[1:] {
		if v == "--merge" || v == "-m" {
			mergeFlag = true
		} else if v == "--force" || v == "-f" {
			forceFlag = true
		}
	}

	if mergeFlag && forceFlag {
		fmt.Println("Cannot MERGE and FORCE at the same time. Defaulting to MERGE.")
		forceFlag = false
	} else if mergeFlag {
		fmt.Println("MERGE FLAG IS SET")
	} else if forceFlag {
		fmt.Println("FORCE FLAG IS SET")
	}

	generateDBCode(mergeFlag, forceFlag)
	generateLogicCode(mergeFlag, forceFlag)
}
