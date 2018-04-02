// This program generates URN helper functions.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const sourceFileTemplateText = `package %s // auto-generated using make - DO NOT EDIT!

import (
	"strings"

	"github.com/voicera/gooseberry/urn"
)
{{range $funcName, $namespaceID := .}}
// New{{$funcName}}URN creates a new URN with the "{{$namespaceID}}"
// namespace ID.
func New{{$funcName}}URN(namespaceSpecificString string) *urn.URN {
	return urn.NewURN("{{$namespaceID}}", namespaceSpecificString)
}

// Is{{$funcName}}URN determines whether the specified URN uses
// "{{$namespaceID}}" as its namespace ID.
func Is{{$funcName}}URN(u *urn.URN) bool {
	return strings.EqualFold(u.GetNamespaceID(), "{{$namespaceID}}")
}

// Is{{$funcName}}URNWithValue determines whether the specified URN uses
// "{{$namespaceID}}" as its namespace ID and the specified
// namespaceSpecificString as its namespace-specific string.
func Is{{$funcName}}URNWithValue(u *urn.URN, namespaceSpecificString string) bool {
	return Is{{$funcName}}URN(u) && strings.EqualFold(u.GetNamespaceSpecificString(), namespaceSpecificString)
}
{{end}}`

var (
	sourceFileTemplate *template.Template
	// We autogenerate URN functions from the map below; add a new namespace ID
	// here and map it to a function name suffix.
	funcNameSuffixToNamespaceID = map[string]string{}
)

func init() {
	packageName := flag.String("p", "urns", "name of auto-generated package.")
	mapping := flag.String("m", "", "func names suffixes to namespace IDs; e.g., 'User=user Email=email'")
	flag.Parse()
	if *mapping == "" {
		flag.Usage()
		fmt.Println(`Example: go run scripts/urns/main.go -p urns -m "User=user Email=email"`)
		os.Exit(1)
	}
	sourceFileTemplate = template.Must(template.New("").Parse(fmt.Sprintf(sourceFileTemplateText, *packageName)))
	for _, m := range strings.Split(*mapping, " ") {
		indexOfMapper := strings.IndexRune(m, '=')
		funcNameSuffixToNamespaceID[m[:indexOfMapper]] = m[indexOfMapper+1:]
	}
}

func main() {
	if err := sourceFileTemplate.Execute(os.Stdout, funcNameSuffixToNamespaceID); err != nil {
		panic(err)
	}
}
