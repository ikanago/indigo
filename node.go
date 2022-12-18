package main

type Expr interface {
	// Returns the corresponding token to this node.
	token() *Token
	// Output assembly corresponding to this node.
	emit()
}

type FunctionDecl struct {
	tok        *Token
	Name       string
	Parameters []*Variable
	ReturnType *Type
	Body       *Block
	Scope      *Scope
}

type Block struct {
	tok  *Token
	Body []Expr
}

type Return struct {
	tok  *Token
	Node Expr
}

type Assign struct {
	tok *Token
	Lhs Expr
	Rhs Expr
}

type AddOp struct {
	tok *Token
	Lhs Expr
	Rhs Expr
}

// Variable is considered a tag for a memory region with type information.
// `offset` is determined in code generation step.
type Variable struct {
	tok  *Token
	Name string
	// Offset from stack pointer after function's prelude.
	Offset int
	Ty     *Type
}

type Identifier struct {
	tok      *Token
	Name     string
	Variable *Variable
}

type IntLiteral struct {
	tok   *Token
	Value string
}

type BoolLiteral struct {
	tok   *Token
	Value bool
}

type FunctionCall struct {
	tok       *Token
	Function  *FunctionDecl
	Arguments []Expr
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

func (node *FunctionCall) Name() string {
	return node.token().value
}
