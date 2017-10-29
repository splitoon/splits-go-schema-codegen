// Validating that the generated files were untampered. If there are changes,
// then the code generator will fail so manual changes are not overwritten.

package graphql

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

// ValidateGraphQLSchemas validates the correctness of a generated file in the
// graphql package.
func ValidateGraphQLSchemas(
	schemas []cg.Schema,
	packageName string,
	mergeFlag bool,
	forceFlag bool,
) (map[string][]string, error) {
	filesRead := map[string]bool{}
	manualParts := map[string][]string{}

	// Validate the signatures of all _type files
	destination := os.Args[1] + "/api/" + packageName + "/resolvers/"

	files, _ := ioutil.ReadDir(destination)
	for _, f := range files {
		if !f.IsDir() {

			if strings.HasPrefix(f.Name(), "type_") {

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
					content = []byte(strings.Replace(string(content), "%s", "%%s", -1))
					sum := md5.Sum([]byte(content))
					expectedSignature := hex.EncodeToString([]byte(sum[:]))

					if signature != "" && signature != expectedSignature && !mergeFlag &&
						!forceFlag {
						return manualParts, errors.New("Invalid file signature in " +
							filePath + "\nExpected '" + expectedSignature + "' and got '" +
							signature + "'")
					}
				}
			}
		}
	}

	// Read the specific files
	destinations := []string{
		os.Args[1] + "/api/" + packageName + "/resolvers/dataloader_batcher.go",
		os.Args[1] + "/api/" + packageName + "/schema.go",
	}
	for _, destination = range destinations {
		if _, ok := filesRead[destination]; !ok {
			filesRead[destination] = true

			// Validate the signature of the file
			filePath := destination
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
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
			content = []byte(strings.Replace(string(content), "%s", "%%s", -1))
			sum := md5.Sum([]byte(content))
			expectedSignature := hex.EncodeToString([]byte(sum[:]))

			if signature != "" && signature != expectedSignature && !mergeFlag &&
				!forceFlag {
				return manualParts, errors.New("Invalid file signature in " +
					filePath + "\nExpected '" + expectedSignature + "' and got '" +
					signature + "'")
			}
		}
	}
	return manualParts, nil
}
