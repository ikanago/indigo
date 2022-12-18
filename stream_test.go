package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteStreamNextWithoutNewLine(t *testing.T) {
	source := "abc"
	stream := NewByteStream(source)

	b, ok := stream.NextByte()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 1, Column: 1}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('b'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 1, Column: 2}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('c'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 1, Column: 3}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte(0), b)
	assert.False(t, ok)
	assert.Equal(t, &Position{Line: 1, Column: 3}, stream.currentPosition)
}

func TestByteStreamNextWithNewLine(t *testing.T) {
	source := "a\nbc\nd"
	stream := NewByteStream(source)

	b, ok := stream.NextByte()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 1, Column: 1}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 2, Column: 0}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('b'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 2, Column: 1}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('c'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 2, Column: 2}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 3, Column: 0}, stream.currentPosition)

	b, ok = stream.NextByte()
	assert.Equal(t, byte('d'), b)
	assert.True(t, ok)
	assert.Equal(t, &Position{Line: 3, Column: 1}, stream.currentPosition)
}
