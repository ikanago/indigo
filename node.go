package main

type Expr interface {
	// Returns the corresponding token to this node.
	token() *Token
	// Output assembly corresponding to this node.
	emit()
}

type FunctionDecl struct {
	tok  *Token
	name string
	body Expr
}

type Block struct {
	tok      *Token
	body     []Expr
	localEnv *LocalEnv
}

type Return struct {
	tok  *Token
	node Expr
}

type Assign struct {
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
	tok *Token
	// Offset from stack pointer after function's prelude.
	offset int
}

type Identifier struct {
	tok    *Token
	offset int
}

type IntLiteral struct {
	tok *Token
}

type BoolLiteral struct {
	tok   *Token
	value bool
}

func (node *FunctionDecl) token() *Token { return node.tok }
func (node *Block) token() *Token        { return node.tok }
func (node *Return) token() *Token       { return node.tok }
func (node *Assign) token() *Token       { return node.tok }
func (node *AddOp) token() *Token        { return node.tok }
func (node *Variable) token() *Token     { return node.tok }
func (node *Identifier) token() *Token   { return node.tok }
func (node *IntLiteral) token() *Token   { return node.tok }
func (node *BoolLiteral) token() *Token  { return node.tok }
