// Validating that the generated files were untampered. If there are changes,
// then the code generator will fail so manual changes are not overwritten.

package logic

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	cg "splits-go-schema-codegen/codegen"
)

// ValidateLogicSchemas validates the correctness of a generated file in the
// logic package.
func ValidateLogicSchemas(
	schemas []cg.Schema,
	packageName string,
	mergeFlag bool,
	forceFlag bool,
) (map[string][]string, error) {
	filesRead := map[string]bool{}
	manualParts := map[string][]string{}

	// Validate the signatures of all _node and _edge files
	destination := os.Args[1] + "/" + packageName + "/"
	files, _ := ioutil.ReadDir(destination)
	for _, f := range files {
		if !f.IsDir() {

			if strings.HasSuffix(f.Name(), ".go") {

				if _, ok := filesRead[f.Name()]; !ok {
					filesRead[f.Name()] = true

					// Validate the signature of the file
					filePath := destination + f.Name()
					content, err := ioutil.ReadFile(filePath)
					if err != nil {
						return manualParts, errors.New("Cannot read file: " + filePath)
					}
					manualPiece := cg.ExtractManualSections(string(content))
					manualParts[filePath] = manualPiece

					// Remove manual components
					content = []byte(cg.ReplaceAllStringSubmatchFunc(
						cg.ManualExtractor,
						string(content),
						func(groups []string) string {
							return cg.StartManual + groups[2] + cg.EndManual
						},
					))

					index := strings.Index(string(content), "\n") + 1
					firstLine := string(content[:index])
					content = content[index:]
					extractSignature := regexp.MustCompile(`\([a-zA-Z0-9]+\)`)
					signature := extractSignature.FindString(firstLine)
					if signature == "" {
						log.Printf("No file signature found, so overwriting: %s\n",
							filePath)
					} else {
						signature = signature[1 : len(signature)-1]
					}
					sum := md5.Sum([]byte(content))
					expectedSignature := hex.EncodeToString([]byte(sum[:]))

					if signature != "" && signature != expectedSignature && !mergeFlag &&
						!forceFlag {
						return manualParts, errors.New("Invalid file signature in " + filePath +
							"\nExpected '" + expectedSignature + "' and got '" + signature +
							"'")
					}
				}
			}
		}
	}
	// Add edge pointers
	for _, s := range schemas {
		for _, e := range s.GetEdges() {
			e.ToNode.AddEdgePointer(e)
		}
	}
	return manualParts, nil
}
