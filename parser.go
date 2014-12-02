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
}

// parses vcs file to internal data structures (Element)
func (parser *PHPParser) Elements(file *vcs.File) []Element {
	code := file.Content()

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
	endNodes := parser.readBlockIntoCfg(cfg, block, []*gs.Vertex{cfg.CreateAndAddToGraph("start")})

	//add final return statement if missing
	exitNode := cfg.CreateAndAddToGraph(parser.createId("return", cfg))
	for _, endNode := range endNodes {
		if strings.HasSuffix(endNode.ID, "return") == false {
			cfg.Connect(endNode, exitNode, 1)
		}
	}

	return cfg
}

// reads a block into a given control flow graph
func (parser *PHPParser) readBlockIntoCfg(cfg *gs.Graph, block *ast.Block, startNodes []*gs.Vertex) []*gs.Vertex {

	var endNodes []*gs.Vertex



	for _, statement := range block.Statements {

		switch t := statement.(type) {

		case ast.ExpressionStmt, ast.EchoStmt, ast.BreakStmt:
			endNodes = parser.readSimpleStmtIntoCfg(cfg, fmt.Sprintf("%T", statement), startNodes)

		case ast.ReturnStmt, ast.ThrowStmt:
			//return statements couldn't be followed by another node, so no endNodes will be empty
			endNodes = []*gs.Vertex{}
			parser.readSimpleStmtIntoCfg(cfg, fmt.Sprintf("%T", statement), startNodes)

		case *ast.IfStmt:
			endNodes = parser.readIfStmtIntoCfg(cfg, statement.(*ast.IfStmt), startNodes)

		case *ast.SwitchStmt:
			endNodes = parser.readSwitchStmtIntoCfg(cfg, statement.(*ast.SwitchStmt), startNodes)

		case ast.SwitchStmt:
			switchStatement := statement.(ast.SwitchStmt)
			endNodes = parser.readSwitchStmtIntoCfg(cfg, &switchStatement, startNodes)

		case *ast.ForeachStmt:
			endNodes = parser.readLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.ForeachStmt).LoopBlock.(*ast.Block), startNodes)

		case *ast.ForStmt:
			endNodes = parser.readLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.ForStmt).LoopBlock.(*ast.Block), startNodes)

		case *ast.WhileStmt:
			endNodes = parser.readLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.WhileStmt).LoopBlock.(*ast.Block), startNodes)

		default:
			log.Fatalf("Unhandled type %T", t)

		}

		startNodes = endNodes
	}

	return endNodes
}

func (parser *PHPParser) readLoopIntoCfg(cfg *gs.Graph, label string, block *ast.Block, startNodes []*gs.Vertex) []*gs.Vertex {

	id := parser.createId(label, cfg)
	headNode := cfg.CreateAndAddToGraph(id)

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, headNode, 1)
		}
	}

	endNodes := parser.readBlockIntoCfg(cfg, block, []*gs.Vertex{headNode})
	footNode := cfg.CreateAndAddToGraph(id + "_end")
	cfg.Connect(footNode, headNode, 1)

	if len(endNodes) > 0 {
		for _, endNode := range endNodes {
			cfg.Connect(endNode, footNode, 1)
		}
	}

	return []*gs.Vertex{headNode}
}

// reads a simple statement into a control flow graph
func (parser *PHPParser) readSimpleStmtIntoCfg(cfg *gs.Graph, label string, startNodes []*gs.Vertex) []*gs.Vertex {

	node := cfg.CreateAndAddToGraph(parser.createId(label, cfg))

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node, 1)
		}
	}

	return []*gs.Vertex{node}
}


func (parser *PHPParser) readSwitchStmtIntoCfg(cfg *gs.Graph, switchStmt *ast.SwitchStmt, startNodes []*gs.Vertex) []*gs.Vertex {
	node := cfg.CreateAndAddToGraph(parser.createId("switch", cfg))

	endNodes := []*gs.Vertex{}
	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node, 1)
		}
	}

	openNodes := []*gs.Vertex{}

	//handle switch case
	for _, switchCase := range switchStmt.Cases {
		caseNode := cfg.CreateAndAddToGraph(parser.createId("case", cfg));
		cfg.Connect(node, caseNode, 1)

		if len(openNodes) > 0 {
			for _, openNode := range openNodes {
				cfg.Connect(openNode, caseNode, 1)
			}
			openNodes = []*gs.Vertex{}
		}

		caseNodes := parser.readBlockIntoCfg(cfg, &switchCase.Block, []*gs.Vertex{caseNode})
		//for an empty switch-case
		if len(caseNodes) == 0 {
			openNodes = append(openNodes, caseNode)
		}

		for _, caseNode := range caseNodes {
			if strings.HasSuffix(caseNode.ID, "break") == false {
				openNodes = append(openNodes, caseNode)
			}else {
				endNodes = append(endNodes, caseNode)
			}
		}

	}

	if switchStmt.DefaultCase != nil {
		//handle default case
		defaultNode := cfg.CreateAndAddToGraph(parser.createId("default", cfg));
		cfg.Connect(node, defaultNode, 1)
		if len(openNodes) > 0 {
			for _, openNode := range openNodes {
				cfg.Connect(openNode, defaultNode, 1)
			}
			openNodes = []*gs.Vertex{}
		}
		defaultCaseNode := parser.readBlockIntoCfg(cfg, switchStmt.DefaultCase, []*gs.Vertex{defaultNode})

		endNodes = append(endNodes, defaultCaseNode...)
	}

	return endNodes
}

//reads a if statement into given cfg struct
func (parser *PHPParser) readIfStmtIntoCfg(cfg *gs.Graph, ifStmt *ast.IfStmt, startNodes []*gs.Vertex) []*gs.Vertex {

	node := cfg.CreateAndAddToGraph(parser.createId("if", cfg))
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

func (parser *PHPParser) createId(label string, cfg *gs.Graph) string {
	label = strings.ToLower(label)
	label = strings.Replace(label, "*", "", -1)
	label = strings.Replace(label, "ast.", "", -1)
	label = strings.Replace(label, "stmt", "", -1)
	return fmt.Sprintf("n%d_%s", cfg.GetEdgesSize()+1, label)
}
