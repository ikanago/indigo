package main

type Expr interface {
	// Returns the corresponding token to this node.
	token() *Token
	// Output assembly corresponding to this node.
	emit()
}

type ShortVarDecl struct {
	tok *Token
	lhs Expr
	rhs Expr
}

type AddOp struct {
	tok *Token
	lhs Expr
	rhs Expr
}

type Variable struct {
	tok    *Token
	offset int
}

type IntLiteral struct {
	tok *Token
}

func (node *ShortVarDecl) token() *Token { return node.tok }
func (node *AddOp) token() *Token        { return node.tok }
func (node *Variable) token() *Token     { return node.tok }
func (node *IntLiteral) token() *Token   { return node.tok }
