package main

import (
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
			return nil, fmt.Errorf("%s: not enough return values\n\thave: ()\n\twant: (%s)", expr.token().pos.toString(), actualType.Name)
		}
		if returnType != nil && actualType == nil {
			return nil, fmt.Errorf("%s: too many return values\n\thave: (%s)\n\twant: ()", expr.token().pos.toString(), returnType.Name)
		}
		if returnType != actualType {
			return nil, fmt.Errorf("%s: cannot use %s as %s in return statement", expr.token().pos.toString(), returnType.Name, actualType.Name)
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
			if variable, ok := scope.GetExpr(node.Name); !ok {
				return nil, fmt.Errorf("%s: undefined: %s", node.token().pos.toString(), node.Name)
			} else if variable := variable.(*Variable); variable.Ty.isUnresolved() {
				return nil, fmt.Errorf("%s: undefined: %s", node.token().pos.toString(), node.Name)
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
			return nil, fmt.Errorf("%s: non-name on left side of :=", expr.Lhs.token().pos.toString())
		}
		if variable, ok := scope.GetExpr(variable.Name); ok {
			if variable := variable.(*Variable); variable.Ty.isUnresolved() {
				variable.Ty = rhsType
				return rhsType, nil
			} else {
				return nil, fmt.Errorf("%s: no new variables on left side of :=", expr.Lhs.token().pos.toString())
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
			return nil, fmt.Errorf("%s: invalid operation: adding different types", expr.token().pos.toString())
		}
		return lhsType, nil
	case *Identifier:
		variable, ok := scope.GetExpr(expr.Name)
		if !ok {
			return nil, fmt.Errorf("%s: undefined: %s", variable.token().pos.toString(), expr.Name)
		}
		if variable, ok := variable.(*Variable); ok {
			expr.Variable = variable
			return variable.Ty, nil
		}
		return nil, fmt.Errorf("%s: unexpected %s, expecting variable", variable.token().pos.toString(), expr.Name)
	case *IntLiteral:
		return &TypeInt, nil
	case *BoolLiteral:
		return &TypeBool, nil
	case *FunctionCall:
		maybeFunctionDecl, ok := scope.GetExpr(expr.Name())
		if !ok {
			return nil, fmt.Errorf("%s: undefined: %s", expr.token().pos.toString(), expr.Name())
		}
		if function, ok := maybeFunctionDecl.(*FunctionDecl); ok {
			expr.Function = function
			return function.ReturnType, nil
		}
		return nil, fmt.Errorf("%s: invalid operation: cannot call non-function %s", expr.token().pos.toString(), expr.Name())
	}
	return nil, nil
}
