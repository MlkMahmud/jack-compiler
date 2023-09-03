package symboltable

import (
	"fmt"
)

type Symbol struct {}

type SymbolTable struct {
	Enclosing *SymbolTable
	Values    map[string]Symbol
}

func NewSymbolTable(enclosing *SymbolTable) *SymbolTable {
	return &SymbolTable{ Enclosing: enclosing }
}

func (table *SymbolTable) Add(id string, symbol Symbol) {
	_, ok := table.Values[id]

	if ok {
		panic(fmt.Sprintf("SyntaxError: Identifier '%s' has already been declared", id))
	}

	table.Values[id] = symbol
}

func (table *SymbolTable) Get(id string) Symbol {
	symbol, ok := table.Values[id]

	if ok {
		return symbol
	}

	if table.Enclosing != nil {
		return table.Enclosing.Get(id)
	}

	panic(fmt.Sprintf("ReferenceError: '%s' is not defined.", id))
}
