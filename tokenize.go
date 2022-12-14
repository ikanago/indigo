package main

import "fmt"

type TokenKind int

const (
	TOKEN_INT = iota
	TOKEN_IDENTIFIER
	// Symbols
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_COLONEQUAL
	TOKEN_COMMA
	// Keywords
	TOKEN_FUNC
	TOKEN_RETURN
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

func initKeywordMap() map[string]TokenKind {
	return map[string]TokenKind{
		"func":   TOKEN_FUNC,
		"return": TOKEN_RETURN,
	}
}

func Tokenize(source string) (*TokenStream, error) {
	keywordMap := initKeywordMap()
	current := 0
	var tokens []Token
	for {
		if current >= len(source) {
			break
		}

		currentByte := source[current]

		if isWhitespace(currentByte) {
			current += 1
		} else if currentByte == '\n' {
			if shouldInsertSemicolon(tokens) {
				tokens = append(tokens, Token{kind: TOKEN_SEMICOLON, value: ";"})
			}
			current += 1
		} else if currentByte == '(' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_LPAREN, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if currentByte == ')' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_RPAREN, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if currentByte == '{' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_LBRACE, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if currentByte == '}' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_RBRACE, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if currentByte == '+' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_PLUS, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if expectString(source, current, ":=") {
			token := Token{kind: TOKEN_COLONEQUAL, value: ":="}
			tokens = append(tokens, token)
			current += 2
		} else if currentByte == ',' {
			symbol := source[current:(current + 1)]
			token := Token{kind: TOKEN_COMMA, value: symbol}
			tokens = append(tokens, token)
			current += len(symbol)
		} else if isDigit(currentByte) {
			digits := readDigit(source, current)
			token := Token{kind: TOKEN_INT, value: digits}
			tokens = append(tokens, token)
			current += len(digits)
		} else if isLetter(currentByte) {
			identifier := readIdentifier(source, current)
			var token Token
			if kindKeyword, ok := keywordMap[identifier]; ok {
				token = Token{kind: kindKeyword, value: identifier}
			} else {
				token = Token{kind: TOKEN_IDENTIFIER, value: identifier}
			}
			tokens = append(tokens, token)
			current += len(identifier)
		} else {
			err := fmt.Errorf("unknown character: %c", currentByte)
			return nil, err
		}
	}

	tokens = append(tokens, Token{kind: TOKEN_EOF})

	return &TokenStream{tokens: tokens, index: 0}, nil
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

func expectString(source string, startIndex int, expected string) bool {
	remains := len(source) - 1 - startIndex
	if remains < len(expected) {
		return false
	}

	return source[startIndex:(startIndex+len(expected))] == expected
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLetter(c byte) bool {
	return isAlpha(c) || c == '_'
}

func isAlpha(c byte) bool {
	isLower := 'a' <= c && c <= 'z'
	isUpper := 'A' <= c && c <= 'Z'
	return isLower || isUpper
}

// Determine whether a semicolon should be inserted or not.
// Refer to this page for the rule: https://go.dev/ref/spec#Expression:~:text=a%20valid%20token.-,Semicolons,-The%20formal%20syntax
func shouldInsertSemicolon(tokens []Token) bool {
	if len(tokens) == 0 {
		return false
	}

	switch tokens[len(tokens)-1].kind {
	case TOKEN_IDENTIFIER, TOKEN_INT, TOKEN_RPAREN, TOKEN_RBRACE:
		return true
	default:
		return false
	}
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

func readIdentifier(source string, startIndex int) string {
	if !isLetter(source[startIndex]) {
		return ""
	}

	identifierLen := 1
	for {
		currentIndex := startIndex + identifierLen
		if currentIndex >= len(source) {
			break
		}

		if !isLetter(source[currentIndex]) && !isDigit(source[currentIndex]) {
			break
		}

		identifierLen += 1
	}

	return source[startIndex:(startIndex + identifierLen)]
}
