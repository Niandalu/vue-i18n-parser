package tree

import (
	"fmt"
	"log"
)

func Flatten(tree map[string]interface{}) map[string]string {
	result := map[string]string{}
	todo := NewStack()

	for k, v := range tree {
		todo.Push(StackItem{k, v})
	}

	for {
		node, err := todo.Pop()

		if err != nil {
			break
		}

		switch t := node.v.(type) {
		case string:
			result[node.k] = t

		case map[string]interface{}:
			for nestedK, nestedV := range t {
				todo.Push(StackItem{
					node.k + "." + nestedK,
					nestedV,
				})
			}
		case map[string]string:
			for nestedK, nestedV := range t {
				todo.Push(StackItem{
					node.k + "." + nestedK,
					nestedV,
				})
			}
		case map[interface{}]interface{}:
			for nestedK, nestedV := range t {

				todo.Push(StackItem{
					node.k + "." + fmt.Sprintf("%v", nestedK),
					nestedV,
				})
			}
		default:
			log.Printf("Unrecognized %v", t)
		}
	}

	return result
}
