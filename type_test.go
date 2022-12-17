package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeResolveVariable(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nx := 1\nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].body.body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.lhs.(*Variable).ty)
	ret := ast.funcs[0].body.body[1].(*Return)
	assert.Equal(t, &TypeInt, ret.node.(*Identifier).variable.ty)
}

func TestTypeResolveAdd(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nx := 1 \nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].body.body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.lhs.(*Variable).ty)
}

func TestTypeCallFunctionWithoutArgument(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nx := f() + 1 \nreturn x\n}\n func f() int {\nreturn 2\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].body.body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.lhs.(*Variable).ty)
	fuctionCall := assign.rhs.(*AddOp).lhs.(*FunctionCall)
	assert.Equal(t, &TypeInt, fuctionCall.function.returnType)

	x, ok := ast.funcs[0].scope.GetExpr("x")
	assert.True(t, ok)
	assert.Equal(t, 0, x.(*Variable).offset)
}

func TestTypeCallFunctionWithArgument(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nx := f(1)\nreturn x\n}\nfunc f(a int) int {\nreturn a\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	a, ok := ast.funcs[1].scope.GetExpr("a")
	assert.True(t, ok)
	assert.Equal(t, 0, a.(*Variable).offset)
}

func TestReturnUndefinedVariable(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nreturn abc\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "undefined: abc")
}

func TestNotEnoughReturnType(t *testing.T) {
	stream, _ := Tokenize("func f() bool {}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "not enough return values\n\thave: ()\n\twant: (bool)")
}

func TestTooManyReturnType(t *testing.T) {
	stream, _ := Tokenize("func f() {\nx := 1\nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "too many return values\n\thave: (int)\n\twant: ()")
}

func TestDiffentReturnType(t *testing.T) {
	stream, _ := Tokenize("func f() bool {\nx := 1\nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "cannot use int as bool in return statement")
}

func TestDifferentTypeAdd(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nreturn 1 + true\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "invalid operation: adding different types")
}

func TestAddingNil(t *testing.T) {
	// `returnType` node of `f` is nil because of its declaration.
	stream, _ := Tokenize("func main() int {\nreturn 1 + f(1, 2)\n}\nfunc f(x int, y int) {\nz := x + y\nreturn z\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "invalid operation: adding different types")
}
