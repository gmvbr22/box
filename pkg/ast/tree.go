package ast

import "fmt"

const (
	Keyword    = "Keyword"
	Identifier = "Identifier"
	Operator   = "Operator"
	Separator  = "Separator"
	EOF        = "EOF"
)

type Token struct {
	Type   string
	Value  string
	Line   int
	Column int
}

type SingleRule struct {
	Type  string
	Value string
	Child interface{}
}

var program = &SingleRule{
	Type:  Keyword,
	Value: "package",
	Child: &SingleRule{
		Type: Identifier,
		Child: &SingleRule{
			Type: Identifier,
		},
	},
}

var read interface{} = nil

func ReadToken(token *Token) {
	if read == nil {
		read = program
	}
	if rule, isRule := read.(SingleRule); isRule {
		if token.Type != rule.Type {
			panic(fmt.Sprintf("Expected %v but got %v [Line: %d, Col: %d]",
				rule.Type, token.Type, token.Line, token.Column,
			))
		}
		if rule.Value != "" && rule.Value != token.Value {
			panic(fmt.Sprintf("Expected %v but got %v [Line: %d, Col: %d]",
				rule.Value, token.Value, token.Line, token.Column,
			))
		}
		read = rule.Child
		fmt.Printf("Line: %d, Offset: %d, Type: %s, Value: %s\n", token.Line, token.Column, token.Type, token.Value)
		return
	}
}
