package analyzer

import (
	"fmt"
	"github.com/gyuho/goraph/graph/gs"
	"github.com/jochil/vcs"
	"github.com/stephens2424/php"
	"github.com/stephens2424/php/ast"
	"log"
	"strings"
)

//interface for encapsulating the language specific parser
type Parser interface {
	Elements(file vcs.File) []Element
}

// struct for the php parser (implemented against Parser interface)
type PHPParser struct {
	nodeCounter uint
}

// parses vcs file to internal data structures (Element)
func (parser *PHPParser) Elements(file *vcs.File) []Element {
	parser.nodeCounter = 0
	code := file.Content()

	log.Println(code)

	realParser := php.NewParser(code)
	nodes, err := realParser.Parse()

	if err != nil {
		log.Fatal(err)
	}

	elements := make([]Element, 0, 1)

	for _, node := range nodes {
		switch node.(type) {

		case ast.Class:
			element := parser.readClass(node.(ast.Class))
			elements = append(elements, &element)

		case *ast.FunctionStmt:
			function := node.(*ast.FunctionStmt)
			element := parser.readFunction(function.Name, function.Body)
			elements = append(elements, &element)

		}
	}
	return elements
}

// convert class data structure of the language specific parser to the internal data structure
func (parser *PHPParser) readClass(class ast.Class) Class {
	element := Class{}
	element.Name = class.Name

	methods := make([]Function, 0, 1)
	for _, method := range class.Methods {
		functionName := class.Name + "->" + method.Name
		methods = append(methods, parser.readFunction(functionName, method.Body))
	}

	element.Methods = methods

	return element
}

// convert function data structure of the language specific parser to the internal data structure
func (parser *PHPParser) readFunction(name string, body *ast.Block) Function {
	element := Function{}
	element.Name = name
	element.CFG = parser.buildCFG(body)
	dumpCFG(name, element.CFG)

	return element
}

// creating the control flow graph for a block struct from language specific parser
func (parser *PHPParser) buildCFG(block *ast.Block) *gs.Graph {
	cfg := gs.NewGraph()
	startNode := gs.NewVertex("start")
	cfg.AddVertex(startNode)
	parser.readBlockIntoCfg(cfg, block, []*gs.Vertex{startNode})
	return cfg
}

// reads a block into a given control flow graph
func (parser *PHPParser) readBlockIntoCfg(cfg *gs.Graph, block *ast.Block, startNodes []*gs.Vertex) []*gs.Vertex {

	var endNodes []*gs.Vertex

	for _, statement := range block.Statements {

		switch t := statement.(type) {

		case ast.ExpressionStmt, ast.EchoStmt:
			endNodes = parser.readSimpleStmtIntoCfg(cfg, fmt.Sprintf("%T", statement), startNodes)

		case ast.ReturnStmt:
			//return statements couldn't be followed by another node, so no endNodes will be empty
			endNodes = []*gs.Vertex{}
			parser.readSimpleStmtIntoCfg(cfg, fmt.Sprintf("%T", statement), startNodes)

		case *ast.IfStmt:
			endNodes = parser.readIfStmtIntoCfg(cfg, statement.(*ast.IfStmt), startNodes)

		default:
			log.Fatalf("Unhandled type %T", t)

		}

		startNodes = endNodes
	}

	return endNodes
}

// reads a simple statement into a control flow graph
func (parser *PHPParser) readSimpleStmtIntoCfg(cfg *gs.Graph, label string, startNodes []*gs.Vertex) []*gs.Vertex {

	node := gs.NewVertex(parser.createId(label))
	cfg.AddVertex(node)

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node, 1)
		}
	}

	return []*gs.Vertex{node}
}

//reads a if statement into given cfg struct
func (parser *PHPParser) readIfStmtIntoCfg(cfg *gs.Graph, ifStmt *ast.IfStmt, startNodes []*gs.Vertex) []*gs.Vertex {

	node := gs.NewVertex(parser.createId(fmt.Sprintf("%T", ifStmt)))

	cfg.AddVertex(node)

	endNodes := []*gs.Vertex{}

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node, 1)
		}
	}

	// handle true branch
	trueBranch := ifStmt.TrueBranch
	switch ifTrueType := trueBranch.(type) {
	case *ast.Block:
		trueEndNodes := parser.readBlockIntoCfg(cfg, trueBranch.(*ast.Block), []*gs.Vertex{node})
		endNodes = append(endNodes, trueEndNodes...)

	default:
		log.Fatalf("invalid if branch of type %T", ifTrueType)
	}

	//handle false branch
	falseBranch := ifStmt.FalseBranch
	switch ifFalseType := falseBranch.(type) {

	case ast.Block: //no/empty else
		endNodes = append(endNodes, node)

	case *ast.Block: //else
		falseEndNodes := parser.readBlockIntoCfg(cfg, falseBranch.(*ast.Block), []*gs.Vertex{node})
		endNodes = append(endNodes, falseEndNodes...)

	case *ast.IfStmt: //elseif
		falseEndNodes := parser.readIfStmtIntoCfg(cfg, falseBranch.(*ast.IfStmt), []*gs.Vertex{node})
		endNodes = append(endNodes, falseEndNodes...)

	default:
		log.Fatalf("invalid if branch of type %T", ifFalseType)
	}

	return endNodes
}

func (parser *PHPParser) createId(label string) string {
	label = strings.Replace(label, "*", "", -1)
	label = strings.Replace(label, "ast.", "", -1)
	parser.nodeCounter++
	return fmt.Sprintf("%s%d", label, parser.nodeCounter)
}
