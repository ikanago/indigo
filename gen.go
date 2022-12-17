package main

import (
	"fmt"
)

const fp = "x29"

var argumentRegisters = []string{"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7"}

func Generate(ast *Ast) {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align 2")
	for _, node := range ast.funcs {
		fmt.Println()
		node.emit()
	}
}

func code(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Printf("\t%s\n", s)
}

func comment(msg string, a ...any) {
	s := fmt.Sprintf(msg, a...)
	fmt.Printf("\t;%s\n", s)
}

func (expr *FunctionDecl) emit() {
	var functionName string
	if expr.name == "main" {
		functionName = "_main"
	} else {
		functionName = expr.name
	}
	fmt.Printf(".globl %s\n", functionName)
	fmt.Printf("%s:\n", functionName)

	save_frame_pointer_and_link_register()

	totalOffset := 0
	for name, expr := range expr.scope.exprs {
		if variable, ok := expr.(*Variable); ok {
			variable.offset = totalOffset
			comment("offset of %s: %d", name, variable.offset)
			totalOffset += variable.ty.GetSize()
		}
	}

	code("sub sp, sp, #%d", totalOffset)
	code("mov %s, sp", fp)

	for i, parameter := range expr.parameters {
		code("str %s, [%s, #%d]", argumentRegisters[i], fp, parameter.offset)
	}
	expr.body.emit()

	code("add sp, sp, #%d", totalOffset)
	restore_frame_pointer_and_link_register()
	code("ret")
}

func save_frame_pointer_and_link_register() {
	code("stp %s, x30, [sp, -32]!", fp)
}

func restore_frame_pointer_and_link_register() {
	code("ldp %s, x30, [sp], 32", fp)
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
	expr.lhs.emit()
	expr.rhs.emit()
	comment("assign")
	generatePop("x1")
	generatePop("x2")
	code("str x1, [x2]")
}

func (expr *AddOp) emit() {
	expr.lhs.emit()
	expr.rhs.emit()
	comment("add")
	generatePop("x2")
	generatePop("x1")
	code("add x0, x1, x2")
	generatePush("x0")
}

func (expr *Variable) emit() {
	comment("variable: %s", expr.Name())
	code("add x0, %s, #%d", fp, expr.offset)
	generatePush("x0")
}

func (expr *Identifier) emit() {
	comment("identifier: %s", expr.Name())
	code("add x1, %s, #%d", fp, expr.variable.offset)
	code("ldr x0, [x1]")
	generatePush("x0")
}

func (expr *IntLiteral) emit() {
	comment("int literal")
	code("mov x0, #%s", expr.tok.value)
	generatePush("x0")
}

func (expr *BoolLiteral) emit() {
	comment("bool literal")
	if expr.value {
		code("mov x0, #1")
	} else {
		code("mov x0, #0")
	}
	generatePush("x0")
}

func (expr *FunctionCall) emit() {
	comment("function call")
	for i := len(expr.arguments) - 1; i >= 0; i-- {
		expr.arguments[i].emit()
		generatePop("x0")
		code("mov x%d, x0", i)
	}
	code("bl %s", expr.function.name)
	generatePush("x0")
	comment("function call end")
}

func generatePush(register string) {
	code("str %s, [sp, #-16]!", register)
}

func generatePop(register string) {
	code("ldr %s, [sp], #16", register)
}
