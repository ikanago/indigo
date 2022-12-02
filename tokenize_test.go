package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeMultipleLine(t *testing.T) {
	stream, err := Tokenize("42\n43")
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]Token{
			{kind: TOKEN_INT, value: "42"},
			{kind: TOKEN_SEMICOLON, value: ";"},
			{kind: TOKEN_INT, value: "43"},
			{kind: TOKEN_EOF, value: ""},
		},
		stream.tokens,
	)
}

func TestTokenizeShortVarDecl(t *testing.T) {
	stream, err := Tokenize("xy := 42")
	assert.NoError(t, err)
	assert.Equal(t,
		[]Token{
			{kind: TOKEN_IDENTIFIER, value: "xy"},
			{kind: TOKEN_COLONEQUAL, value: ":="},
			{kind: TOKEN_INT, value: "42"},
			{kind: TOKEN_EOF, value: ""},
		},
		stream.tokens,
	)
}

func TestTokenizeInt(t *testing.T) {
	stream, err := Tokenize("42")
	assert.NoError(t, err)
	assert.Equal(t,
		[]Token{
			{kind: TOKEN_INT, value: "42"},
			{kind: TOKEN_EOF, value: ""},
		},
		stream.tokens,
	)
}

func TestTokenizeSymbols(t *testing.T) {
	stream, err := Tokenize("+")
	assert.NoError(t, err)
	assert.Equal(t, Token{kind: TOKEN_PLUS, value: "+"}, stream.tokens[0])
}

func TestTokenizeKeywords(t *testing.T) {
	stream, err := Tokenize("return abc")
	assert.NoError(t, err)
	assert.Equal(t,
		[]Token{
			{kind: TOKEN_RETURN, value: "return"},
			{kind: TOKEN_IDENTIFIER, value: "abc"},
			{kind: TOKEN_EOF, value: ""},
		},
		stream.tokens,
	)
}

func TestTokenizeIdentifier(t *testing.T) {
	stream, err := Tokenize("abc12")
	assert.NoError(t, err)
	assert.Equal(t, Token{kind: TOKEN_IDENTIFIER, value: "abc12"}, stream.tokens[0])
}
