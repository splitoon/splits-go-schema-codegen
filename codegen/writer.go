// Generic writing utility functions

package codegen

import (
	"bytes"
	"text/template"
)

// =============================================================================
// Template executor
// =============================================================================

// ExecTemplate executes a template to return a string.
func ExecTemplate(
	tmplStr string,
	tmplName string,
	data interface{},
) string {
	tmpl, err := template.New(tmplName).Parse(tmplStr)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}
	return tpl.String()
}
