package main

type Scope struct {
	exprs map[string]Expr  // key: name, value: corresponding `Expr` in the AST
	types map[string]*Type // key: name, value: defined type
	outer *Scope
}

func NewScope(outer *Scope) *Scope {
	return &Scope{exprs: map[string]Expr{}, types: map[string]*Type{}, outer: outer}
}

func NewGlobalScope() *Scope {
	return &Scope{
		exprs: map[string]Expr{},
		types: map[string]*Type{"int": &TypeInt, "bool": &TypeBool},
		outer: nil,
	}
}

func (scope *Scope) ExistsExpr(name string) bool {
	_, exists := scope.GetExpr(name)
	return exists
}

func (scope *Scope) InsertExpr(name string, expr Expr) {
	scope.exprs[name] = expr
}

func (scope *Scope) GetExpr(name string) (Expr, bool) {
	expr, ok := scope.exprs[name]
	if ok {
		return expr, ok
	}
	if scope.outer != nil {
		return scope.outer.GetExpr(name)
	}
	return nil, false
}

func (scope *Scope) ExistsType(name string) bool {
	_, exists := scope.types[name]
	return exists
}

func (scope *Scope) InsertType(name string, ty *Type) {
	scope.types[name] = ty
}

func (scope *Scope) GetType(name string) (*Type, bool) {
	ty, ok := scope.types[name]
	if ok {
		return ty, ok
	}
	if scope.outer != nil {
		return scope.outer.GetType(name)
	}
	return nil, false
}
