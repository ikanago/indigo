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
		ast.nodes[0],
		&Assign{
			tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
			lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
			rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
		},
	)
}

func TestShortVarDeclAndAdd(t *testing.T) {
	stream, _ := Tokenize("xy := 1 + 2\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.nodes[0],
		&Assign{
			tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
			lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
			rhs: &AddOp{
				tok: &Token{kind: TOKEN_PLUS, value: "+"},
				lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
				rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
			},
		},
	)
}

func TestShortVarDeclAndReturn(t *testing.T) {
	stream, _ := Tokenize("xy := 1 + 2\nxy\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.nodes,
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
			&Variable{
				tok:    &Token{kind: TOKEN_IDENTIFIER, value: "xy"},
				offset: 0,
			},
		},
	)
}

func TestParseAddOp(t *testing.T) {
	stream, _ := Tokenize("1+2\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.nodes[0],
		&AddOp{
			tok: &Token{kind: TOKEN_PLUS, value: "+"},
			lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
			rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
		},
	)
}

func TestParse2AddOp(t *testing.T) {
	stream, _ := Tokenize("1+2+3\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.nodes[0],
		&AddOp{
			tok: &Token{kind: TOKEN_PLUS, value: "+"},
			lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
			rhs: &AddOp{
				tok: &Token{kind: TOKEN_PLUS, value: "+"},
				lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
				rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "3"}},
			},
		},
	)
}
