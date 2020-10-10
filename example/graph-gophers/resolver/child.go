package resolver

import "github.com/graph-gophers/graphql-go"

// Child represent child node
type Child struct {
	id   graphql.ID
	name string
	age  int32
}

// ID resolve id field
func (c Child) ID() graphql.ID {
	return c.id
}

// Name resolve name field
func (c Child) Name() string {
	return c.name
}

// Age resolve age field
func (c Child) Age() int32 {
	return c.age
}
