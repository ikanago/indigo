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
		globalScope: NewGlobalScope(),
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
	if token.Value != expected {
		return nil, fmt.Errorf("%s: unexpected %s, expecting %s", token.pos.toString(), token.Value, expected)
	}
	parser.skip()
	return token, nil
}

func (parser *parser) consumeString(expected string) error {
	_, err := parser.expectString(expected)
	return err
}

func (parser *parser) parse() (*Ast, error) {
	var functions []*FunctionDecl
	for {
		if parser.peek().Kind == TOKEN_EOF {
			break
		}

		function, err := parser.topLevelDecl()
		if err != nil {
			return nil, err
		}
		if err := parser.consumeString(";"); err != nil {
			return nil, err
		}
		functions = append(functions, function)
	}
	return &Ast{functions}, nil
}

func (parser *parser) topLevelDecl() (*FunctionDecl, error) {
	token := parser.peek()
	if token.Kind == TOKEN_FUNC {
		return parser.functionDecl()
	} else {
		return nil, fmt.Errorf("%s: syntax error: non-declaration statement outside function body: %s", token.pos.toString(), token.Value)
	}
}

func (parser *parser) stmt() (Expr, error) {
	token := parser.peek()
	if token.Kind == TOKEN_RETURN {
		parser.skip()
		node, err := parser.addOp()
		if err != nil {
			return nil, err
		}
		return &Return{tok: token, Node: node}, nil
	}

	node, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	token = parser.peek()
	switch token.Kind {
	case TOKEN_COLONEQUAL:
		return parser.shortVarDecl(node)
	default:
		return node, nil
	}
}

func (parser *parser) functionDecl() (*FunctionDecl, error) {
	parser.localScope = NewScope(parser.globalScope)

	tokenFunc, _ := parser.expectString("func")

	token := parser.peek()
	if token.Kind != TOKEN_IDENTIFIER {
		err := fmt.Errorf("%s: unexpected %s, expecting name", token.pos.toString(), token.Value)
		return nil, err
	}
	name := token.Value
	parser.skip()

	parameters, returnType, err := parser.signiture()
	if err != nil {
		return nil, err
	}

	body, err := parser.block()
	if err != nil {
		return nil, err
	}

	function := &FunctionDecl{
		tok:        tokenFunc,
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		Scope:      parser.localScope,
	}
	parser.globalScope.InsertExpr(name, function)
	return function, nil
}

func (parser *parser) signiture() ([]*Variable, *Type, error) {
	if err := parser.consumeString("("); err != nil {
		return nil, nil, err
	}

	parameters := []*Variable{}
	if parser.peek().Kind != TOKEN_RPAREN {
		if parameter, err := parser.parameterDecl(); err != nil {
			return nil, nil, err
		} else {
			parameters = append(parameters, parameter)
		}
		for {
			if parser.peek().Kind == TOKEN_RPAREN {
				break
			}
			if err := parser.consumeString(","); err != nil {
				return nil, nil, err
			}
			if parameter, err := parser.parameterDecl(); err != nil {
				return nil, nil, err
			} else {
				parameters = append(parameters, parameter)
			}
		}
	}

	if err := parser.consumeString(")"); err != nil {
		return nil, nil, err
	}

	returnType, err := parser.parseType()
	if err != nil {
		return nil, nil, err
	}
	if returnType != nil && !parser.globalScope.ExistsType(returnType.Name) {
		return nil, nil, fmt.Errorf("undefined: %s", returnType.Name)
	}
	return parameters, returnType, nil
}

func (parser *parser) parameterDecl() (*Variable, error) {
	parameterToken := parser.peek()
	if parameterToken.Kind == TOKEN_IDENTIFIER {
		parser.skip()
		ty, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		parameter := &Variable{tok: parameterToken, Name: parameterToken.Value, Ty: ty}
		name := parameterToken.Value
		if parser.localScope.ExistsExpr(name) {
			return nil, fmt.Errorf("%s: %s redeclared in this block", parameterToken.pos.toString(), name)
		}
		parser.localScope.InsertExpr(name, parameter)
		return parameter, nil
	}
	return nil, fmt.Errorf("%s: unexpected %s, expected )", parameterToken.pos.toString(), parameterToken.Value)
}

func (parser *parser) parseType() (*Type, error) {
	token := parser.peek()
	switch token.Kind {
	case TOKEN_LBRACE:
		return nil, nil
	case TOKEN_IDENTIFIER:
		parser.skip()
		// Assume all types are defined so far.
		ty, _ := parser.globalScope.GetType(token.Value)
		return ty, nil
	}
	return nil, errors.New("expecting type")
}

func (parser *parser) block() (*Block, error) {
	lbraceToken, err := parser.expectString("{")
	if err != nil {
		return nil, err
	}

	var body []Expr
	for {
		if parser.peek().Kind == TOKEN_RBRACE {
			parser.skip()
			break
		}
		if parser.peek().Kind == TOKEN_EOF {
			return nil, fmt.Errorf("%s: unexpected EOF, expecting }", parser.peek().pos.toString())
		}

		node, err := parser.stmt()
		if err != nil {
			return nil, err
		}
		body = append(body, node)
		// TODO: Semicolon can be omitted after } for one-line function definition.
		if err := parser.consumeString(";"); err != nil {
			return nil, err
		}
	}

	return &Block{tok: lbraceToken, Body: body}, nil
}

func (parser *parser) shortVarDecl(lhs Expr) (Expr, error) {
	if _, ok := lhs.(*Identifier); !ok {
		err := fmt.Errorf("%s: unexpected %s, expecting variable", lhs.token().pos.toString(), lhs.token().Value)
		return nil, err
	}

	lhsVar := &Variable{tok: lhs.token(), Name: lhs.token().Value, Ty: &TypeUnresolved}
	if parser.localScope.ExistsExpr(lhsVar.Name) {
		return nil, errors.New("no new variables on left side of :=")
	}
	parser.localScope.InsertExpr(lhsVar.Name, lhsVar)

	parser.consumeString(":=")
	rhs, err := parser.addOp()
	if err != nil {
		return nil, err
	}

	return &Assign{tok: &Token{Kind: TOKEN_COLONEQUAL, Value: ":="}, Lhs: lhsVar, Rhs: rhs}, nil
}

func (parser *parser) addOp() (Expr, error) {
	lhs, err := parser.primaryExpr()
	if err != nil {
		return nil, err
	}

	token := parser.peek()
	switch token.Kind {
	case TOKEN_PLUS:
		parser.skip()
		rhs, err := parser.addOp()
		if err != nil {
			return nil, err
		}
		return &AddOp{tok: token, Lhs: lhs, Rhs: rhs}, nil
	default:
		return lhs, nil
	}
}

func (parser *parser) primaryExpr() (Expr, error) {
	token := parser.peek()
	switch token.Kind {
	case TOKEN_INT:
		parser.skip()
		return &IntLiteral{tok: token, Value: token.Value}, nil
	case TOKEN_IDENTIFIER:
		parser.skip()
		if token.Value == "true" {
			return &BoolLiteral{tok: token, Value: true}, nil
		} else if token.Value == "false" {
			return &BoolLiteral{tok: token, Value: false}, nil
		}

		if parser.peek().Kind == TOKEN_LPAREN {
			return parser.functionCall(token)
		}

		return &Identifier{tok: token, Name: token.Value}, nil
	}
	return nil, fmt.Errorf("%s: unexpected %s, expecting primary expression", token.pos.toString(), token.Value)
}

func (parser *parser) functionCall(token *Token) (Expr, error) {
	if err := parser.consumeString("("); err != nil {
		return nil, err
	}
	arguments := []Expr{}
	if parser.peek().Kind != TOKEN_RPAREN {
		if argument, err := parser.addOp(); err != nil {
			return nil, err
		} else {
			arguments = append(arguments, argument)
		}
		for {
			if parser.peek().Kind == TOKEN_RPAREN {
				break
			}
			if err := parser.consumeString(","); err != nil {
				return nil, err
			}
			if argument, err := parser.addOp(); err != nil {
				return nil, err
			} else {
				arguments = append(arguments, argument)
			}
		}
	}
	if err := parser.consumeString(")"); err != nil {
		return nil, err
	}
	return &FunctionCall{tok: token, Arguments: arguments}, nil
}
