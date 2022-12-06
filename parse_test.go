package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncDef(t *testing.T) {
	stream, _ := Tokenize("func main(){\nabc := 3\nreturn abc\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)

	assert.Equal(
		t,
		&Block{
			tok: &Token{kind: TOKEN_LBRACE, value: "{"},
			body: []Expr{
				&Assign{
					tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
					lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "abc"}, offset: 0, ty: &TypeUnresolved},
					rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "3"}},
				},
				&Return{
					tok: &Token{kind: TOKEN_RETURN, value: "return"},
					node: &Identifier{
						tok:      &Token{kind: TOKEN_IDENTIFIER, value: "abc"},
						variable: nil,
					},
				},
			},
		},
		ast.funcs[0].body,
	)
}

func TestFuncReturnType(t *testing.T) {
	stream, _ := Tokenize("func f() int {\nreturn 3\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(t, &TypeInt, ast.funcs[0].returnType)
	assert.Equal(
		t,
		&Block{
			tok: &Token{kind: TOKEN_LBRACE, value: "{"},
			body: []Expr{
				&Return{
					tok: &Token{kind: TOKEN_RETURN, value: "return"},
					node: &IntLiteral{
						tok: &Token{kind: TOKEN_INT, value: "3"},
					},
				},
			},
		},
		ast.funcs[0].body,
	)
}

func TestBool(t *testing.T) {
	stream, _ := Tokenize("func main(){\nreturn true\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&Block{
			tok: &Token{kind: TOKEN_LBRACE, value: "{"},
			body: []Expr{
				&Return{
					tok: &Token{kind: TOKEN_RETURN, value: "return"},
					node: &BoolLiteral{
						tok:   &Token{kind: TOKEN_IDENTIFIER, value: "true"},
						value: true,
					},
				},
			},
		},
		ast.funcs[0].body,
	)
}

func TestShortVarDeclAndAdd(t *testing.T) {
	stream, _ := Tokenize("func main(){\nxy := 1 + 2 + 3\nreturn xy\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)

	assert.Equal(
		t,
		&Block{
			tok: &Token{kind: TOKEN_LBRACE, value: "{"},
			body: []Expr{
				&Assign{
					tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
					lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0, ty: &TypeUnresolved},
					rhs: &AddOp{
						tok: &Token{kind: TOKEN_PLUS, value: "+"},
						lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
						rhs: &AddOp{
							tok: &Token{kind: TOKEN_PLUS, value: "+"},
							lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
							rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "3"}},
						},
					},
				},
				&Return{
					tok: &Token{kind: TOKEN_RETURN, value: "return"},
					node: &Identifier{
						tok:      &Token{kind: TOKEN_IDENTIFIER, value: "xy"},
						variable: nil,
					},
				},
			},
		},
		ast.funcs[0].body,
	)
}

func TestLhsOfShortVarDeclIsNotIdentifier(t *testing.T) {
	stream, _ := Tokenize("func main(){\n1 := 2\n}\n")
	_, err := Parse(stream)
	assert.Error(t, err)
}

func TestNoNewVariableOnRhsOfShortVarDecl(t *testing.T) {
	stream, _ := Tokenize("func main(){\nx := 1\nx := 2\n}\n")
	_, err := Parse(stream)
	assert.EqualError(t, err, "no new variables on left side of :=")
}
