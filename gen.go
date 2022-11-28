package main

import "fmt"

func Generate(ast *Ast) {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align	2")
	fmt.Println(".globl _main")
	fmt.Println("_main:")
	ast.nodes[0].emit()
	generatePop("x0")
	fmt.Println("\tret")
}

func (expr *ShortVarDecl) emit() {
}

func (expr *AddOp) emit() {
	expr.lhs.emit()
	expr.rhs.emit()
	generatePop("x1")
	generatePop("x2")
	fmt.Printf("\tadd x0, x1, x2\n")
	generatePush("x0")
}

func (expr *Variable) emit() {
}

func (expr *IntLiteral) emit() {
	fmt.Printf("\tmov x0, #%s\n", expr.tok.value)
	generatePush("x0")
}

func generatePush(register string) {
	fmt.Printf("\tstr %s, [sp, #-16]!\n", register)
}

func generatePop(register string) {
	fmt.Printf("\tldr %s, [sp], #16\n", register)
}
