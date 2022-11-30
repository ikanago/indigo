package main

import "fmt"

type Ast struct {
	nodes    []Expr
	localEnv *LocalEnv
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

func (parser *parser) consumeString(expected string) error {
	token := parser.peek()
	if token.value != expected {
		return fmt.Errorf("expected %s, but got %s", expected, token.value)
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
	return &Ast{nodes: nodes, localEnv: parser.localEnv}, nil
}

func (parser *parser) stmt() (Expr, error) {
	node, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	token := parser.peek()
	switch token.kind {
	case TOKEN_COLONEQUAL:
		return parser.shortVarDecl(node)
	default:
		return node, nil
	}
}

func (parser *parser) shortVarDecl(lhs Expr) (Expr, error) {
	switch lhs.(type) {
	case *Variable:
		lhs = parser.newVariable(lhs.(*Variable))
	default:
		err := fmt.Errorf("expected variable, but got %s", lhs.token().value)
		return nil, err
	}

	if err := parser.consumeString(":="); err != nil {
		return nil, err
	}
	rhs, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	return &Assign{tok: &Token{kind: TOKEN_COLONEQUAL, value: ":="}, lhs: lhs, rhs: rhs}, nil
}

func (parser *parser) newVariable(identifier *Variable) Expr {
	offset := parser.localEnv.insertVariable(identifier.token().value)
	return &Variable{tok: identifier.token(), offset: offset}
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
		if offset, ok := parser.localEnv.variables[token.value]; ok {
			return &Identifier{tok: token, offset: offset}, nil
		}
		return &Variable{tok: token}, nil
	default:
		err := fmt.Errorf("expected int literal, but got %s", token.value)
		return nil, err
	}
}
