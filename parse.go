package main

import (
	"errors"
	"fmt"
)

type Ast struct {
	funcs []*FunctionDecl
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
	localScope  *Scope
	globalScope *Scope
}

func makeParser(tokenStream *TokenStream) *parser {
	return &parser{
		tokenStream: tokenStream,
		localScope:  nil,
		globalScope: newGlobalScope(),
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

func (parser *parser) expectString(expected string) (*Token, error) {
	token := parser.peek()
	if token.value != expected {
		return nil, fmt.Errorf("unexpected %s, expecting %s", token.value, expected)
	}
	parser.skip()
	return token, nil
}

func (parser *parser) consumeString(expected string) error {
	_, err := parser.expectString(expected)
	return err
}

func (parser *parser) parse() (*Ast, error) {
	var funcs []*FunctionDecl
	for {
		if parser.peek().kind == TOKEN_EOF {
			break
		}

		f, err := parser.topLevelDecl()
		if err != nil {
			return nil, err
		}
		if err := parser.consumeString(";"); err != nil {
			return nil, err
		}
		funcs = append(funcs, f)
	}
	return &Ast{funcs}, nil
}

func (parser *parser) topLevelDecl() (*FunctionDecl, error) {
	token := parser.peek()
	if token.kind == TOKEN_FUNC {
		return parser.functionDecl()
	} else {
		return nil, errors.New("syntax error: non-declaration statement outside function body")
	}
}

func (parser *parser) stmt() (Expr, error) {
	token := parser.peek()
	if token.kind == TOKEN_RETURN {
		parser.skip()
		node, err := parser.addOp()
		if err != nil {
			return nil, err
		}
		return &Return{tok: token, node: node}, nil
	}

	node, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	token = parser.peek()
	switch token.kind {
	case TOKEN_COLONEQUAL:
		return parser.shortVarDecl(node)
	default:
		return node, nil
	}
}

func (parser *parser) functionDecl() (*FunctionDecl, error) {
	tokenFunc, _ := parser.expectString("func")

	token := parser.peek()
	if token.kind != TOKEN_IDENTIFIER {
		err := fmt.Errorf("unexpected %s, expecting name", token.value)
		return nil, err
	}
	name := token.value
	parser.skip()

	returnType, err := parser.signiture()
	if err != nil {
		return nil, err
	}

	body, err := parser.block()
	if err != nil {
		return nil, err
	}

	return &FunctionDecl{tok: tokenFunc, name: name, returnType: returnType, body: body, scope: parser.localScope}, nil
}

func (parser *parser) signiture() (*Type, error) {
	if err := parser.consumeString("("); err != nil {
		return nil, err
	}

	if err := parser.consumeString(")"); err != nil {
		return nil, err
	}

	return parser.parseType()
}

func (parser *parser) parseType() (*Type, error) {
	token := parser.peek()
	switch token.kind {
	case TOKEN_LBRACE:
		return nil, nil
	case TOKEN_IDENTIFIER:
		parser.skip()
		// Assume all types are defined so far.
		ty, _ := parser.globalScope.GetType(token.value)
		return ty, nil
	}
	return nil, errors.New("expecting type")
}

func (parser *parser) block() (*Block, error) {
	lbraceToken, err := parser.expectString("{")
	if err != nil {
		return nil, err
	}

	parser.localScope = newScope()
	var body []Expr
	for {
		if parser.peek().kind == TOKEN_RBRACE {
			parser.skip()
			break
		}
		if parser.peek().kind == TOKEN_EOF {
			return nil, fmt.Errorf("unexpected EOF, expecting }")
		}

		node, err := parser.stmt()
		if err != nil {
			return nil, err
		}
		body = append(body, node)
		if err := parser.consumeString(";"); err != nil {
			return nil, err
		}
	}

	return &Block{tok: lbraceToken, body: body}, nil
}

func (parser *parser) shortVarDecl(lhs Expr) (Expr, error) {
	if _, ok := lhs.(*Identifier); ok {
		if !parser.localScope.ExistsVariable(lhs.token().value) {
			lhs = &Variable{tok: lhs.token(), ty: &TypeUnresolved}
			parser.localScope.InsertVariable(lhs.token().value, lhs.(*Variable))
		} else {
			return nil, errors.New("no new variables on left side of :=")
		}
	} else {
		err := fmt.Errorf("unexpected %s, expecting variable", lhs.token().value)
		return nil, err
	}

	parser.consumeString(":=")
	rhs, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	return &Assign{tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="}, lhs: lhs, rhs: rhs}, nil
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
		if token.value == "true" {
			return &BoolLiteral{tok: token, value: true}, nil
		} else if token.value == "false" {
			return &BoolLiteral{tok: token, value: false}, nil
		}

		return &Identifier{tok: token}, nil
	}
	err := fmt.Errorf("unexpected %s, expecting primary expression", token.value)
	return nil, err
}
