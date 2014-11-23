package analyzer

import (

	"github.com/jochil/vcs"
	"log"
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php"
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

	//TODO remove, just here for testing
	code = `<?php

function dummy($start, $end){

	$sum = 0;
	for( $i = $start; $i <= $end; $i++ ){
		$sum += $i;
	}
	return $sum;
}

echo dummy(1, 3);`
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
			element := parser.readFunction(node.(*ast.FunctionStmt))
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
		functionElement := Function{}
		functionElement.Name = method.Name

		methods = append(methods, functionElement)
	}

	element.Methods = methods

	return element
}

// convert function data structure of the language specific parser to the internal data structure
func (parser *PHPParser) readFunction(function *ast.FunctionStmt) Function {
	element := Function{}
	element.Name = function.Name
	element.CFG = parser.buildCFG(function.Body)

	return element
}

// creating the control flow graph for a block struct from language specific parser
func (parser *PHPParser) buildCFG(block *ast.Block) (ControlFlowGraph) {
	cfg := NewControlFlowGraph()
	parser.readBlockIntoCfg(cfg, block, nil)
	return *cfg
}

// reads a block into a given control flow graph
func (parser *PHPParser) readBlockIntoCfg(cfg *ControlFlowGraph, block *ast.Block, startNodes []*Node) ([]*Node) {

	var endNodes []*Node

	for _, statement := range block.Statements {

		switch t := statement.(type) {

		case ast.ExpressionStmt, ast.EchoStmt:
			endNodes = parser.readSimpleStmtIntoCfg(cfg, &statement, startNodes)

		case ast.ReturnStmt:
			//return statements couldn't be followed by another node, so no endNodes will be empty
			parser.readSimpleStmtIntoCfg(cfg, &statement, startNodes)

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
func (parser *PHPParser) readSimpleStmtIntoCfg(cfg *ControlFlowGraph, statement *ast.Statement, startNodes []*Node) ([]*Node) {

	node := NewNode("statement")
	cfg.Add(node)

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node)
		}
	}

	return []*Node{node}
}

//reads a if statement into given cfg struct
func (parser *PHPParser) readIfStmtIntoCfg(cfg *ControlFlowGraph, ifStmt *ast.IfStmt, startNodes []*Node) ([]*Node) {

	node := NewNode("if")
	cfg.Add(node)

	endNodes := []*Node{}

	// connect end nodes
	if len(startNodes) > 0 {
		for _, parentNode := range startNodes {
			cfg.Connect(parentNode, node)
		}
	}

	// handle true branch
	trueBranch := ifStmt.TrueBranch
	switch ifTrueType := trueBranch.(type){
	case *ast.Block:
		trueEndNodes := parser.readBlockIntoCfg(cfg, trueBranch.(*ast.Block), []*Node{node})
		endNodes = append(endNodes, trueEndNodes...)

	default:
		log.Fatalf("invalid if branch of type %T", ifTrueType)
	}


	//handle false branch
	falseBranch := ifStmt.FalseBranch
	switch ifFalseType := falseBranch.(type){

	case ast.Block: //no/empty else
		endNodes = append(endNodes, node)

	case *ast.Block: //else
		falseEndNodes := parser.readBlockIntoCfg(cfg, falseBranch.(*ast.Block), []*Node{node})
		endNodes = append(endNodes, falseEndNodes...)

	case *ast.IfStmt: //elseif
		falseEndNodes := parser.readIfStmtIntoCfg(cfg, falseBranch.(*ast.IfStmt), []*Node{node})
		endNodes = append(endNodes, falseEndNodes...)

	default:
		log.Fatalf("invalid if branch of type %T", ifFalseType)
	}

	return endNodes
}

//TODO remove, just here for debugging & testing
func TestParser(repo *vcs.Repository) {

	for _, commit := range repo.Commits {
		for _, file := range commit.Files {
			if Filter.ValidExtension(file.Path) {

				parser := PHPParser{}
				parser.Elements(file)

				return
			}
		}
	}

}
