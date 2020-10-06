package gqldeduplicator

import (
	"encoding/json"
	"fmt"
)

const isDeflatedKey = "__is_deflated_key__"

type (
	// DeflateResult represent deflated object result
	DeflateResult struct {
		Data       []byte
		IsDeflated bool
	}
)

// Deflate deflate similar object in graphql response by id as default identifier.
// Use deep first search (DFS) algorithm to walk over nodes and memoize object.
// If object appeared or memoized before, then it will deflated.
func Deflate(data []byte) (*DeflateResult, error) {
	var node interface{}
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	memoize := make(map[string]bool)
	resultByte, err := json.Marshal(deflate(node, memoize, "id", ""))
	if err != nil {
		return nil, err
	}

	return &DeflateResult{
		Data:       resultByte,
		IsDeflated: memoize[isDeflatedKey],
	}, nil
}

// DeflateWithCustomIdentifier deflate similar object in graphql response by identifier.
// Use deep first search (DFS) algorithm to walk over nodes and memoize object.
// If object appeared or memoized before, then it will deflated.
func DeflateWithCustomIdentifier(data []byte, identifier string) (*DeflateResult, error) {
	var node interface{}
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	memoize := make(map[string]bool)
	resultByte, err := json.Marshal(deflate(node, memoize, identifier, ""))
	if err != nil {
		return nil, err
	}

	return &DeflateResult{
		Data:       resultByte,
		IsDeflated: memoize[isDeflatedKey],
	}, nil
}

func deflate(node interface{}, memoize map[string]bool, identifier, path string) interface{} {
	switch value := node.(type) {
	case []interface{}:
		for i, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[i] = deflate(v, memoize, identifier, path)
			default:
				value[i] = v
			}
		}
		return value
	case map[string]interface{}:
		if value != nil && value[identifier] != nil && value["__typename"] != nil {
			key := fmt.Sprintf("%s,%v,%v", path, value["__typename"], value[identifier])
			if memoize[key] {
				memoize[isDeflatedKey] = true
				return map[string]interface{}{
					identifier:   value[identifier],
					"__typename": value["__typename"],
				}
			}

			memoize[key] = true
		}

		for k, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[k] = deflate(v, memoize, identifier, path+","+k)
			default:
				value[k] = v
			}
		}

		return value
	}

	return node
}
