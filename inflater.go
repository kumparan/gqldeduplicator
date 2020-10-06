package gqldeduplicator

import (
	"encoding/json"
	"fmt"
)

const isInflatedKey = "__is_inflated_key__"

type (
	// InflateResult represent inflated object result
	InflateResult struct {
		Data       []byte
		IsInflated bool
	}
)

// Inflate inflate similar object in graphql response by id as identifier.
// Use deep first search (DFS) algorithm to walk over nodes and memoize object.
// If object appeared or memoized before, then it will inflated.
func Inflate(data []byte) (*InflateResult, error) {
	var node interface{}
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	memoize := make(map[string]interface{})
	resultByte, err := json.Marshal(inflate(node, memoize, "id", ""))
	if err != nil {
		return nil, err
	}

	return &InflateResult{
		Data:       resultByte,
		IsInflated: memoize[isInflatedKey] != nil,
	}, nil
}

// InflateWithCustomIdentifier inflate similar object in graphql response by identifier.
// Use deep first search (DFS) algorithm to walk over nodes and memoize object.
// If object appeared or memoized before, then it will inflated.
func InflateWithCustomIdentifier(data []byte, identifier string) (*InflateResult, error) {
	var node interface{}
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	memoize := make(map[string]interface{})
	resultByte, err := json.Marshal(inflate(node, memoize, identifier, ""))
	if err != nil {
		return nil, err
	}

	return &InflateResult{
		Data:       resultByte,
		IsInflated: memoize[isInflatedKey] != nil,
	}, nil
}

func inflate(node interface{}, memoize map[string]interface{}, identifier, path string) interface{} {
	switch value := node.(type) {
	case []interface{}:
		for i, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[i] = inflate(v, memoize, identifier, path)
			default:
				value[i] = v
			}
		}
		return value
	case map[string]interface{}:
		if value != nil && value[identifier] != nil && value["__typename"] != nil {
			key := fmt.Sprintf("%s,%v,%v", path, value["__typename"], value[identifier])
			if memoize[key] != nil {
				memoize[isInflatedKey] = true
				return memoize[key]
			}

			memoize[key] = value
		}

		for k, v := range value {
			switch v.(type) {
			case []interface{}, map[string]interface{}:
				value[k] = inflate(v, memoize, identifier, path+","+k)
			default:
				value[k] = v
			}
		}

		return value
	}

	return node
}
