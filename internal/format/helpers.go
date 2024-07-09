package format

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
	cel_common "github.com/google/cel-go/common"
	cel_parser "github.com/google/cel-go/parser"
)

func getField(node ast.Node, field string) (ast.Node, bool) {
	switch rootNode := node.(type) {
	case *ast.MappingNode:
		for _, v := range rootNode.Values {
			if v.Key.String() == field {
				return v.Value, true
			}
		}
	case *ast.MappingValueNode:
		if rootNode.Key.String() == field {
			return rootNode.Value, true
		}
	}

	return nil, false
}

func orderFields(n ast.Node, order []string) {
	node, ok := n.(*ast.MappingNode)
	if !ok {
		return
	}

	orderedValues := make([]*ast.MappingValueNode, len(order))
	otherValues := make([]*ast.MappingValueNode, 0, len(node.Values))

	indexLookup := make(map[string]int, len(order))
	for i, key := range order {
		indexLookup[key] = i
	}

	for _, node := range node.Values {
		index, ok := indexLookup[node.Key.String()]
		if !ok {
			otherValues = append(otherValues, node)
			continue
		}
		orderedValues[index] = node
	}

	mergedValues := make([]*ast.MappingValueNode, 0, len(node.Values))
	for _, node := range orderedValues {
		if node != nil {
			mergedValues = append(mergedValues, node)
		}
	}
	node.Values = append(mergedValues, otherValues...)
}

func parseCEL(expr string) (string, error) {
	cel, err := cel_parser.NewParser()
	if err != nil {
		return "", fmt.Errorf("cel parser: %w", err)
	}

	ast, issues := cel.Parse(cel_common.NewStringSource(expr, "<expr>"))
	if len(issues.GetErrors()) > 0 {
		return "", fmt.Errorf("cel parser: %s", issues.ToDisplayString())
	}
	return cel_parser.Unparse(ast.Expr(), ast.SourceInfo())
}
