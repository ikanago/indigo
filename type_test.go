package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeResolveVariable(t *testing.T) {
	stream := NewByteStream("func main() int {\nx := 1\nreturn x\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].Body.Body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.Lhs.(*Variable).Ty)
	ret := ast.funcs[0].Body.Body[1].(*Return)
	assert.Equal(t, &TypeInt, ret.Node.(*Identifier).Variable.Ty)
}

func TestTypeResolveAdd(t *testing.T) {
	stream := NewByteStream("func main() int {\nx := 1 \nreturn x\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].Body.Body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.Lhs.(*Variable).Ty)
}

func TestTypeCallFunctionWithoutArgument(t *testing.T) {
	stream := NewByteStream("func main() int {\nx := f() + 1 \nreturn x\n}\n func f() int {\nreturn 2\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].Body.Body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.Lhs.(*Variable).Ty)
	fuctionCall := assign.Rhs.(*AddOp).Lhs.(*FunctionCall)
	assert.Equal(t, &TypeInt, fuctionCall.Function.ReturnType)

	x, ok := ast.funcs[0].Scope.GetExpr("x")
	assert.True(t, ok)
	assert.Equal(t, 0, x.(*Variable).Offset)
}

func TestTypeCallFunctionWithArgument(t *testing.T) {
	stream := NewByteStream("func main() int {\nx := f(1)\nreturn x\n}\nfunc f(a int) int {\nreturn a\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	a, ok := ast.funcs[1].Scope.GetExpr("a")
	assert.True(t, ok)
	assert.Equal(t, 0, a.(*Variable).Offset)
}

func TestReturnUndefinedVariable(t *testing.T) {
	stream := NewByteStream("func main() int {\nreturn abc\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "2:8: undefined: abc")
}

func TestNotEnoughReturnType(t *testing.T) {
	stream := NewByteStream("func f() bool {}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "1:1: not enough return values\n\thave: ()\n\twant: (bool)")
}

func TestTooManyReturnType(t *testing.T) {
	stream := NewByteStream("func f() {\nx := 1\nreturn x\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "1:1: too many return values\n\thave: (int)\n\twant: ()")
}

func TestDiffentReturnType(t *testing.T) {
	stream := NewByteStream("func f() bool {\nx := 1\nreturn x\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "1:1: cannot use int as bool in return statement")
}

func TestDifferentTypeAdd(t *testing.T) {
	stream := NewByteStream("func main() int {\nreturn 1 + true\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "2:10: invalid operation: adding different types")
}

func TestAddingNil(t *testing.T) {
	// `returnType` node of `f` is nil because of its declaration.
	stream := NewByteStream("func main() int {\nreturn 1 + f(1, 2)\n}\nfunc f(x int, y int) {\nz := x + y\nreturn z\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "2:10: invalid operation: adding different types")
}
