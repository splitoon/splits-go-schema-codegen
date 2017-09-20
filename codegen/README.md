# splits-go-schema-codegen

This package is for generating some basic code for the splits db layer. It generates query, mutator, and deleter code, as well as constraints and indices for some schema.

## Usage
1. Write a schema defined in splits-go-api/db/models/schemas.
2. Add that schema in the `schemas` var in splits-go-schema-codegen/main.go.
3. Execute `./scripts/go-run.sh` to read in the schemas and generate code.
