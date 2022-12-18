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
		dumped += dln(level+1, "name: %s", expr.Name)
		dumped += dln(level+1, "parameters: [")
		for _, parameter := range expr.Parameters {
			dumped += dumpExpr(level+2, parameter)
		}
		dumped += dln(level+1, "]")
		dumped += dln(level+1, "returnType: %s", dumpType(expr.ReturnType))
		dumped += dln(level+1, "body: \n%s", dumpExpr(level+2, expr.Body))
		dumped += dln(level, "}")
	case *Block:
		dumped += dln(level, "Block: {")
		for _, e := range expr.Body {
			dumped += dumpExpr(level+1, e)
		}
		dumped += dln(level, "}")
	case *Return:
		dumped += dln(level, "Return: {")
		dumped += dumpExpr(level+1, expr.Node)
		dumped += dln(level, "}")
	case *Assign:
		dumped += dln(level, "Assign: {")
		dumped += d(level+1, "lhs:\n%s", dumpExpr(level+2, expr.Lhs))
		dumped += d(level+1, "rhs:\n%s", dumpExpr(level+2, expr.Rhs))
		dumped += dln(level, "}")
	case *AddOp:
		dumped += dln(level, "AddOp: {")
		dumped += d(level+1, "lhs\n%s", dumpExpr(level+2, expr.Lhs))
		dumped += d(level+1, "rhs\n%s", dumpExpr(level+2, expr.Rhs))
		dumped += dln(level, "}")
	case *Variable:
		dumped += dln(level, "Variable: { name: %s, type: %s }", expr.Name, dumpType(expr.Ty))
	case *Identifier:
		dumped += dln(level, "Identifier: { name: %s, type: %s }", expr.Name(), dumpType(expr.Variable.Ty))
	case *IntLiteral:
		dumped += dln(level, "IntLiteral: %s", expr.token().value)
	case *BoolLiteral:
		dumped += dln(level, "BoolLiteral: %t", expr.Value)
	case *FunctionCall:
		dumped += dln(level, "FunctionCall: {")
		dumped += dln(level+1, "name: %s", expr.Function.Name)
		dumped += dln(level+1, "arguments: [")
		for _, argument := range expr.Arguments {
			dumped += dumpExpr(level+2, argument)
		}
		dumped += dln(level+1, "]")
		dumped += dln(level, "}")
	}
	return dumped
}

func dumpType(ty *Type) string {
	return ty.Name
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
