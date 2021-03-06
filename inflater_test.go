package gqldeduplicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInflate(t *testing.T) {
	tests := []struct {
		Name     string
		Expected []byte
		Given    []byte
		Inflated bool
	}{
		{
			Name:     "should inflate 1st child",
			Inflated: true,
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					},
					{
						"__typename": "Parent",
						"id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"id": "1"
						},
						"another_child": {
							"__typename": "Child",
							"id": "2",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					}
				]
			}`),
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					},
					{
						"__typename": "Parent",
						"id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"id": "2",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					}
				]
			}`),
		},
		{
			Name:     "should inflate nth child",
			Inflated: true,
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					},
					{
						"__typename": "Parent",
						"id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"id": "2",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"id": "1"
							}
						}
					}
				]
			}`),
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"id": "1",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					},
					{
						"__typename": "Parent",
						"id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"id": "2",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					}
				]
			}`),
		},
		{
			Name:     "should not inflate single result",
			Inflated: false,
			Given: []byte(`
			{
				"root": true
			}`),
			Expected: []byte(`
			{
				"root": true
			}`),
		},
		{
			Name:     "should not inflate null result",
			Inflated: false,
			Given:    []byte(`null`),
			Expected: []byte(`null`),
		},
		{
			Name:     "should not inflate object without typename and id",
			Inflated: false,
			Given: []byte(`
			{
				"id": "1",
				"foo": "bar"
			}`),
			Expected: []byte(`
			{
				"id": "1",
				"foo": "bar"
			}`),
		},
		{
			Name:     "should not inflate first object",
			Inflated: true,
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "foo",
						"id": 1,
						"name": "foo"
					},
					{
						"__typename": "foo",
						"id": 1
					}
				]
			}`),
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "foo",
						"id": 1,
						"name": "foo"
					},
					{
						"__typename": "foo",
						"id": 1,
						"name": "foo"
					}
				]
			}`),
		},
		{
			Name:     "should not inflate array of string",
			Inflated: false,
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"names": [
						"foo",
						"bar"
					]
				}
			}`),
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"names": [
						"foo",
						"bar"
					]
				}
			}`),
		},
		{
			Name:     "should not inflate array of number",
			Inflated: false,
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"numbers": [
						1,
						2
					]
				}
			}`),
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"numbers": [
						1,
						2
					]
				}
			}`),
		},
		{
			Name:     "should not inflate array of boolean",
			Inflated: false,
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"values": [
						true,
						false
					]
				}
			}`),
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"values": [
						true,
						false
					]
				}
			}`),
		},
		{
			Name:     "should not inflate nested array",
			Inflated: false,
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"values": [
						[
							true,
							false
						],
						[
							true,
							false
						]
					]
				}
			}`),
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"id": 1,
					"values": [
						[
							true,
							false
						],
						[
							true,
							false
						]
					]
				}
			}`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := Inflate(test.Given)
			assert.NoError(t, err)
			assert.JSONEq(t, string(test.Expected), string(result.Data))
			assert.Equal(t, test.Inflated, result.Inflated)
		})
	}

	t.Run("should return error on invalid json", func(t *testing.T) {
		result, err := Inflate([]byte(`{`))
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInflateIdentifier(t *testing.T) {
	tests := []struct {
		Name       string
		Expected   []byte
		Given      []byte
		Inflated   bool
		Identifier string
	}{
		{
			Name:       "should inflate 1st child",
			Inflated:   true,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"node_id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					},
					{
						"__typename": "Parent",
						"node_id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"node_id": "2",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					}
				]
			}`),
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"node_id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1"
						},
						"another_child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					},
					{
						"__typename": "Parent",
						"node_id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"node_id": "1"
						},
						"another_child": {
							"__typename": "Child",
							"node_id": "2",
							"field_1": "field 1",
							"field_2": "field 2"
						}
					}
				]
			}`),
		},
		{
			Name:       "should inflate nth child",
			Inflated:   true,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"node_id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"node_id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					},
					{
						"__typename": "Parent",
						"node_id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"node_id": "2",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"node_id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					}
				]
			}`),
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "Parent",
						"node_id": "1",
						"name": "parent 1",
						"child": {
							"__typename": "Child",
							"node_id": "1",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"node_id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
							}
						}
					},
					{
						"__typename": "Parent",
						"node_id": "2",
						"name": "parent 2",
						"child": {
							"__typename": "Child",
							"node_id": "2",
							"field_1": "field 1",
							"another_child": {
								"__typename": "AnotherChild",
								"node_id": "1"
							}
						}
					}
				]
			}`),
		},
		{
			Name:       "should not inflate single result",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": true
			}`),
			Given: []byte(`
			{
				"root": true
			}`),
		},
		{
			Name:       "should not inflate null result",
			Inflated:   false,
			Identifier: "node_id",
			Expected:   []byte(`null`),
			Given:      []byte(`null`),
		},
		{
			Name:       "should not inflate object without typename and node_id",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"node_id": "1",
				"foo": "bar"
			}`),
			Given: []byte(`
			{
				"node_id": "1",
				"foo": "bar"
			}`),
		},
		{
			Name:       "should not inflate first object",
			Inflated:   true,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": [
					{
						"__typename": "foo",
						"node_id": 1,
						"name": "foo"
					},
					{
						"__typename": "foo",
						"node_id": 1,
						"name": "foo"
					}
				]
			}`),
			Given: []byte(`
			{
				"root": [
					{
						"__typename": "foo",
						"node_id": 1,
						"name": "foo"
					},
					{
						"__typename": "foo",
						"node_id": 1
					}
				]
			}`),
		},
		{
			Name:       "should not inflate array of string",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"names": [
						"foo",
						"bar"
					]
				}
			}`),
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"names": [
						"foo",
						"bar"
					]
				}
			}`),
		},
		{
			Name:       "should not inflate array of number",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"numbers": [
						1,
						2
					]
				}
			}`),
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"numbers": [
						1,
						2
					]
				}
			}`),
		},
		{
			Name:       "should not inflate array of boolean",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"values": [
						true,
						false
					]
				}
			}`),
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"values": [
						true,
						false
					]
				}
			}`),
		},
		{
			Name:       "should not inflate nested array",
			Inflated:   false,
			Identifier: "node_id",
			Expected: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"values": [
						[
							true,
							false
						],
						[
							true,
							false
						]
					]
				}
			}`),
			Given: []byte(`
			{
				"root": {
					"__typename": "foo",
					"node_id": 1,
					"values": [
						[
							true,
							false
						],
						[
							true,
							false
						]
					]
				}
			}`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := InflateWithCustomIdentifier(test.Given, test.Identifier)
			assert.NoError(t, err)
			assert.JSONEq(t, string(test.Expected), string(result.Data))
			assert.Equal(t, test.Inflated, result.Inflated)
		})
	}

	t.Run("should return error on invalid json", func(t *testing.T) {
		result, err := InflateWithCustomIdentifier([]byte(`{`), "node_id")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
