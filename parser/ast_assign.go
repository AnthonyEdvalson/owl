package parser

import (
	"strings"

	"github.com/AnthonyEdvalson/owl/lexer"
)

/*
assignment = AssignName(name Name)
           | AssignList(parts []assignment)
		   | AssignIndex(target assignment, index expression)
		   | AssignAttribute(target assignment, attribute string)
		   | AssignMap(keyAssign assignment)
		   | AssignSpread(target assignment)
*/

type Assign interface {
	Node
	enforeceAssign()
}

type AssignName struct {
	Name  string
	token lexer.Token
}

type AssignList struct {
	Parts []Assign
	token lexer.Token
}

type AssignIndex struct {
	Target Assign
	Index  Expression
	token  lexer.Token
}

type AssignAttribute struct {
	Target     Expression
	Attribute  string
	IsDeep     bool
	IsCoalesce bool
	token      lexer.Token
}

type AssignMap struct {
	KeyAssign Assign
	token     lexer.Token
}

type AssignSpread struct {
	Target Assign
	token  lexer.Token
}

type AssignNull struct {
	token lexer.Token
}

func (a *AssignName) ToString() string {
	var b strings.Builder

	b.WriteString(a.Name)

	return b.String()
}

func (a *AssignList) ToString() string {
	var b strings.Builder

	for i, part := range a.Parts {
		if i > 0 {
			b.WriteString(", ")
		}

		if part != nil {
			b.WriteString(part.ToString())
		} else {
			b.WriteString("nil")
		}
	}

	return b.String()
}

func (a *AssignIndex) ToString() string {
	var b strings.Builder

	b.WriteString(a.Target.ToString())
	b.WriteString("[")
	b.WriteString(a.Index.ToString())
	b.WriteString("]")

	return b.String()
}

func (a *AssignAttribute) ToString() string {
	var b strings.Builder

	b.WriteString(a.Target.ToString())
	if a.IsDeep {
		b.WriteString("::")
	} else {
		b.WriteString(".")
	}
	b.WriteString(a.Attribute)

	return b.String()
}

func (a *AssignMap) ToString() string {
	var b strings.Builder

	b.WriteString("{")
	b.WriteString(a.KeyAssign.ToString())
	b.WriteString("}")

	return b.String()
}

func (a *AssignSpread) ToString() string {
	var b strings.Builder

	b.WriteString("...")
	b.WriteString(a.Target.ToString())

	return b.String()
}

func (a *AssignNull) ToString() string {
	return "<>"
}

func (a *AssignName) enforeceAssign()      {}
func (a *AssignList) enforeceAssign()      {}
func (a *AssignIndex) enforeceAssign()     {}
func (a *AssignAttribute) enforeceAssign() {}
func (a *AssignMap) enforeceAssign()       {}
func (a *AssignSpread) enforeceAssign()    {}
func (a *AssignNull) enforeceAssign()      {}

func (n *AssignName) Token() lexer.Token      { return n.token }
func (n *AssignList) Token() lexer.Token      { return n.token }
func (n *AssignIndex) Token() lexer.Token     { return n.token }
func (n *AssignAttribute) Token() lexer.Token { return n.token }
func (n *AssignMap) Token() lexer.Token       { return n.token }
func (n *AssignSpread) Token() lexer.Token    { return n.token }
func (n *AssignNull) Token() lexer.Token      { return n.token }
