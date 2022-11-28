package main

import "fmt"

type Ast struct {
	root Expr
}

func Parse(tokenStream *TokenStream) (*Ast, error) {
	parser := makeParser(tokenStream)
	ast, err := parser.parse()
	if err != nil {
		return nil, err
	}
	return ast, nil
}

type parser struct {
	tokenStream *TokenStream
}

func makeParser(tokenStream *TokenStream) *parser {
	return &parser{tokenStream: tokenStream}
}

func (parser *parser) peek() *Token {
	current := parser.tokenStream.index
	return &parser.tokenStream.tokens[current]
}

func (parser *parser) skip() {
	if !parser.tokenStream.IsEnd() {
		parser.tokenStream.index += 1
	}
}

func (parser *parser) parse() (*Ast, error) {
	root, err := parser.addOp()
	if err != nil {
		return nil, err
	}
	return &Ast{root: root}, nil
}

func (parser *parser) addOp() (Expr, error) {
	lhs, err := parser.intLiteral()
	if err != nil {
		return nil, err
	}

	token := parser.peek()
	if token.kind == TOKEN_PLUS {
		parser.skip()
		rhs, err := parser.addOp()
		if err != nil {
			return nil, err
		}
		return &AddOp{tok: token, lhs: lhs, rhs: rhs}, nil
	} else {
		return lhs, nil
	}
}

func (parser *parser) intLiteral() (Expr, error) {
	token := parser.peek()
	switch token.kind {
	case TOKEN_INT:
		parser.skip()
		return &IntLiteral{tok: token}, nil
	default:
		err := fmt.Errorf("Expected int literal, but got %s", token.value)
		return nil, err
	}
}
