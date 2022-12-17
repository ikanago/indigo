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

func isSameType(ty *Type, other *Type) bool {
	if ty == nil || other == nil {
		return false
	}
	return ty.id == other.id
}

func (ty *Type) isUnresolved() bool {
	return ty.id == TypeIdUnresolved
}

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
		returnType, err := InferTypeForNode(expr.body, scope)
		if err != nil {
			return nil, err
		}
		actualType := expr.returnType
		if returnType == nil && actualType != nil {
			return nil, fmt.Errorf("not enough return values\n\thave: ()\n\twant: (%s)", actualType.name)
		}
		if returnType != nil && actualType == nil {
			return nil, fmt.Errorf("too many return values\n\thave: (%s)\n\twant: ()", returnType.name)
		}
		if returnType != actualType {
			return nil, fmt.Errorf("cannot use %s as %s in return statement", returnType.name, actualType.name)
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
			if variable, ok := scope.GetExpr(node.Name()); !ok {
				return nil, fmt.Errorf("undefined: %s", node.Name())
			} else if variable := variable.(*Variable); variable.ty.isUnresolved() {
				return nil, fmt.Errorf("undefined: %s", node.Name())
			}
		}
		return InferTypeForNode(expr.node, scope)
	case *Assign:
		rhsType, err := InferTypeForNode(expr.rhs, scope)
		if err != nil {
			return nil, err
		}
		variable, ok := expr.lhs.(*Variable)
		if !ok {
			return nil, errors.New("non-name on left side of :=")
		}
		if variable, ok := scope.GetExpr(variable.Name()); ok {
			if variable := variable.(*Variable); variable.ty.isUnresolved() {
				variable.ty = rhsType
				return rhsType, nil
			} else {
				return nil, errors.New("no new variables on left side of :=")
			}
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
		if !isSameType(lhsType, rhsType) {
			return nil, errors.New("invalid operation: adding different types")
		}
		return lhsType, nil
	case *Identifier:
		variable, ok := scope.GetExpr(expr.Name())
		if !ok {
			return nil, fmt.Errorf("undefined: %s", expr.Name())
		}
		if variable, ok := variable.(*Variable); ok {
			expr.variable = variable
			return variable.ty, nil
		}
		return nil, fmt.Errorf("unexpected %s, expecting variable", expr.Name())
	case *IntLiteral:
		return &TypeInt, nil
	case *BoolLiteral:
		return &TypeBool, nil
	case *FunctionCall:
		node, ok := scope.GetExpr(expr.Name())
		if !ok {
			return nil, fmt.Errorf("undefined: %s", expr.Name())
		}
		if function, ok := node.(*FunctionDecl); ok {
			expr.function = function
			return function.returnType, nil
		}
		return nil, fmt.Errorf("invalid operation: cannot call non-function %s", expr.Name())
	}
	return nil, nil
}
