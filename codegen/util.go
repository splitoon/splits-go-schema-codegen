// Utility functions for

package codegen

import (
	"regexp"
	"strings"
)

// ExtractManualSections returns the strings that represent the manual sections
// in a file.
func ExtractManualSections(content string) []string {
	sections := []string{}
	re := regexp.MustCompile(
		"(?sU)" +
			strings.Replace(StartManual, "*", "\\*", -1) + "\\n?" +
			"(.*)" +
			"\\n\\s*" + strings.Replace(EndManual, "*", "\\*", -1),
	)
	res := re.FindAllStringSubmatch(content, -1)
	for _, v := range res {
		if len(v) == 2 {
			s := v[1]
			sections = append(sections, s)
		} else {
			sections = append(sections, "")
		}
	}
	return sections
}

// StartManual is a string constant that represents the start of a manul block.
const StartManual = "// * START MANUAL SECTION *"

// EndManual is a string consctant that represents the end of a manual block.
const EndManual = "// * END MANUAL SECTION *"

// ManualExtractor is a regex that matches the manual pieces of code
var ManualExtractor = regexp.MustCompile(
	"(?sU)" +
		strings.Replace(StartManual, "*", "\\*", -1) + "\\n?" +
		"(.*)" +
		"(\\n\\s*)" + strings.Replace(EndManual, "*", "\\*", -1),
)

// ReplaceAllStringSubmatchFunc finds submatches and replaces the submatches
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}
		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
