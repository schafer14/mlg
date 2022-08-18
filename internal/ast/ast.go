package ast

import (
	"github.com/flynn/json5"
	"golang.org/x/mod/modfile"
)

// AST represents a syntax tree of the Lamdba definition.
type AST struct {
	Module       *modfile.File
	Namespace    string       `json:"namespace" validate:"required"`
	Name         string       `json:"name" validate:"required"`
	Description  string       `json:"description"`
	Dependencies []Dependency `json:"dependencies"`
	Triggers     []Trigger    `json:"triggers"`
	Emits        []string     `json:"emits"`
}

// Trigger is how the function is triggered.
type Trigger struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

// Dependency is a list of dependencies required.
type Dependency struct {
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
}

// Unmarshal a json5 into an AST.
func Unmarshal(in []byte) (AST, error) {
	var ast AST
	err := json5.Unmarshal(in, &ast)
	return ast, err

}
