package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveVariable(t *testing.T) {
	stream, _ := Tokenize("func main(){\nx := 1\nreturn x\n}\n")
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
	stream, _ := Tokenize("func main(){\nx := 1 \nreturn x\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.NoError(t, err)

	assign := ast.funcs[0].body.body[0].(*Assign)
	assert.Equal(t, &TypeInt, assign.lhs.(*Variable).ty)
}

func TestReturnUndefinedVariable(t *testing.T) {
	stream, _ := Tokenize("func main(){\nreturn abc\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "undefined: abc")
}

func TestDifferentTypeAdd(t *testing.T) {
	stream, _ := Tokenize("func main(){\nreturn 1 + true\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	err = ast.InferType()
	assert.EqualError(t, err, "invalid operation: adding different types")
}
