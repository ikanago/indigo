package main

import "fmt"

type TokenKind int

const (
	TOKEN_INT = iota
	TOKEN_EOF
)

type Token struct {
	kind  TokenKind
	value string
}

type TokenStream struct {
	tokens []Token
	index  int
}

func (stream *TokenStream) IsEnd() bool {
	token := stream.tokens[stream.index]
	return token.kind == TOKEN_EOF
}

func Tokenize(source string) (*TokenStream, error) {
	current := 0
	var tokens []Token
	for {
		if current >= len(source) {
			break
		}

		currentByte := source[current]

		if isDigit(currentByte) {
			digits := readDigit(source, current)
			token := Token{kind: TOKEN_INT, value: digits}
			tokens = append(tokens, token)
			current += len(digits)
		} else {
			err := fmt.Errorf("Unknown character: %c", currentByte)
			return nil, err
		}
	}

	tokens = append(tokens, Token{kind: TOKEN_EOF})

	return &TokenStream{tokens: tokens, index: 0}, nil
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func readDigit(source string, startIndex int) string {
	digitsLen := 0
	for {
		if startIndex+digitsLen >= len(source) {
			break
		}

		if !isDigit(source[startIndex+digitsLen]) {
			break
		}

		digitsLen += 1
	}
	return source[startIndex:(startIndex + digitsLen)]
}
