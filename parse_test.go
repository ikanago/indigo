package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

var opts = []cmp.Option{
	cmpopts.IgnoreUnexported(Token{}),
	cmpopts.IgnoreUnexported(FunctionDecl{}),
	cmpopts.IgnoreUnexported(Block{}),
	cmpopts.IgnoreUnexported(Return{}),
	cmpopts.IgnoreUnexported(Assign{}),
	cmpopts.IgnoreUnexported(AddOp{}),
	cmpopts.IgnoreUnexported(Variable{}),
	cmpopts.IgnoreUnexported(Identifier{}),
	cmpopts.IgnoreUnexported(IntLiteral{}),
	cmpopts.IgnoreUnexported(BoolLiteral{}),
	cmpopts.IgnoreUnexported(FunctionCall{}),
}

func TestFuncDef(t *testing.T) {
	stream := NewByteStream("func main(){\nabc := 3\nreturn abc\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)

	d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Assign{
					Lhs: &Variable{Name: "abc", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &IntLiteral{Value: "3"},
				},
				&Return{Node: &Identifier{Name: "abc", Variable: nil}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestCallFunctionWithoutArgument(t *testing.T) {
	stream := NewByteStream("func main() int{\n x := f()\nreturn x\n}\nfunc f() int {\nreturn 3\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	assert.Equal(t, &TypeInt, ast.funcs[0].ReturnType)
	d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Assign{
					Lhs: &Variable{Name: "x", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &FunctionCall{Arguments: []Expr{}},
				},
				&Return{Node: &Identifier{Name: "x"}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}

	assert.Equal(t, &TypeInt, ast.funcs[1].ReturnType)
	d = cmp.Diff(
		&Block{
			Body: []Expr{
				&Return{Node: &IntLiteral{Value: "3"}},
			},
		},
		ast.funcs[1].Body,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestFuncReturnType(t *testing.T) {
	stream := NewByteStream("func f() int {\nreturn 3\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	assert.Equal(t, &TypeInt, ast.funcs[0].ReturnType)
	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Return{Node: &IntLiteral{Value: "3"}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestFunctionWithOneArgument(t *testing.T) {
	stream := NewByteStream("func f(a int) int {\nreturn a\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	assert.Equal(t, &TypeInt, ast.funcs[0].ReturnType)
	if d := cmp.Diff(
		[]*Variable{
			{Name: "a", Offset: 0, Ty: &TypeInt},
		},
		ast.funcs[0].Parameters,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}

	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Return{Node: &Identifier{Name: "a"}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestFunctionWithArguments(t *testing.T) {
	stream := NewByteStream("func f(a int, b int) int {\nreturn a + b\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	assert.Equal(t, &TypeInt, ast.funcs[0].ReturnType)
	if d := cmp.Diff(
		[]*Variable{
			{Name: "a", Offset: 0, Ty: &TypeInt},
			{Name: "b", Offset: 0, Ty: &TypeInt},
		},
		ast.funcs[0].Parameters,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}

	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Return{
					Node: &AddOp{
						Lhs: &Identifier{Name: "a"},
						Rhs: &Identifier{Name: "b"},
					},
				},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestCallFunctionWithArgument(t *testing.T) {
	stream := NewByteStream("func main() {\nx := f(1)\nreturn x\n}\nfunc f(a int) int {\nreturn a\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Assign{
					Lhs: &Variable{Name: "x", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &FunctionCall{
						Arguments: []Expr{
							&IntLiteral{Value: "1"},
						},
					},
				},
				&Return{Node: &Identifier{Name: "x"}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestCallFunctionWithArguments(t *testing.T) {
	stream := NewByteStream("func main() {\nx := 1\ny := f(x, 2 + 3)\nreturn y\n}\nfunc f(a int, b int) int {\nreturn a + b\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Assign{
					Lhs: &Variable{Name: "x", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &IntLiteral{Value: "1"},
				},
				&Assign{
					Lhs: &Variable{Name: "y", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &FunctionCall{
						Arguments: []Expr{
							&Identifier{Name: "x"},
							&AddOp{
								Lhs: &IntLiteral{Value: "2"},
								Rhs: &IntLiteral{Value: "3"},
							},
						},
					},
				},
				&Return{
					Node: &Identifier{Name: "y"},
				},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestBool(t *testing.T) {
	stream := NewByteStream("func main(){\nreturn true\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)
	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Return{
					Node: &BoolLiteral{Value: true},
				},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestShortVarDeclAndAdd(t *testing.T) {
	stream := NewByteStream("func main(){\nxy := 1 + 2 + 3\nreturn xy\n}\n")
	tokenStream, _ := Tokenize(stream)
	ast, err := Parse(tokenStream)
	assert.NoError(t, err)

	if d := cmp.Diff(
		&Block{
			Body: []Expr{
				&Assign{
					Lhs: &Variable{Name: "xy", Offset: 0, Ty: &TypeUnresolved},
					Rhs: &AddOp{
						Lhs: &IntLiteral{Value: "1"},
						Rhs: &AddOp{
							Lhs: &IntLiteral{Value: "2"},
							Rhs: &IntLiteral{Value: "3"},
						},
					},
				},
				&Return{Node: &Identifier{Name: "xy"}},
			},
		},
		ast.funcs[0].Body,
		opts...,
	); len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestLhsOfShortVarDeclIsNotIdentifier(t *testing.T) {
	stream := NewByteStream("func main(){\n1 := 2\n}\n")
	tokenStream, _ := Tokenize(stream)
	_, err := Parse(tokenStream)
	assert.Error(t, err)
}

func TestNoNewVariableOnRhsOfShortVarDecl(t *testing.T) {
	stream := NewByteStream("func main(){\nx := 1\nx := 2\n}\n")
	tokenStream, _ := Tokenize(stream)
	_, err := Parse(tokenStream)
	assert.EqualError(t, err, "no new variables on left side of :=")
}
