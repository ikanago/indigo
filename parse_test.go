package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortVarDecl(t *testing.T) {
	stream, _ := Tokenize("xy := 1\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&Assign{
			tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
			lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
			rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
		},
		ast.nodes[0],
	)
}

func TestShortVarDeclAndAdd(t *testing.T) {
	stream, _ := Tokenize("xy := 1 + 2\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&Assign{
			tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
			lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
			rhs: &AddOp{
				tok: &Token{kind: TOKEN_PLUS, value: "+"},
				lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
				rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
			},
		},
		ast.nodes[0],
	)
}

func TestShortVarDeclAndReturn(t *testing.T) {
	stream, _ := Tokenize("xy := 1 + 2\nreturn xy\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]Expr{
			&Assign{
				tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
				lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
				rhs: &AddOp{
					tok: &Token{kind: TOKEN_PLUS, value: "+"},
					lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
					rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
				},
			},
			&Return{
				tok: &Token{kind: TOKEN_RETURN, value: "return"},
				node: &Identifier{
					tok:    &Token{kind: TOKEN_IDENTIFIER, value: "xy"},
					offset: 0,
				},
			},
		},
		ast.nodes,
	)
}

func TestParseAddOp(t *testing.T) {
	stream, _ := Tokenize("1+2\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&AddOp{
			tok: &Token{kind: TOKEN_PLUS, value: "+"},
			lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
			rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
		},
		ast.nodes[0],
	)
}

func TestParse2AddOp(t *testing.T) {
	stream, _ := Tokenize("1+2+3\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&AddOp{
			tok: &Token{kind: TOKEN_PLUS, value: "+"},
			lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
			rhs: &AddOp{
				tok: &Token{kind: TOKEN_PLUS, value: "+"},
				lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
				rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "3"}},
			},
		},
		ast.nodes[0],
	)
}
