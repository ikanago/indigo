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
		&FunctionDecl{
			tok:  &Token{kind: TOKEN_FUNC, value: "func"},
			name: "main",
			body: &Block{
				tok: &Token{kind: TOKEN_LBRACE, value: "{"},
				body: []Expr{
					&Assign{
						tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
						lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "abc"}, offset: 0},
						rhs: &IntLiteral{tok: &Token{kind: TOKEN_INT, value: "3"}},
					},
					&Return{
						tok: &Token{kind: TOKEN_RETURN, value: "return"},
						node: &Identifier{
							tok:    &Token{kind: TOKEN_IDENTIFIER, value: "abc"},
							offset: 0,
						},
					},
				},
				localEnv: &LocalEnv{variables: map[string]int{"abc": 0}, totalOffset: 16},
			},
		},
		ast.funcs[0],
	)
}

func TestBool(t *testing.T) {
	stream, _ := Tokenize("func main(){\nreturn true\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		ast.funcs[0],
		&FunctionDecl{
			tok:  &Token{kind: TOKEN_FUNC, value: "func"},
			name: "main",
			body: &Block{
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
				localEnv: &LocalEnv{variables: map[string]int{}, totalOffset: 0},
			},
		},
	)
}

func TestShortVarDeclAndAdd(t *testing.T) {
	stream, _ := Tokenize("func main(){\nxy := 1 + 2 + 3\nreturn xy\n}\n")
	ast, err := Parse(stream)
	assert.NoError(t, err)
	assert.Equal(
		t,
		&FunctionDecl{
			tok:  &Token{kind: TOKEN_FUNC, value: "func"},
			name: "main",
			body: &Block{
				tok: &Token{kind: TOKEN_LBRACE, value: "{"},
				body: []Expr{
					&Assign{
						tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="},
						lhs: &Variable{tok: &Token{kind: TOKEN_IDENTIFIER, value: "xy"}, offset: 0},
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
							tok:    &Token{kind: TOKEN_IDENTIFIER, value: "xy"},
							offset: 0,
						},
					},
				},
				localEnv: &LocalEnv{variables: map[string]int{"xy": 0}, totalOffset: 16},
			},
		},
		ast.funcs[0],
	)
}
