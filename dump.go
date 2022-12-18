package main

import (
	"fmt"
	"strings"
)

func (ast *Ast) Dump() string {
	dumped := ""
	for _, function := range ast.funcs {
		dumped += dumpExpr(0, function)
	}
	return dumped
}

func dumpExpr(level int, expr Expr) string {
	dumped := ""
	switch expr := expr.(type) {
	case *FunctionDecl:
		dumped += dln(level, "FunctionDecl: {")
		dumped += dln(level+1, "name: %s", expr.name)
		dumped += dln(level+1, "parameters: [")
		for _, parameter := range expr.parameters {
			dumped += dumpExpr(level+2, parameter)
		}
		dumped += dln(level+1, "]")
		dumped += dln(level+1, "returnType: %s", dumpType(expr.returnType))
		dumped += dln(level+1, "body: \n%s", dumpExpr(level+2, expr.body))
		dumped += dln(level, "}")
	case *Block:
		dumped += dln(level, "Block: {")
		for _, e := range expr.body {
			dumped += dumpExpr(level+1, e)
		}
		dumped += dln(level, "}")
	case *Return:
		dumped += dln(level, "Return: {")
		dumped += dumpExpr(level+1, expr.node)
		dumped += dln(level, "}")
	case *Assign:
		dumped += dln(level, "Assign: {")
		dumped += d(level+1, "lhs:\n%s", dumpExpr(level+2, expr.lhs))
		dumped += d(level+1, "rhs:\n%s", dumpExpr(level+2, expr.rhs))
		dumped += dln(level, "}")
	case *AddOp:
		dumped += dln(level, "AddOp: {")
		dumped += d(level+1, "lhs\n%s", dumpExpr(level+2, expr.lhs))
		dumped += d(level+1, "rhs\n%s", dumpExpr(level+2, expr.rhs))
		dumped += dln(level, "}")
	case *Variable:
		dumped += dln(level, "Variable: { name: %s, type: %s }", expr.Name(), dumpType(expr.ty))
	case *Identifier:
		dumped += dln(level, "Identifier: { name: %s, type: %s }", expr.Name(), dumpType(expr.variable.ty))
	case *IntLiteral:
		dumped += dln(level, "IntLiteral: %s", expr.token().value)
	case *BoolLiteral:
		dumped += dln(level, "BoolLiteral: %t", expr.value)
	case *FunctionCall:
		dumped += dln(level, "FunctionCall: {")
		dumped += dln(level+1, "name: %s", expr.function.name)
		dumped += dln(level+1, "arguments: [")
		for _, argument := range expr.arguments {
			dumped += dumpExpr(level+2, argument)
		}
		dumped += dln(level+1, "]")
		dumped += dln(level, "}")
	}
	return dumped
}

func dumpType(ty *Type) string {
	return ty.name
}

func dln(level int, s string, a ...any) string {
	return fmt.Sprintln(d(level, s, a...))
}

func d(level int, s string, a ...any) string {
	const width = 2
	indent := strings.Repeat(" ", level*width)
	r := fmt.Sprintf(s, a...)
	return fmt.Sprintf("%s%s", indent, r)
}
