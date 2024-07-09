package format

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

var (
	stringQuoteStyle = token.DoubleQuoteType

	// `variables` order
	variablesFieldOrder = []string{
		"import",
		"local",
	}
)

type formattter func(*ast.DocumentNode) error

func YAML(input []byte) ([]byte, error) {
	formatters := []formattter{
		formatStringQuotes(),
		formatVariables(),
		formatDerivedRoles(),
	}

	file, err := parser.ParseBytes(input, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, doc := range file.Docs {
		for _, f := range formatters {
			if err := f(doc); err != nil {
				return nil, err
			}
		}
	}

	return []byte(file.String()), nil
}

// formatStringQuotes forces a consistent string quote style.
func formatStringQuotes() formattter {
	return func(doc *ast.DocumentNode) error {
		nodes := ast.Filter(ast.StringType, doc)
		for _, node := range nodes {
			if stringNode, ok := node.(*ast.StringNode); ok {
				switch stringNode.Token.Type {
				case token.SingleQuoteType:
					stringNode.Token.Type = stringQuoteStyle
				case token.DoubleQuoteType:
					stringNode.Token.Type = stringQuoteStyle
				}
			}
		}
		return nil
	}
}

// formatVariables orders `variables` field & formats CEL.
func formatVariables() formattter {
	searchFields := []string{
		"derivedRoles",
		"resourcePolicy",
		"principalPolicy",
	}

	return func(doc *ast.DocumentNode) error {
		documentNode, ok := doc.Body.(*ast.MappingNode)
		if !ok {
			return nil
		}

		for _, search := range searchFields {
			searchNode, searchOK := getField(documentNode, search)
			if !searchOK {
				continue
			}

			variablesNode, variablesOK := getField(searchNode, "variables")
			if !variablesOK {
				continue
			}

			// order `{search}.variables`
			orderFields(variablesNode, variablesFieldOrder)

			// format `{search}.variables.local` expressions
			localVariablesNode, localVariablesOK := getField(variablesNode, "local")
			if !localVariablesOK {
				continue
			}

			switch rootNode := localVariablesNode.(type) {
			case *ast.MappingNode:
				for _, v := range rootNode.Values {
					if exprNode, ok := v.Value.(*ast.StringNode); ok {
						expr, err := parseCEL(exprNode.Value)
						if err != nil {
							return err
						}
						exprNode.Value = expr
					}
				}
			case *ast.MappingValueNode:
				if exprNode, ok := rootNode.Value.(*ast.StringNode); ok {
					expr, err := parseCEL(exprNode.Value)
					if err != nil {
						return err
					}
					exprNode.Value = expr
				}
			default:
				continue
			}
		}
		return nil
	}
}
