package tree

import (
	"log"
	"reflect"
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
				keyType := reflect.TypeOf(nestedK).Kind()
				if keyType != reflect.String {
					log.Printf("%v", node)
					log.Fatalf("ERROR: Only string is accepted, but %v is %v.", nestedK, keyType)
				}

				todo.Push(StackItem{
					node.k + "." + nestedK.(string),
					nestedV,
				})
			}
		default:
			log.Printf("Unrecognized %v", t)
		}
	}

	return result
}
