package main

type ByteStream struct {
	source          string
	currentIndex    int
	currentPosition *Position // Position of most recently emitted byte.
}

func NewByteStream(source string) *ByteStream {
	return &ByteStream{
		source:          source,
		currentIndex:    0,
		currentPosition: NewPosition(),
	}
}

// NextByte returns the upcoming byte and whether it exists or not.
func (stream *ByteStream) NextByte() (byte, bool) {
	if stream.currentIndex >= len(stream.source) {
		return byte(0), false
	}
	b := stream.source[stream.currentIndex]
	stream.currentIndex += 1
	if b == '\n' {
		stream.currentPosition.nextLine()
	} else {
		stream.currentPosition.step()
	}
	return b, true
}

type Position struct {
	Line   int // 1-indexed.
	Column int // 1-indexed.
}

func NewPosition() *Position {
	return &Position{Line: 1, Column: 0}
}

func (position *Position) step() {
	position.Column += 1
}

func (position *Position) nextLine() {
	position.Line += 1
	position.Column = 0
}
