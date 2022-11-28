package main

type Node interface {
	// Returns the corresponding token to this node.
	token() *Token
}

type Expr interface {
	token() *Token
}

type IntLiteral struct {
	tok *Token
}

func (node *IntLiteral) token() *Token { return node.tok }
