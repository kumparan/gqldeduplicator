package resolver

// Resolver represent graphql resolver
type Resolver struct{}

// Parents resolve parents query
func (r Resolver) Parents() []Parent {
	child := &Child{
		id:   "1",
		name: "child 1",
		age:  10,
	}
	return []Parent{
		{id: "1", name: "parent 1", child: child},
		{id: "2", name: "parent 2", child: child},
		{id: "3", name: "parent 3", child: child},
		{id: "4", name: "parent 4", child: nil},
		{id: "5", name: "parent 5", child: &Child{id: "2", name: "child 2", age: 9}},
	}
}
