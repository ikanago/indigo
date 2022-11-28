package main

type Expr interface {
	// Returns the corresponding token to this node.
	token() *Token
	// Output assembly corresponding to this node.
	emit()
}

type AddOp struct {
	tok *Token
	lhs Expr
	rhs Expr
}

type IntLiteral struct {
	tok *Token
}

func (node *AddOp) token() *Token      { return node.tok }
func (node *IntLiteral) token() *Token { return node.tok }
