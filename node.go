package main

type Expr interface {
	// Returns the corresponding token to this node.
	token() *Token
	// Output assembly corresponding to this node.
	emit()
}

type FunctionDecl struct {
	tok        *Token
	name       string
	returnType *Type
	body       *Block
	scope      *Scope
}

type Block struct {
	tok  *Token
	body []Expr
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

// Variable is considered a tag for a memory region with type information.
// `offset` is determined in code generation step.
type Variable struct {
	tok *Token
	// Offset from stack pointer after function's prelude.
	offset int
	ty     *Type
}

type Identifier struct {
	tok      *Token
	variable *Variable
}

type IntLiteral struct {
	tok *Token
}

type BoolLiteral struct {
	tok   *Token
	value bool
}

type FunctionCall struct {
	tok      *Token
	function *FunctionDecl
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
func (node *FunctionCall) token() *Token { return node.tok }
