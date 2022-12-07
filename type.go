package main

import (
	"errors"
	"fmt"
)

type TypeID int

const (
	TypeIdUnresolved = iota
	// TypeFunction
	TypeIdInt
	TypeIdBool
)

type Type struct {
	id   TypeID
	size int // Size on a memory in bytes.
	name string
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

// func NewFunctionType(returnType *Type) *Type {
// 	return &Type{
// 		id:         TypeFunction,
// 		size:       0,
// 		name:       fmt.Sprintf("func () %s", returnType.name),
// 		returnType: returnType,
// 	}
// }

var TypeUnresolved = Type{id: TypeIdUnresolved, size: 0}
var TypeBool = Type{id: TypeIdInt, size: 16, name: "bool"}
var TypeInt = Type{id: TypeIdBool, size: 16, name: "int"}

func (ast *Ast) InferType() error {
	for _, f := range ast.funcs {
		if _, err := InferTypeForNode(f, f.scope); err != nil {
			return err
		}
	}
	return nil
}

// Traverse AST and determine a type for defined variables.
// Returns pointer to a determined `Type`.
func InferTypeForNode(expr Expr, scope *Scope) (*Type, error) {
	switch expr := expr.(type) {
	case *FunctionDecl:
		if returnType, err := InferTypeForNode(expr.body, scope); err != nil {
			return nil, err
		} else {
			if returnType == nil && expr.returnType != nil {
				return nil, fmt.Errorf("not enough return values\n\thave: ()\n\twant: (%s)", expr.returnType.name)
			}
			if returnType != nil && expr.returnType == nil {
				return nil, fmt.Errorf("too many return values\n\thave: (%s)\n\twant: ()", returnType.name)
			}
			if returnType != expr.returnType {
				return nil, fmt.Errorf("cannot use %s as %s in return statement", returnType.name, expr.returnType.name)
			}
		}
	case *Block:
		var returnType *Type
		for _, node := range expr.body {
			if ty, err := InferTypeForNode(node, scope); err != nil {
				return nil, err
			} else if _, ok := node.(*Return); ok {
				returnType = ty
			}
		}
		return returnType, nil
	case *Return:
		if node, ok := expr.node.(*Identifier); ok {
			if variable, ok := scope.GetExpr(node.token().value); !ok {
				return nil, fmt.Errorf("undefined: %s", node.token().value)
			} else if variable := variable.(*Variable); variable.ty.isUnresolved() {
				return nil, fmt.Errorf("undefined: %s", node.token().value)
			}
		}
		return InferTypeForNode(expr.node, scope)
	case *Assign:
		rhsType, err := InferTypeForNode(expr.rhs, scope)
		if err != nil {
			return nil, err
		}
		if variable, ok := expr.lhs.(*Variable); ok {
			if variable, ok := scope.GetExpr(variable.token().value); ok {
				if variable := variable.(*Variable); variable.ty.isUnresolved() {
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
		lhsType, err := InferTypeForNode(expr.lhs, scope)
		if err != nil {
			return nil, err
		}
		rhsType, err := InferTypeForNode(expr.rhs, scope)
		if err != nil {
			return nil, err
		}
		if !lhsType.isSameType(*rhsType) {
			return nil, errors.New("invalid operation: adding different types")
		}
		return lhsType, nil
	case *Identifier:
		if variable, ok := scope.GetExpr(expr.token().value); ok {
			if variable, ok := variable.(*Variable); ok {
				expr.variable = variable
				return variable.ty, nil
			} else {
				return nil, fmt.Errorf("unexpected %s, expecting variable", expr.token().value)
			}
		}
		return nil, fmt.Errorf("undefined: %s", expr.token().value)
	case *IntLiteral:
		return &TypeInt, nil
	case *BoolLiteral:
		return &TypeBool, nil
	case *FunctionCall:
		if node, ok := scope.GetExpr(expr.token().value); ok {
			if function, ok := node.(*FunctionDecl); ok {
				expr.function = function
				return function.returnType, nil
			}
			return nil, fmt.Errorf("invalid operation: cannot call non-function %s", expr.token().value)
		}
		return nil, fmt.Errorf("undefined: %s", expr.token().value)
	}
	return nil, nil
}
