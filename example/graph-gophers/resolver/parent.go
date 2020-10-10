package resolver

import "github.com/graph-gophers/graphql-go"

// Parent represent parent node
type Parent struct {
	id    graphql.ID
	name  string
	child *Child
}

// ID resolve id field
func (p Parent) ID() graphql.ID {
	return p.id
}

// Name resolve name field
func (p Parent) Name() string {
	return p.name
}

// Child resolve child node
func (p Parent) Child() *Child {
	return p.child
}
