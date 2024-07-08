package format

import "github.com/goccy/go-yaml/ast"

var (
	derivedRolesDocumentOrder = []string{
		"apiVersion",
		"description",
		"derivedRoles",
	}

	// `derivedRoles` order
	derivedRolesFieldOrder = []string{
		"name",
		"variables",
		"definitions",
	}

	// `derivedRoles.definitions` order
	derivedRolesDefinitionsOrder = []string{
		"name",
		"parentRoles",
		"condition",
	}
)

func formatDerivedRoles() formattter {
	return func(doc *ast.DocumentNode) error {
		document, ok := doc.Body.(*ast.MappingNode)
		if !ok {
			return nil
		}
		derivedRoles, ok := getMapField(document, "derivedRoles")
		if !ok {
			// not a `derivedRoles` policy
			return nil
		}

		// order document
		remapFields(document, derivedRolesDocumentOrder)

		// order `derivedRoles`
		remapFields(derivedRoles, derivedRolesFieldOrder)

		// order `derivedRoles.definitions`
		eachMapSequence(derivedRoles, "definitions", func(node *ast.MappingNode) {
			remapFields(node, derivedRolesDefinitionsOrder)
		})

		return nil
	}
}
