package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddOp(t *testing.T) {
	stream, _ := Tokenize("1+2")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.root,
		&AddOp{
			tok: &Token{kind: TOKEN_PLUS, value: "+"},
			lhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "1"}},
			rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "2"}},
		},
	)
}

func TestParse2AddOp(t *testing.T) {
	stream, _ := Tokenize("1+2+3")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.root,
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
