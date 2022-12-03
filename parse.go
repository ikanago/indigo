package main

import "fmt"

type Ast struct {
	funcs []Expr
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

func newLocalEnv() *LocalEnv {
	return &LocalEnv{variables: map[string]int{}, totalOffset: 0}
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
		localEnv:    nil,
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
	var funcs []Expr
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

func (parser *parser) topLevelDecl() (Expr, error) {
	token := parser.peek()
	if token.kind == TOKEN_FUNC {
		return parser.functionDecl()
	} else {
		fmt.Printf("%v", token)
		return nil, fmt.Errorf("syntax error: non-declaration statement outside function body")
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

func (parser *parser) functionDecl() (Expr, error) {
	tokenFunc, _ := parser.expectString("func")

	token := parser.peek()
	if token.kind != TOKEN_IDENTIFIER {
		err := fmt.Errorf("unexpected %s, expecting name", token.value)
		return nil, err
	}
	name := token.value
	parser.skip()

	if err := parser.consumeString("("); err != nil {
		return nil, err
	}

	if err := parser.consumeString(")"); err != nil {
		return nil, err
	}

	body, err := parser.block()
	if err != nil {
		return nil, err
	}

	return &FunctionDecl{tok: tokenFunc, name: name, body: body}, nil
}

func (parser *parser) block() (Expr, error) {
	lbraceToken, err := parser.expectString("{")
	if err != nil {
		return nil, err
	}

	parser.localEnv = newLocalEnv()
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

	return &Block{tok: lbraceToken, body: body, localEnv: parser.localEnv}, nil
}

func (parser *parser) shortVarDecl(lhs Expr) (Expr, error) {
	if _, ok := lhs.(*Variable); ok {
		lhs = parser.newVariable(lhs.(*Variable))
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
		if token.value == "true" {
			return &BoolLiteral{tok: token, value: true}, nil
		} else if token.value == "false" {
			return &BoolLiteral{tok: token, value: false}, nil
		}

		if offset, ok := parser.localEnv.variables[token.value]; ok {
			return &Identifier{tok: token, offset: offset}, nil
		}
		return &Variable{tok: token}, nil
	default:
		err := fmt.Errorf("unexpected %s, expecting primary expression", token.value)
		return nil, err
	}
}
