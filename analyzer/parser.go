package analyzer

import (
	"fmt"
	"github.com/gyuho/goraph/graph/gs"
	"github.com/jochil/scabov/vcs"
	"github.com/stephens2424/php"
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/lexer"
	"github.com/stephens2424/php/token"
	"log"
	"strings"
)

//interface for encapsulating the language specific parser
type Parser interface {
	Elements(file *vcs.File) []Element
	Functions(file *vcs.File) map[string]Function
	UpdateLanguageUsage(langUsage LanguageUsage, file *vcs.File)
}

func NewParser() Parser {
	if Filter.Lang() == vcs.PHP {
		return &PHPParser{}
	}
	return nil
}

// struct for the php parser (implemented against Parser interface)
type PHPParser struct {
}

//TODO add to interface
func (parser *PHPParser) UpdateLanguageUsage(langUsage LanguageUsage, file *vcs.File) {
	code := file.Content()

	lex := lexer.NewLexer(code)
	for {
		if item := lex.Next(); item.Typ == token.EOF {
			break
		} else {
			langUsage.(*PHPLanguageUsage).AddItem(item)
		}
	}

}

//TODO remove redundant code (see Elements())
func (parser *PHPParser) Functions(file *vcs.File) map[string]Function {
	nodes := parser.parseFile(file)
	functions := map[string]Function{}

	for _, node := range nodes {
		switch node.(type) {

		case ast.Class, *ast.Class:
			element := parser.readClass(node.(*ast.Class))
			for _, function := range element.Methods {
				functions[function.Name] = function
			}

		case *ast.FunctionStmt:
			function := node.(*ast.FunctionStmt)
			element := parser.readFunction(function.Name, function.Body)
			functions[element.Name] = element

		}
	}
	return functions
}

// parses vcs file to internal data structures (Element)
func (parser *PHPParser) Elements(file *vcs.File) []Element {

	nodes := parser.parseFile(file)
	elements := make([]Element, 0, 1)

	for _, node := range nodes {
		switch node.(type) {

		case ast.Class, *ast.Class:
			element := parser.readClass(node.(*ast.Class))
			elements = append(elements, &element)

		case *ast.FunctionStmt:
			function := node.(*ast.FunctionStmt)
			element := parser.readFunction(function.Name, function.Body)
			elements = append(elements, &element)

		}
	}
	return elements
}

func (parser *PHPParser) parseFile(file *vcs.File) []ast.Node {
	code := file.Content()
	realParser := php.NewParser(code)
	nodes, err := realParser.Parse()

	//log.Println(file.StoragePath)

	if err != nil {
		//log.Fatal(err)
		return []ast.Node{}
	}
	return nodes
}

// convert class data structure of the language specific parser to the internal data structure
func (parser *PHPParser) readClass(class *ast.Class) Class {
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

	if body != nil {
		element.NumNodes = countNodes(body.Children())
	}
	return element
}

func countNodes(nodes []ast.Node) int {
	count := len(nodes)
	for _, node := range nodes {
		if node != nil {
			count += countNodes(node.Children())
		}
	}
	return count
}

// creating the control flow graph for a block struct from language specific parser
func (parser *PHPParser) buildCFG(block *ast.Block) *gs.Graph {
	cfg := gs.NewGraph()
	startNode := cfg.CreateAndAddToGraph("start")

	endNodes := parser.readBlockIntoCfg(cfg, block, []*gs.Vertex{startNode})
	exitNode := cfg.CreateAndAddToGraph("exit")
	//connect all endNodes to the exit node
	for _, endNode := range endNodes {
		cfg.Connect(endNode, exitNode, 1)
	}
	cfg.Connect(exitNode, startNode, 1)
	return cfg
}

func (parser *PHPParser) readStatementIntoCfg(cfg *gs.Graph, statement ast.Statement, startNodes []*gs.Vertex) []*gs.Vertex {

	var endNodes []*gs.Vertex

	switch t := statement.(type) {

	case ast.Block:
		endNodes = startNodes

	case *ast.Block:
		endNodes = parser.readBlockIntoCfg(cfg, statement.(*ast.Block), startNodes)

	case ast.ExpressionStmt, ast.EchoStmt, ast.BreakStmt, *ast.BreakStmt, ast.ReturnStmt, ast.ThrowStmt, *ast.ReturnStmt,
		*ast.EmptyStatement, *ast.ExitStmt, ast.AssignmentExpression, ast.BinaryExpression:
		endNodes = parser.readSimpleStmtIntoCfg(cfg, fmt.Sprintf("%T", statement), startNodes)

	case *ast.IfStmt:
		endNodes = parser.readIfStmtIntoCfg(cfg, statement.(*ast.IfStmt), startNodes)

	case *ast.SwitchStmt:
		endNodes = parser.readSwitchStmtIntoCfg(cfg, statement.(*ast.SwitchStmt), startNodes)

	case ast.SwitchStmt:
		switchStatement := statement.(ast.SwitchStmt)
		endNodes = parser.readSwitchStmtIntoCfg(cfg, &switchStatement, startNodes)

	case *ast.ForeachStmt:
		endNodes = parser.readHeadLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.ForeachStmt).LoopBlock.(ast.Statement), startNodes)

	case *ast.ForStmt:
		endNodes = parser.readHeadLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.ForStmt).LoopBlock.(ast.Statement), startNodes)

	case *ast.WhileStmt:
		endNodes = parser.readHeadLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.WhileStmt).LoopBlock.(ast.Statement), startNodes)

	case *ast.DoWhileStmt:
		endNodes = parser.readFootLoopIntoCfg(cfg, fmt.Sprintf("%T", statement), statement.(*ast.DoWhileStmt).LoopBlock.(ast.Statement), startNodes)

	case *ast.TryStmt:
		endNodes = parser.readTryCatchIntoCfg(cfg, statement.(*ast.TryStmt), startNodes)

	case *ast.ContinueStmt:
		//TODO implement

	default:
		log.Fatalf("Unhandled type %T", t)

	}
	return endNodes
}

// reads a block into a given control flow graph
func (parser *PHPParser) readBlockIntoCfg(cfg *gs.Graph, block *ast.Block, startNodes []*gs.Vertex) []*gs.Vertex {

	var endNodes []*gs.Vertex

	if block != nil {

		for _, statement := range block.Statements {
			endNodes = parser.readStatementIntoCfg(cfg, statement, startNodes)
			startNodes = endNodes
		}
	}
	return endNodes
}

func (parser *PHPParser) readFootLoopIntoCfg(cfg *gs.Graph, label string, statement ast.Statement, startNodes []*gs.Vertex) []*gs.Vertex {

	id := parser.createId(label, cfg)
	headNode := cfg.CreateAndAddToGraph(id)

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, headNode, 1)
		}
	}

	endNodes := parser.readStatementIntoCfg(cfg, statement, []*gs.Vertex{headNode})

	footNode := cfg.CreateAndAddToGraph(id + "_end")
	cfg.Connect(footNode, headNode, 1)

	if len(endNodes) > 0 {
		for _, endNode := range endNodes {
			cfg.Connect(endNode, footNode, 1)
		}
	}

	return []*gs.Vertex{headNode}
}

func (parser *PHPParser) readHeadLoopIntoCfg(cfg *gs.Graph, label string, statement ast.Statement, startNodes []*gs.Vertex) []*gs.Vertex {

	id := parser.createId(label, cfg)
	headNode := cfg.CreateAndAddToGraph(id)

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, headNode, 1)
		}
	}

	endNodes := parser.readStatementIntoCfg(cfg, statement, []*gs.Vertex{headNode})

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
		caseNode := cfg.CreateAndAddToGraph(parser.createId("case", cfg))
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
			} else {
				endNodes = append(endNodes, caseNode)
			}
		}

	}

	if switchStmt.DefaultCase != nil {
		//handle default case
		defaultNode := cfg.CreateAndAddToGraph(parser.createId("default", cfg))
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

func (parser *PHPParser) readTryCatchIntoCfg(cfg *gs.Graph, tryStmt *ast.TryStmt, startNodes []*gs.Vertex) []*gs.Vertex {
	node := cfg.CreateAndAddToGraph(parser.createId("if", cfg))
	endNodes := []*gs.Vertex{}

	// connect nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node, 1)
		}
	}

	tryEndNodes := parser.readBlockIntoCfg(cfg, tryStmt.TryBlock, []*gs.Vertex{node})

	allCatchEndNodes := []*gs.Vertex{}
	for _, catchStmt := range tryStmt.CatchStmts {
		catchEndNodes := parser.readBlockIntoCfg(cfg, catchStmt.CatchBlock, []*gs.Vertex{node})
		allCatchEndNodes = append(allCatchEndNodes, catchEndNodes...)
	}

	endNodes = append(endNodes, tryEndNodes...)
	endNodes = append(endNodes, allCatchEndNodes...)

	if tryStmt.FinallyBlock != nil {
		finallyEndNodes := parser.readBlockIntoCfg(cfg, tryStmt.FinallyBlock, endNodes)
		endNodes = finallyEndNodes
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
	trueEndNodes := parser.readStatementIntoCfg(cfg, ifStmt.TrueBranch, []*gs.Vertex{node})
	endNodes = append(endNodes, trueEndNodes...)

	falseEndNodes := parser.readStatementIntoCfg(cfg, ifStmt.FalseBranch, []*gs.Vertex{node})
	endNodes = append(endNodes, falseEndNodes...)

	return endNodes
}

func (parser *PHPParser) createId(label string, cfg *gs.Graph) string {
	label = strings.ToLower(label)
	label = strings.Replace(label, "*", "", -1)
	label = strings.Replace(label, "ast.", "", -1)
	label = strings.Replace(label, "stmt", "", -1)
	return fmt.Sprintf("n%d_%s", cfg.GetEdgesSize()+1, label)
}
