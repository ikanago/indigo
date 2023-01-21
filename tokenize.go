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
	Kind  TokenKind
	Value string
	pos   Position
}

type TokenStream struct {
	tokens []Token
	index  int
}

func (stream *TokenStream) IsEnd() bool {
	token := stream.tokens[stream.index]
	return token.Kind == TOKEN_EOF
}

// func mergeToken(token *Token, other *Token) *Token {}

func initKeywordMap() map[string]TokenKind {
	return map[string]TokenKind{
		"func":   TOKEN_FUNC,
		"return": TOKEN_RETURN,
	}
}

func Tokenize(stream *ByteStream) (*TokenStream, error) {
	keywordMap := initKeywordMap()
	var tokens []Token
	for {
		currentByte, ok := stream.get()
		if !ok {
			break
		}
		pos := stream.CurrentPosition

		if isWhitespace(currentByte) {
			continue
		} else if currentByte == '\n' {
			if shouldInsertSemicolon(tokens) {
				tokens = append(tokens, Token{Kind: TOKEN_SEMICOLON, Value: ";"})
			}
		} else if currentByte == '(' {
			token := Token{
				Kind:  TOKEN_LPAREN,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if currentByte == ')' {
			token := Token{
				Kind:  TOKEN_RPAREN,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if currentByte == '{' {
			token := Token{
				Kind:  TOKEN_LBRACE,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if currentByte == '}' {
			token := Token{
				Kind:  TOKEN_RBRACE,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if currentByte == '+' {
			token := Token{
				Kind:  TOKEN_PLUS,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if currentByte == ':' {
			c, ok := stream.get()
			if ok && c == '=' {
				token := Token{
					Kind:  TOKEN_COLONEQUAL,
					Value: ":=",
					pos:   pos,
				}
				tokens = append(tokens, token)
			} else {
				err := fmt.Errorf("unknown character: %c", currentByte)
				return nil, err
			}
		} else if currentByte == ',' {
			token := Token{
				Kind:  TOKEN_COMMA,
				Value: string(currentByte),
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if isDigit(currentByte) {
			digits := stream.readDigit(currentByte)
			token := Token{
				Kind:  TOKEN_INT,
				Value: digits,
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else if isLetter(currentByte) {
			identifier := stream.readIdentifier(currentByte)
			var kind TokenKind
			if kindKeyword, ok := keywordMap[identifier]; ok {
				kind = kindKeyword
			} else {
				kind = TOKEN_IDENTIFIER
			}
			token := Token{
				Kind:  kind,
				Value: identifier,
				pos:   pos,
			}
			tokens = append(tokens, token)
		} else {
			err := fmt.Errorf("unknown character: %c", currentByte)
			return nil, err
		}
	}

	tokens = append(tokens, Token{Kind: TOKEN_EOF})

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

	switch tokens[len(tokens)-1].Kind {
	case TOKEN_IDENTIFIER, TOKEN_INT, TOKEN_RPAREN, TOKEN_RBRACE:
		return true
	default:
		return false
	}
}

func (stream *ByteStream) readDigit(c0 byte) string {
	digits := []byte{c0}
	for {
		c, ok := stream.get()
		if !ok {
			break
		}

		if !isDigit(c) {
			stream.unget()
			break
		}

		digits = append(digits, c)
	}
	return string(digits)
}

func (stream *ByteStream) readIdentifier(c0 byte) string {
	identifier := []byte{c0}
	for {
		c, ok := stream.get()
		if !ok {
			break
		}

		if !isLetter(c) && !isDigit(c) {
			stream.unget()
			break
		}

		identifier = append(identifier, c)
	}
	return string(identifier)
}
