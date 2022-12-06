package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveVariable(t *testing.T) {
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

func TestResolveAdd(t *testing.T) {
	stream, _ := Tokenize("func main() int {\nx := 1 \nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].body.body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.lhs.(*Variable).ty)
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
