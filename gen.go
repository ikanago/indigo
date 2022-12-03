package main

import "fmt"

func Generate(ast *Ast) {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align	2")
	for _, node := range ast.funcs {
		node.emit()
	}
}

func comment(msg string) {
	fmt.Printf("\t;%s\n", msg)
}

func (expr *FunctionDecl) emit() {
	var name string
	if expr.name == "main" {
		name = "_main"
	} else {
		name = expr.name
	}
	fmt.Printf(".globl %s\n", name)
	fmt.Printf("%s:\n", name)

	body := expr.body.(*Block)
	totalOffset := body.localEnv.totalOffset
	fmt.Printf("\tsub sp, sp, #%d\n", totalOffset)
	fmt.Println("\tmov x7, sp")
	body.emit()
	fmt.Printf("\tadd sp, sp, #%d\n", totalOffset)
	fmt.Println("\tret")
}

func (expr *Block) emit() {
	for _, stmt := range expr.body {
		stmt.emit()
	}
}

func (expr *Return) emit() {
	expr.node.emit()
	generatePop("x0")
}

func (expr *Assign) emit() {
	comment("assign")
	expr.lhs.emit()
	expr.rhs.emit()
	generatePop("x1")
	generatePop("x2")
	fmt.Println("\tstr x1, [x2]")
}

func (expr *AddOp) emit() {
	comment("add")
	expr.lhs.emit()
	expr.rhs.emit()
	generatePop("x2")
	generatePop("x1")
	fmt.Printf("\tadd x0, x1, x2\n")
	generatePush("x0")
}

func (expr *Variable) emit() {
	comment("variable")
	fmt.Printf("\tadd x0, x7, #%d\n", expr.offset)
	generatePush("x0")
}

func (expr *Identifier) emit() {
	comment("identifier")
	fmt.Printf("\tadd x1, x7, #%d\n", expr.offset)
	fmt.Println("\tldr x0, [x1]")
	generatePush("x0")
}

func (expr *IntLiteral) emit() {
	comment("int literal")
	fmt.Printf("\tmov x0, #%s\n", expr.tok.value)
	generatePush("x0")
}

func (expr *BoolLiteral) emit() {
	comment("bool literal")
	if expr.value {
		fmt.Println("\tmov x0, #1")
	} else {
		fmt.Println("\tmov x0, #0")
	}
	generatePush("x0")
}

func generatePush(register string) {
	fmt.Printf("\tstr %s, [sp, #-16]!\n", register)
}

func generatePop(register string) {
	fmt.Printf("\tldr %s, [sp], #16\n", register)
}
