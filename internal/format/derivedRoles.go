package format

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
)

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
		documentNode, ok := doc.Body.(*ast.MappingNode)
		if !ok {
			return nil
		}

		derivedRolesNode, derivedRolesOK := getField(documentNode, "derivedRoles")
		if !derivedRolesOK {
			// not a `derivedRoles` policy
			return nil
		}

		// order policy
		orderFields(documentNode, derivedRolesDocumentOrder)

		// order `derivedRoles`
		orderFields(derivedRolesNode, derivedRolesFieldOrder)

		// order `derivedRoles.definitions`
		definitionsNode, definitionsOK := getField(derivedRolesNode, "definitions")
		if !definitionsOK {
			return nil
		}
		definitionsSequenceNode, definitionsSequenceOK := definitionsNode.(*ast.SequenceNode)
		if !definitionsSequenceOK {
			return nil
		}
		for _, v := range definitionsSequenceNode.Values {
			orderFields(v, derivedRolesDefinitionsOrder)

			// format `derivedRoles.definitions[*].condition` expressions
			conditionNode, conditionOK := getField(v, "condition")
			if !conditionOK {
				continue
			}

			// walk all string nodes even though we get back key fields
			// this allows us to more easily handle complex nested structures
			conditionStringNodes := ast.Filter(ast.StringType, conditionNode)
			if len(conditionStringNodes) == 0 {
				continue
			}
			for _, field := range conditionStringNodes {
				conditionFieldNode, ok := field.(*ast.StringNode)
				if !ok {
					continue
				}

				// only look to format `expr` fields
				// e.g: $.derivedRoles.definitions[0].condition.match.expr
				if !strings.HasSuffix(conditionFieldNode.GetPath(), ".expr") {
					continue
				}

				expr, err := parseCEL(conditionFieldNode.Value)
				if err != nil {
					return err
				}
				conditionFieldNode.Value = expr
			}
		}

		return nil
	}
}
