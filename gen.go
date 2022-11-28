package main

import "fmt"

func Generate(ast *Ast) {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align	2")
	fmt.Println(".globl _main")
	fmt.Println("_main:")
	generateNode(ast.root)
	fmt.Println("  ret")
}

func generateNode(node Node) {
	switch node.(type) {
	case *IntLiteral:
		fmt.Printf("  mov x0, %s\n", node.token().value)
	}
}
