package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeInt(t *testing.T) {
	stream, err := Tokenize("42")
	assert.NoError(t, err)
	assert.Equal(t, stream.tokens[0], Token{kind: TOKEN_INT, value: "42"})
	assert.Equal(t, stream.tokens[1], Token{kind: TOKEN_EOF})
}

func TestUnknown(t *testing.T) {
	_, err := Tokenize("42a")
	assert.Error(t, err)
}
