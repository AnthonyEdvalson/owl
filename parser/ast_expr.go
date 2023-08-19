package parser

import (
	"owl/lexer"
	"strconv"
	"strings"
)

/*
expression = AssignExpression(target assign, op binop, value expr)
		   | BinOp(left expr, op binop, right expr)
           | UnaryOp(op unaryop, right expr)
		   | FunctionCall(target expr, args []expr)
		   | FunctionDef(args []string, body expr)
		   | IfExpression(test expr, iftrue expr, iffalse expr)
		   | Map(keys []expr, values []expr)
		   | Set(values []expr)
		   | List(values []expr)
		   | Constant(value constant)
		   | Attrubute(target expr, attr string)
		   | Index(target expr, index expr)
		   | Name(name string)
		   | Spread(target expr)
*/

type Expression interface {
	Node
	enforceExpression()
}

type AssignExpression struct {
	Target Assign
	Op     string
	Value  Expression
	token  lexer.Token
}

type BinOp struct {
	Left  Expression
	Op    string
	Right Expression
	token lexer.Token
}

type List struct {
	Parts []Expression
	token lexer.Token
}

type UnaryOp struct {
	Op    string
	Value Expression
	token lexer.Token
}

type IncDec struct {
	Target Assign
	Op     string
	token  lexer.Token
}

type FunctionCall struct {
	Target     Expression
	Arg        Expression
	IsCoalesce bool
	token      lexer.Token
}

type FunctionDef struct {
	Arg       Assign
	Condition Expression
	Body      []Statement
	Else      *FunctionDef
	token     lexer.Token
}

type IfExpression struct {
	Test    Expression
	IfTrue  Expression
	IfFalse Expression
	token   lexer.Token
}

type Map struct {
	Keys   []string
	Values []Expression
	token  lexer.Token
}

type Set struct {
	Values []Expression
	token  lexer.Token
}

type Const struct {
	Value interface{}
	token lexer.Token
}

type Null struct {
	token lexer.Token
}

type Attribute struct {
	Target     Expression
	Attribute  string
	IsDeep     bool
	IsCoalesce bool
	token      lexer.Token
}

type Index struct {
	Target Expression
	Index  Expression
	token  lexer.Token
}

type Slice struct {
	Target Expression
	Start  Expression
	End    Expression
	token  lexer.Token
}

type Name struct {
	Name  string
	token lexer.Token
}

type Spread struct {
	Target Expression
	token  lexer.Token
}

type Overload struct {
	Cases []FunctionDef
	token lexer.Token
}

func (a *AssignExpression) ToString() string {
	var b strings.Builder

	if a.Target != nil {
		b.WriteString(a.Target.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(" ")
	b.WriteString(a.Op)
	b.WriteString(" ")

	if a.Value != nil {
		b.WriteString(a.Value.ToString())
	} else {
		b.WriteString("nil")
	}

	return b.String()
}

func (bo *BinOp) ToString() string {
	var b strings.Builder

	b.WriteString("(")
	if bo.Left != nil {
		b.WriteString(bo.Left.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(" ")
	b.WriteString(bo.Op)
	b.WriteString(" ")

	if bo.Right != nil {
		b.WriteString(bo.Right.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(")")

	return b.String()
}

func (l *List) ToString() string {
	var b strings.Builder

	b.WriteString("[")

	for i, part := range l.Parts {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(part.ToString())
	}

	b.WriteString("]")

	return b.String()
}

func (u *UnaryOp) ToString() string {
	var b strings.Builder

	b.WriteString("(")

	b.WriteString(u.Op)
	if u.Value != nil {
		b.WriteString(u.Value.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(")")

	return b.String()
}

func (o *IncDec) ToString() string {
	var b strings.Builder

	if o.Target != nil {
		b.WriteString(o.Target.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(o.Op)

	return b.String()
}

func (f *FunctionCall) ToString() string {
	var b strings.Builder

	if f.Target != nil {
		b.WriteString(f.Target.ToString())
	} else {
		b.WriteString("nil")
	}
	if f.IsCoalesce {
		b.WriteString("?")
	}
	b.WriteString("(")
	if f.Arg != nil {
		b.WriteString(f.Arg.ToString())
	}
	b.WriteString(")")

	return b.String()
}

func (f *FunctionDef) ToString() string {
	var b strings.Builder

	if f.Condition != nil {
		b.WriteString("when ")
		b.WriteString(f.Condition.ToString())
	}

	b.WriteString("(")
	if f.Arg != nil {
		b.WriteString(f.Arg.ToString())
	}
	b.WriteString(") => {\n")

	printBlock(&b, f.Body)

	b.WriteString("}")

	return b.String()
}

func (i *IfExpression) ToString() string {
	var b strings.Builder

	b.WriteString("(")

	if i.Test != nil {
		b.WriteString(i.Test.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(" ? ")

	if i.IfTrue != nil {
		b.WriteString(i.IfTrue.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString(" : ")

	if i.IfFalse != nil {
		b.WriteString(i.IfFalse.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(")")

	return b.String()
}

func (d *Map) ToString() string {
	var b strings.Builder

	b.WriteString("{\n")

	for i, k := range d.Keys {
		b.WriteString(k)

		b.WriteString(": ")

		if d.Values[i] != nil {
			b.WriteString(d.Values[i].ToString())
		} else {
			b.WriteString("nil")
		}

		if i < len(d.Keys)-1 {
			b.WriteString(",")
		}

		b.WriteString("\n")
	}

	b.WriteString("}")

	return b.String()
}

func (s *Set) ToString() string {
	var b strings.Builder

	b.WriteString("{")
	for i, e := range s.Values {
		if e != nil {
			b.WriteString(e.ToString())
		} else {
			b.WriteString("nil")
		}
		if i < len(s.Values)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("}")

	return b.String()
}

func (c *Const) ToString() string {
	switch v := c.Value.(type) {
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return "\"" + v + "\""
	case bool:
		return strconv.FormatBool(v)
	default:
		return "(Unknown Const)"
	}
}

func (n *Null) ToString() string {
	var b strings.Builder

	b.WriteString("null")

	return b.String()
}

func (a *Attribute) ToString() string {
	var b strings.Builder

	if a.Target != nil {
		b.WriteString(a.Target.ToString())
	} else {
		b.WriteString("nil")
	}

	if a.IsCoalesce {
		b.WriteString("?")
	}

	if a.IsDeep {
		b.WriteString("::")
	} else {
		b.WriteString(".")
	}

	b.WriteString(a.Attribute)

	return b.String()
}

func (i *Index) ToString() string {
	var b strings.Builder

	if i.Target != nil {
		b.WriteString(i.Target.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("[")
	if i.Index != nil {
		b.WriteString(i.Index.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("]")

	return b.String()
}

func (i *Slice) ToString() string {
	var b strings.Builder

	if i.Target != nil {
		b.WriteString(i.Target.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("[")
	if i.Start != nil {
		b.WriteString(i.Start.ToString())
	}
	b.WriteString(":")
	if i.End != nil {
		b.WriteString(i.End.ToString())
	}
	b.WriteString("]")

	return b.String()
}

func (n *Name) ToString() string {
	return n.Name
}

func (s *Spread) ToString() string {
	var b strings.Builder

	b.WriteString("...")
	if s.Target != nil {
		b.WriteString(s.Target.ToString())
	} else {
		b.WriteString("nil")
	}

	return b.String()
}

func (o *Overload) ToString() string {
	var b strings.Builder

	b.WriteString("<")

	for i, c := range o.Cases {
		b.WriteString(c.ToString())
		if i < len(o.Cases)-1 {
			b.WriteString(" | ")
		}
	}

	b.WriteString(">")

	return b.String()
}

func (a *AssignExpression) enforceExpression() {}
func (b *BinOp) enforceExpression()            {}
func (c *List) enforceExpression()             {}
func (u *UnaryOp) enforceExpression()          {}
func (i *IncDec) enforceExpression()           {}
func (f *FunctionCall) enforceExpression()     {}
func (f *FunctionDef) enforceExpression()      {}
func (i *IfExpression) enforceExpression()     {}
func (d *Map) enforceExpression()              {}
func (s *Set) enforceExpression()              {}
func (c *Const) enforceExpression()            {}
func (n *Null) enforceExpression()             {}
func (a *Attribute) enforceExpression()        {}
func (i *Index) enforceExpression()            {}
func (s *Slice) enforceExpression()            {}
func (n *Name) enforceExpression()             {}
func (s *Spread) enforceExpression()           {}
func (o *Overload) enforceExpression()         {}

func (n *AssignExpression) Token() lexer.Token { return n.token }
func (n *BinOp) Token() lexer.Token            { return n.token }
func (n *List) Token() lexer.Token             { return n.token }
func (n *UnaryOp) Token() lexer.Token          { return n.token }
func (n *IncDec) Token() lexer.Token           { return n.token }
func (n *FunctionCall) Token() lexer.Token     { return n.token }
func (n *FunctionDef) Token() lexer.Token      { return n.token }
func (n *IfExpression) Token() lexer.Token     { return n.token }
func (n *Map) Token() lexer.Token              { return n.token }
func (n *Set) Token() lexer.Token              { return n.token }
func (n *Const) Token() lexer.Token            { return n.token }
func (n *Null) Token() lexer.Token             { return n.token }
func (n *Attribute) Token() lexer.Token        { return n.token }
func (n *Index) Token() lexer.Token            { return n.token }
func (n *Slice) Token() lexer.Token            { return n.token }
func (n *Name) Token() lexer.Token             { return n.token }
func (n *Spread) Token() lexer.Token           { return n.token }
func (n *Overload) Token() lexer.Token         { return n.token }
