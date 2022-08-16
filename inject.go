package inject

import (
	"bytes"
	"fmt"
	"reflect"
)

// Logger allows for simple logging as inject traverses
// and populates the object graph.
type Logger interface {
	Infof(format string, args ...any)
}

type objectType int

const (
	structPointer objectType = iota + 1
	function
)

// Object is the abstraction of each node in Graph, and the
// objects present in the graph are marked by Object.
type Object struct {
	// The specific object that will be used in the user
	// code can be a structure pointer or a function.
	Value any

	// Use name to distinguish between two objects of the same type.
	Name string

	// Use reflection to get the corresponding Type and Value representation.
	reflectType  reflect.Type
	reflectValue reflect.Value

	// Mark Object type
	t objectType

	// complete marks whether this structure is ready to be reused for injection.
	complete bool
}

// String impl String interface for fmt.Println
func (o *Object) String() string {
	b := &bytes.Buffer{}
	fmt.Fprint(b, o.reflectType)

	return b.String()
}

// Graph uses the incoming structures as a starting
// point to build a graph of the dependencies of
// each structure. And you can assemble them to complete
// the dependency injection.
//
// It should be noted that there is only one global
// instance of each object.
type Graph struct {
	// logger is optional, will output some info in build graph
	logger Logger

	// Use the Provide function to provide a starting point for creating objects.
	startPoint []*Object
	// Maintain a global pool of objects.
	namePool    map[string]*Object
	unnamedPool map[reflect.Type]*Object
}

// Provide registers the starting point for the graph and
// can distinguish different objects of the same type by
// the Name field in Object.
func (g *Graph) Provide(obj ...*Object) error {
	// init object pool
	if g.unnamedPool == nil {
		g.unnamedPool = make(map[reflect.Type]*Object)
	}
	if g.namePool == nil {
		g.namePool = make(map[string]*Object)
	}

	for _, v := range obj {
		// Assemble the object object and set the incoming entry point type and other information.
		v.reflectType = reflect.TypeOf(v)
		v.reflectValue = reflect.ValueOf(v)
		err := isCorrectType(v, v.reflectType)
		if err != nil {
			return err
		}

		// object don't have name
		if v.Name == "" {
			if _, ok := g.unnamedPool[v.reflectType]; ok {
				return fmt.Errorf(
					"unsupport two object have same type, graph already has %v with value %v",
					v.reflectType,
					g.unnamedPool[v.reflectType],
				)
			}

			g.unnamedPool[v.reflectType] = v
			g.startPoint = append(g.startPoint, v)

			continue
		}

		// object with name
		if _, ok := g.namePool[v.Name]; ok {
			return fmt.Errorf(
				"unsupport two object have same name, graph already has %v with name %v",
				v.reflectType,
				v.Name,
			)
		}

		g.namePool[v.Name] = v
		g.startPoint = append(g.startPoint, v)
	}

	return nil
}

// GraphConfig is a parameter to the initialization
// function that configures certain fields in the Graph.
type GraphConfig struct {
	Logger Logger
}

// NewGraph returns a new Graph for dependency injection.
func NewGraph(c *GraphConfig) *Graph {
	return &Graph{
		logger: c.Logger,
	}
}
