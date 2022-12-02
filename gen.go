package main

import "fmt"

func Generate(ast *Ast) {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align	2")
	fmt.Println(".globl _main")
	fmt.Println("_main:")
	fmt.Printf("\tsub sp, sp, #%d\n", ast.localEnv.totalOffset)
	fmt.Println("\tmov x7, sp")
	for _, node := range ast.nodes {
		node.emit()
	}
	fmt.Printf("\tadd sp, sp, #%d\n", ast.localEnv.totalOffset)
	fmt.Println("\tret")
}

func comment(msg string) {
	fmt.Printf("\t;%s\n", msg)
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

func generatePush(register string) {
	fmt.Printf("\tstr %s, [sp, #-16]!\n", register)
}

func generatePop(register string) {
	fmt.Printf("\tldr %s, [sp], #16\n", register)
}
