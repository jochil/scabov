package analyzer

import (
	"fmt"
	"github.com/gyuho/goraph/graph/gs"
)

type Element interface {
	String() string
}

// Struct representing a single function/method parameter
type Parameter struct {
	Name string
}

func (parameter *Parameter) String() string {
	return fmt.Sprintf("Parameter %s", parameter.Name)
}

// Struct for abstraction of a single method or function
type Function struct {
	Name       string
	Parameters []Parameter
	NumNodes   int
	CFG        *gs.Graph
	Hash       string
}

func (function *Function) String() string {
	return fmt.Sprintf("Function %s", function.Name)
}

// Struct Representing a single Class
type Class struct {
	Name    string
	Methods []Function
}

func (class *Class) String() string {
	return fmt.Sprintf("Class %s", class.Name)
}
