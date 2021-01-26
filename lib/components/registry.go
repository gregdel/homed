package components

import "fmt"

var registeredComponents map[string]func() Component

// Register a component
func Register(t Type, fn func() Component) {
	if registeredComponents == nil {
		registeredComponents = map[string]func() Component{}
	}

	name := string(t)
	_, ok := registeredComponents[name]
	if ok {
		err := fmt.Errorf("components: component %s already regitered", name)
		panic(err)
	}

	registeredComponents[name] = fn
}

// newComponent returns a new component from a type name
func newComponent(typeName string) (Component, error) {
	fn, ok := registeredComponents[typeName]
	if !ok {
		return nil, fmt.Errorf("components: component %s is not registered", typeName)
	}

	return fn(), nil
}
