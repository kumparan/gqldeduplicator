package gqldeduplicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeflate(t *testing.T) {
	tests := []struct {
		Name     string
		Given    []byte
		Expected []byte
		Deflated bool
	}{
		{
			Name:     "should deflate 1st child",
			Deflated: true,
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
		},
		{
			Name:     "should deflate nth child",
			Deflated: true,
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
								"id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
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
								"id": "1"
							}
						}
					}
				]
			}`),
		},
		{
			Name:     "should not deflate single result",
			Deflated: false,
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
			Name:     "should not deflate null result",
			Deflated: false,
			Given:    []byte(`null`),
			Expected: []byte(`null`),
		},
		{
			Name:     "should not deflate object without typename and id",
			Deflated: false,
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
			Name:     "should not deflate first object",
			Deflated: true,
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
						"id": 1,
						"name": "foo"
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
						"id": 1
					}
				]
			}`),
		},
		{
			Name:     "should not deflate array of string",
			Deflated: false,
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
			Name:     "should not deflate array of number",
			Deflated: false,
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
			Name:     "should not deflate array of boolean",
			Deflated: false,
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
			Name:     "should not deflate nested array",
			Deflated: false,
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
			result, err := Deflate(test.Given)
			assert.NoError(t, err)
			assert.JSONEq(t, string(test.Expected), string(result.Data))
			assert.Equal(t, test.Deflated, result.Deflated)
		})
	}

	t.Run("should return error on invalid json", func(t *testing.T) {
		result, err := Deflate([]byte(`{`))
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeflateIdentifier(t *testing.T) {
	tests := []struct {
		Name       string
		Given      []byte
		Expected   []byte
		Deflated   bool
		Identifier string
	}{
		{
			Name:       "should deflate 1st child",
			Deflated:   true,
			Identifier: "node_id",
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
			Name:       "should deflate nth child",
			Deflated:   true,
			Identifier: "node_id",
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
								"node_id": "1",
								"field_1": "field 1",
								"field_2": "field 2"
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
			Name:       "should not deflate single result",
			Deflated:   false,
			Identifier: "node_id",
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
			Name:       "should not deflate null result",
			Deflated:   false,
			Identifier: "node_id",
			Given:      []byte(`null`),
			Expected:   []byte(`null`),
		},
		{
			Name:       "should not deflate object without typename and node_id",
			Deflated:   false,
			Identifier: "node_id",
			Given: []byte(`
			{
				"node_id": "1",
				"foo": "bar"
			}`),
			Expected: []byte(`
			{
				"node_id": "1",
				"foo": "bar"
			}`),
		},
		{
			Name:       "should not deflate first object",
			Deflated:   true,
			Identifier: "node_id",
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
						"node_id": 1,
						"name": "foo"
					}
				]
			}`),
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
						"node_id": 1
					}
				]
			}`),
		},
		{
			Name:       "should not deflate array of string",
			Deflated:   false,
			Identifier: "node_id",
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
		},
		{
			Name:       "should not deflate array of number",
			Deflated:   false,
			Identifier: "node_id",
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
		},
		{
			Name:       "should not deflate array of boolean",
			Deflated:   false,
			Identifier: "node_id",
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
		},
		{
			Name:       "should not deflate nested array",
			Deflated:   false,
			Identifier: "node_id",
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
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := DeflateWithCustomIdentifier(test.Given, test.Identifier)
			assert.NoError(t, err)
			assert.JSONEq(t, string(test.Expected), string(result.Data))
			assert.Equal(t, test.Deflated, result.Deflated)
		})
	}

	t.Run("should return error on invalid json", func(t *testing.T) {
		result, err := DeflateWithCustomIdentifier([]byte(`{`), "node_id")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
