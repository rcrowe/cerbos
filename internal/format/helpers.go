package format

import (
	"github.com/goccy/go-yaml/ast"
)

func eachMapSequence(node *ast.MappingNode, field string, cb func(node *ast.MappingNode)) {
	for _, node := range node.Values {
		if node.Key.String() == field {
			sequenceNode, ok := node.Value.(*ast.SequenceNode)
			if !ok {
				continue
			}

			for _, sequenceNode := range sequenceNode.Values {
				mappingNode, ok := sequenceNode.(*ast.MappingNode)
				if !ok {
					continue
				}
				cb(mappingNode)
			}
		}
	}
}

func getMapField(node *ast.MappingNode, field string) (*ast.MappingNode, bool) {
	for _, node := range node.Values {
		if node.Key.String() == field {
			mappingNode, ok := node.Value.(*ast.MappingNode)
			return mappingNode, ok
		}
	}
	return nil, false
}

func remapFields(node *ast.MappingNode, order []string) {
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
