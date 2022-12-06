package main

import (
	"errors"
	"fmt"
)

type TypeID int

const (
	TypeIdUnresolved = iota
	TypeIdInt
	TypeIdBool
)

type Type struct {
	id TypeID
	// Size on a memory in bytes.
	size int
}

func (ty *Type) GetSize() int {
	return ty.size
}

func (ty *Type) isSameType(other Type) bool {
	return ty.id == other.id
}

func (ty *Type) isUnresolved() bool {
	return ty.id == TypeIdUnresolved
}

var TypeUnresolved = Type{id: TypeIdUnresolved, size: 0}
var TypeBool = Type{id: TypeIdInt, size: 16}
var TypeInt = Type{id: TypeIdBool, size: 16}

func (ast *Ast) InferType() error {
	for _, f := range ast.funcs {
		if _, err := InferTypeForNode(f, f.localEnv); err != nil {
			return err
		}
	}
	return nil
}

// Traverse AST and determine a type for defined variables.
// Returns pointer to a determined `Type`.
func InferTypeForNode(expr Expr, env *LocalEnv) (*Type, error) {
	switch expr := expr.(type) {
	case *FunctionDecl:
		// TODO: Check return type
		if _, err := InferTypeForNode(expr.body, env); err != nil {
			return nil, err
		}
	case *Block:
		for _, node := range expr.body {
			if _, err := InferTypeForNode(node, env); err != nil {
				return nil, err
			}
		}
	case *Return:
		if node, ok := expr.node.(*Identifier); ok {
			if variable, ok := env.Get(node.token().value); !ok {
				return nil, fmt.Errorf("undefined: %s", node.token().value)
			} else if variable.ty.isUnresolved() {
				return nil, fmt.Errorf("undefined: %s", node.token().value)
			}
		}
		return InferTypeForNode(expr.node, env)
	case *Assign:
		rhsType, err := InferTypeForNode(expr.rhs, env)
		if err != nil {
			return nil, err
		}
		if variable, ok := expr.lhs.(*Variable); ok {
			if variable, ok := env.Get(variable.token().value); ok {
				if variable.ty.isUnresolved() {
					variable.ty = rhsType
					return rhsType, nil
				} else {
					return nil, errors.New("no new variables on left side of :=")
				}
			}
		} else {
			return nil, errors.New("non-name on left side of :=")
		}
	case *AddOp:
		lhsType, err := InferTypeForNode(expr.lhs, env)
		if err != nil {
			return nil, err
		}
		rhsType, err := InferTypeForNode(expr.rhs, env)
		if err != nil {
			return nil, err
		}
		if !lhsType.isSameType(*rhsType) {
			return nil, errors.New("invalid operation: adding different types")
		}
		return lhsType, nil
	case *Identifier:
		if variable, ok := env.Get(expr.token().value); ok {
			expr.variable = variable
			return variable.ty, nil
		} else {
			return nil, fmt.Errorf("undefined: %s", expr.token().value)
		}
	case *IntLiteral:
		return &TypeInt, nil
	case *BoolLiteral:
		return &TypeBool, nil
	}
	return nil, nil
}
