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
	Id   TypeID
	Size int // Size on a memory in bytes.
	Name string
}

func (ty *Type) GetSize() int {
	return ty.Size
}

func isSameType(ty *Type, other *Type) bool {
	if ty == nil || other == nil {
		return false
	}
	return ty.Id == other.Id
}

func (ty *Type) isUnresolved() bool {
	return ty.Id == TypeIdUnresolved
}

var TypeUnresolved = Type{Id: TypeIdUnresolved, Size: 0}
var TypeBool = Type{Id: TypeIdInt, Size: 16, Name: "bool"}
var TypeInt = Type{Id: TypeIdBool, Size: 16, Name: "int"}

func (ast *Ast) InferType() error {
	for _, f := range ast.funcs {
		if _, err := InferTypeForNode(f, f.Scope); err != nil {
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
		returnType, err := InferTypeForNode(expr.Body, scope)
		if err != nil {
			return nil, err
		}
		actualType := expr.ReturnType
		if returnType == nil && actualType != nil {
			return nil, fmt.Errorf("not enough return values\n\thave: ()\n\twant: (%s)", actualType.Name)
		}
		if returnType != nil && actualType == nil {
			return nil, fmt.Errorf("too many return values\n\thave: (%s)\n\twant: ()", returnType.Name)
		}
		if returnType != actualType {
			return nil, fmt.Errorf("cannot use %s as %s in return statement", returnType.Name, actualType.Name)
		}
	case *Block:
		var returnType *Type
		for _, node := range expr.Body {
			if ty, err := InferTypeForNode(node, scope); err != nil {
				return nil, err
			} else if _, ok := node.(*Return); ok {
				returnType = ty
			}
		}
		return returnType, nil
	case *Return:
		if node, ok := expr.Node.(*Identifier); ok {
			if variable, ok := scope.GetExpr(node.Name()); !ok {
				return nil, fmt.Errorf("undefined: %s", node.Name())
			} else if variable := variable.(*Variable); variable.Ty.isUnresolved() {
				return nil, fmt.Errorf("undefined: %s", node.Name())
			}
		}
		return InferTypeForNode(expr.Node, scope)
	case *Assign:
		rhsType, err := InferTypeForNode(expr.Rhs, scope)
		if err != nil {
			return nil, err
		}
		variable, ok := expr.Lhs.(*Variable)
		if !ok {
			return nil, errors.New("non-name on left side of :=")
		}
		if variable, ok := scope.GetExpr(variable.Name); ok {
			if variable := variable.(*Variable); variable.Ty.isUnresolved() {
				variable.Ty = rhsType
				return rhsType, nil
			} else {
				return nil, errors.New("no new variables on left side of :=")
			}
		}
	case *AddOp:
		lhsType, err := InferTypeForNode(expr.Lhs, scope)
		if err != nil {
			return nil, err
		}
		rhsType, err := InferTypeForNode(expr.Rhs, scope)
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
			expr.Variable = variable
			return variable.Ty, nil
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
			expr.Function = function
			return function.ReturnType, nil
		}
		return nil, fmt.Errorf("invalid operation: cannot call non-function %s", expr.Name())
	}
	return nil, nil
}
