package main

import "fmt"

type Ast struct {
	nodes []Expr
}

func Parse(tokenStream *TokenStream) (*Ast, error) {
	parser := makeParser(tokenStream)
	ast, err := parser.parse()
	if err != nil {
		return nil, err
	}
	return ast, nil
}

type LocalEnv struct {
	variables   map[string]int
	totalOffset int
}

func (env *LocalEnv) insertVariable(name string) int {
	offset := env.totalOffset
	env.variables[name] = offset
	env.totalOffset += 16
	return offset
}

type parser struct {
	tokenStream *TokenStream
	localEnv    *LocalEnv
}

func makeParser(tokenStream *TokenStream) *parser {
	return &parser{
		tokenStream: tokenStream,
		localEnv:    &LocalEnv{variables: map[string]int{}, totalOffset: 0},
	}
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

func (parser *parser) try(fn func() (Expr, error)) (Expr, bool) {
	prevTokenIndex := parser.tokenStream.index
	expr, err := fn()
	if err != nil {
		parser.tokenStream.index = prevTokenIndex
		return nil, false
	}
	return expr, true
}

func (parser *parser) consume(expected Token) error {
	token := parser.peek()
	if token.kind != expected.kind || token.value != expected.value {
		return fmt.Errorf("Expected %s, but got %s", expected.value, token.value)
	}
	return nil
}

func (parser *parser) consumeString(expected string) error {
	token := parser.peek()
	if token.value != expected {
		return fmt.Errorf("Expected %s, but got %s", expected, token.value)
	}
	parser.skip()
	return nil
}

func (parser *parser) parse() (*Ast, error) {
	var nodes []Expr
	for {
		if parser.peek().kind == TOKEN_EOF {
			break
		}
		node, _ := parser.stmt()
		nodes = append(nodes, node)
		if err := parser.consumeString(";"); err != nil {
			return nil, err
		}
	}
	return &Ast{nodes}, nil
}

func (parser *parser) stmt() (Expr, error) {
	if expr, ok := parser.try(parser.shortVarDecl); ok {
		return expr, nil
	}
	if expr, ok := parser.try(parser.addOp); ok {
		return expr, nil
	}
	return nil, nil
}

func (parser *parser) shortVarDecl() (Expr, error) {
	lhs, err := parser.newVariable()
	if err != nil {
		return nil, err
	}
	if err := parser.consumeString(":="); err != nil {
		return nil, err
	}
	rhs, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	return &ShortVarDecl{tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="}, lhs: lhs, rhs: rhs}, nil
}

func (parser *parser) newVariable() (Expr, error) {
	token := parser.peek()
	switch token.kind {
	case TOKEN_IDENTIFIER:
		parser.skip()
		offset := parser.localEnv.insertVariable(token.value)
		return &Variable{tok: token, offset: offset}, nil
	default:
		err := fmt.Errorf("Expected variable, but got %s", token.value)
		return nil, err
	}
}

func (parser *parser) addOp() (Expr, error) {
	lhs, err := parser.primaryExpr()
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

func (parser *parser) primaryExpr() (Expr, error) {
	token := parser.peek()
	switch token.kind {
	case TOKEN_INT:
		parser.skip()
		return &IntLiteral{tok: token}, nil
	case TOKEN_IDENTIFIER:
		parser.skip()
		offset, ok := parser.localEnv.variables[token.value]
		if !ok {
			err := fmt.Errorf("Variable %s is not defined", token.value)
			return nil, err
		}
		return &Variable{tok: token, offset: offset}, nil
	default:
		err := fmt.Errorf("Expected int literal, but got %s", token.value)
		return nil, err
	}
}
