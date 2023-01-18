package main

type ByteStream struct {
	source          string
	currentIndex    int
	CurrentPosition Position // Position of most recently emitted byte.
}

func NewByteStream(source string) *ByteStream {
	return &ByteStream{
		source:          source,
		currentIndex:    0,
		CurrentPosition: NewPosition(),
	}
}

// get returns the upcoming byte and whether it exists or not.
func (stream *ByteStream) get() (byte, bool) {
	if stream.currentIndex >= len(stream.source) {
		return byte(0), false
	}
	b := stream.source[stream.currentIndex]
	stream.currentIndex += 1
	if b == '\n' {
		stream.CurrentPosition = stream.CurrentPosition.nextLine()
	} else {
		stream.CurrentPosition = stream.CurrentPosition.step()
	}
	return b, true
}

func (stream *ByteStream) unget() bool {
	if stream.currentIndex == 0 {
		return false
	}
	stream.currentIndex -= 1
	b := stream.source[stream.currentIndex]
	if b == '\n' {
		stream.CurrentPosition = stream.CurrentPosition.previousLine()
	} else {
		stream.CurrentPosition = stream.CurrentPosition.back()
	}
	return true
}

type Position struct {
	Line   int // 1-indexed.
	Column int // 1-indexed.
}

func NewPosition() Position {
	return Position{Line: 1, Column: 0}
}

func (position Position) step() Position {
	return Position{
		Line:   position.Line,
		Column: position.Column + 1,
	}
}

func (position Position) stepBy(n int) Position {
	return Position{
		Line:   position.Line,
		Column: position.Column + n,
	}
}

func (position Position) back() Position {
	return Position{
		Line:   position.Line,
		Column: position.Column - 1,
	}
}

func (position Position) nextLine() Position {
	return Position{
		Line:   position.Line + 1,
		Column: 0,
	}
}

func (position Position) previousLine() Position {
	return Position{
		Line:   position.Line - 1,
		Column: 0,
	}
}

func (pos Position) isPrecedingTo(other Position) bool {
	if pos.Line == other.Line {
		return pos.Column < other.Column
	}
	return pos.Line < other.Line
}
