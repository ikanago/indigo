package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteStreamNextWithoutNewLine(t *testing.T) {
	source := "abc"
	stream := NewByteStream(source)

	b, ok := stream.get()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 1}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('b'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 2}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('c'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 3}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte(0), b)
	assert.False(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 3}, stream.CurrentPosition)
}

func TestByteStreamNextWithNewLine(t *testing.T) {
	source := "a\nbc\nd"
	stream := NewByteStream(source)

	b, ok := stream.get()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 1}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 2, Column: 0}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('b'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 2, Column: 1}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('c'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 2, Column: 2}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 3, Column: 0}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('d'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 3, Column: 1}, stream.CurrentPosition)
}

func TestUnget(t *testing.T) {
	stream := NewByteStream("abc")
	b, ok := stream.get()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 1}, stream.CurrentPosition)

	ok = stream.unget()
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 0}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('a'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 1}, stream.CurrentPosition)
}

func TestUngetAroundNewLine(t *testing.T) {
	stream := NewByteStream("a\nb")
	stream.get()
	b, ok := stream.get()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 2, Column: 0}, stream.CurrentPosition)

	ok = stream.unget()
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 1, Column: 0}, stream.CurrentPosition)

	b, ok = stream.get()
	assert.Equal(t, byte('\n'), b)
	assert.True(t, ok)
	assert.Equal(t, Position{Line: 2, Column: 0}, stream.CurrentPosition)
}
