package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestTokenizeMultipleLine(t *testing.T) {
	stream := NewByteStream("42\n43")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(
		[]Token{
			{Kind: TOKEN_INT, Value: "42"},
			{Kind: TOKEN_SEMICOLON, Value: ";"},
			{Kind: TOKEN_INT, Value: "43"},
			{Kind: TOKEN_EOF, Value: ""},
		},
		tokenStream.tokens,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestTokenizeShortVarDecl(t *testing.T) {
	stream := NewByteStream("xy := 42")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(
		[]Token{
			{Kind: TOKEN_IDENTIFIER, Value: "xy"},
			{Kind: TOKEN_COLONEQUAL, Value: ":="},
			{Kind: TOKEN_INT, Value: "42"},
			{Kind: TOKEN_EOF, Value: ""},
		},
		tokenStream.tokens,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestTokenizeInt(t *testing.T) {
	stream := NewByteStream("42")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(
		[]Token{
			{Kind: TOKEN_INT, Value: "42"},
			{Kind: TOKEN_EOF, Value: ""},
		},
		tokenStream.tokens,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestTokenizeSymbols(t *testing.T) {
	stream := NewByteStream("+")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(Token{Kind: TOKEN_PLUS, Value: "+"}, tokenStream.tokens[0], opts...)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestTokenizeKeywords(t *testing.T) {
	stream := NewByteStream("func main(){\nreturn abc\n}\n")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(
		[]Token{
			{Kind: TOKEN_FUNC, Value: "func"},
			{Kind: TOKEN_IDENTIFIER, Value: "main"},
			{Kind: TOKEN_LPAREN, Value: "("},
			{Kind: TOKEN_RPAREN, Value: ")"},
			{Kind: TOKEN_LBRACE, Value: "{"},
			{Kind: TOKEN_RETURN, Value: "return"},
			{Kind: TOKEN_IDENTIFIER, Value: "abc"},
			{Kind: TOKEN_SEMICOLON, Value: ";"},
			{Kind: TOKEN_RBRACE, Value: "}"},
			{Kind: TOKEN_SEMICOLON, Value: ";"},
			{Kind: TOKEN_EOF, Value: ""},
		},
		tokenStream.tokens,
		opts...,
	)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}

func TestTokenizeIdentifier(t *testing.T) {
	stream := NewByteStream("abc12")
	tokenStream, err := Tokenize(stream)
	assert.NoError(t, err)
	d := cmp.Diff(Token{Kind: TOKEN_IDENTIFIER, Value: "abc12"}, tokenStream.tokens[0], opts...)
	if len(d) != 0 {
		t.Errorf("(-got +want)\n%s", d)
	}
}
