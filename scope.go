package main

type Scope struct {
	variables map[string]*Variable
	types     map[string]*Type
}

func newScope() *Scope {
	return &Scope{variables: map[string]*Variable{}, types: map[string]*Type{}}
}

func newGlobalScope() *Scope {
	return &Scope{
		variables: map[string]*Variable{},
		types:     map[string]*Type{"int": &TypeInt, "bool": &TypeBool},
	}
}

func (scope *Scope) ExistsVariable(name string) bool {
	_, exists := scope.variables[name]
	return exists
}

func (scope *Scope) InsertVariable(name string, variable *Variable) {
	scope.variables[name] = variable
}

func (scope *Scope) GetVariable(name string) (*Variable, bool) {
	variable, exists := scope.variables[name]
	return variable, exists
}

func (scope *Scope) ExistsType(name string) bool {
	_, exists := scope.types[name]
	return exists
}

func (scope *Scope) InsertType(name string, ty *Type) {
	scope.types[name] = ty
}

func (scope *Scope) GetType(name string) (*Type, bool) {
	ty, exists := scope.types[name]
	return ty, exists
}
