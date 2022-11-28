package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeMultipleLine(t *testing.T) {
	stream, err := Tokenize("42\n43")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens, []Token{
		{kind: TOKEN_INT, value: "42"},
		{kind: TOKEN_SEMICOLON, value: ";"},
		{kind: TOKEN_INT, value: "43"},
		{kind: TOKEN_EOF, value: ""},
	})
}

func TestTokenizeShortVarDecl(t *testing.T) {
	stream, err := Tokenize("xy := 42")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens, []Token{
		{kind: TOKEN_IDENTIFIER, value: "xy"},
		{kind: TOKEN_COLONEQUAL, value: ":="},
		{kind: TOKEN_INT, value: "42"},
		{kind: TOKEN_EOF, value: ""},
	})
}

func TestTokenizeInt(t *testing.T) {
	stream, err := Tokenize("42")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens, []Token{
		{kind: TOKEN_INT, value: "42"},
		{kind: TOKEN_EOF, value: ""},
	})
}

func TestTokenizeSymbols(t *testing.T) {
	stream, err := Tokenize("+")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens[0], Token{kind: TOKEN_PLUS, value: "+"})
}

func TestTokenizeIdentifier(t *testing.T) {
	stream, err := Tokenize("abc12")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens[0], Token{kind: TOKEN_IDENTIFIER, value: "abc12"})
}
