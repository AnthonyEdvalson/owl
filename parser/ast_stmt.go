package parser

import (
	"owl/lexer"
	"strings"
)

/*
statement = Let(target assign, value expr)
		  | For(target expr, iter expr, body []statement)
		  | While(test expr, body []statement)
		  | If(test expr, body []statement, else []statement)
		  | Throw(value expr)
		  | Try(body []statement, catch []statement, finally []statement)
		  | Expression(value expr)
		  | Return(value expr)
		  | Break()
		  | Continue()
		  | Import(name string)
*/

type Statement interface {
	Node
	enforceStatement()
}

type Let struct {
	token   lexer.Token
	Target  Assign
	Value   Expression
}

type For struct {
	token   lexer.Token
	Target Assign
	Iter   Expression
	Body   []Statement
}

type While struct {
	token   lexer.Token
	Test   Expression
	Body   []Statement
}

type If struct {
	token   lexer.Token
	Test   Expression
	Body   []Statement
	Else   []Statement
}

type Throw struct {
	token   lexer.Token
	Value Expression
}

type Try struct {
	token   lexer.Token
	Body   []Statement
	Catch  []Statement
	Finally []Statement
}

type ExpressionStatement struct {
	token   lexer.Token
	Value Expression
}

type Return struct {
	token   lexer.Token
	Value Expression
}

type Break struct{
	token   lexer.Token
}

type Continue struct{
	token   lexer.Token
}

type Import struct {
	token   lexer.Token
	Name    string
}

type Print struct {
	token   lexer.Token
	Value Expression
}

func printBlock(bs *strings.Builder, b []Statement) {
	for _, s := range b {
		if s != nil {
			bs.WriteString(s.ToString())
		} else {
			bs.WriteString("nil")
		}
	}
}

func (l *Let) ToString() string {
	var b strings.Builder

	b.WriteString("let ")

	b.WriteString(l.Target.ToString())

	b.WriteString(" = ")
	if l.Value != nil {
		b.WriteString(l.Value.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("\n")

	return b.String()
}

func (f *For) ToString() string {
	var b strings.Builder

	b.WriteString("for ")
	b.WriteString(f.Target.ToString())
	b.WriteString(" in ")
	if f.Iter != nil {
		b.WriteString(f.Iter.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(" {\n")
	printBlock(&b, f.Body)
	b.WriteString("}\n")

	return b.String()
}

func (w *While) ToString() string {
	var b strings.Builder

	b.WriteString("while ")
	if w.Test != nil {
		b.WriteString(w.Test.ToString())
	} else {
		b.WriteString("nil")
	}

	b.WriteString(" {\n")
	printBlock(&b, w.Body)
	b.WriteString("}\n")

	return b.String()
}

func (i *If) ToString() string {
	var b strings.Builder

	b.WriteString("if ")
	b.WriteString(i.Test.ToString())
	b.WriteString(" {\n")
	printBlock(&b, i.Body)
	b.WriteString("}\n")

	if len(i.Else) > 0 {
		b.WriteString("else {\n")
		printBlock(&b, i.Else)
		b.WriteString("}\n")
	}

	return b.String()
}

func (t *Throw) ToString() string {
	var b strings.Builder

	b.WriteString("throw ")
	b.WriteString(t.Value.ToString())
	b.WriteString("\n")

	return b.String()
}

func (t *Try) ToString() string {
	var b strings.Builder

	b.WriteString("try {\n")
	printBlock(&b, t.Body)
	b.WriteString("}\n")

	b.WriteString("catch {\n")
	printBlock(&b, t.Catch)
	b.WriteString("}\n")

	b.WriteString("finally {\n")
	printBlock(&b, t.Finally)
	b.WriteString("}\n")

	return b.String()
}

func (e *ExpressionStatement) ToString() string {
	if e.Value != nil {
		return e.Value.ToString() + "\n"
	} else {
		return "nil"
	}
}

func (r *Return) ToString() string {
	var b strings.Builder

	b.WriteString("return ")
	if r.Value != nil {
		b.WriteString(r.Value.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("\n")

	return b.String()
}

func (b *Break) ToString() string {
	return "break\n"
}

func (c *Continue) ToString() string {
	return "continue\n"
}

func (i *Import) ToString() string {
	var b strings.Builder

	b.WriteString("import ")
	b.WriteString(i.Name)
	b.WriteString("\n")

	return b.String()
}

func (p *Print) ToString() string {
	var b strings.Builder

	b.WriteString("print ")
	if p.Value != nil {
		b.WriteString(p.Value.ToString())
	} else {
		b.WriteString("nil")
	}
	b.WriteString("\n")

	return b.String()
}

func (s *Let) enforceStatement() {}
func (f *For) enforceStatement() {}
func (w *While) enforceStatement() {}
func (i *If) enforceStatement() {}
func (t *Throw) enforceStatement() {}
func (t *Try) enforceStatement() {}
func (e *ExpressionStatement) enforceStatement() {}
func (r *Return) enforceStatement() {}
func (b *Break) enforceStatement() {}
func (c *Continue) enforceStatement() {}
func (i *Import) enforceStatement() {}
func (p *Print) enforceStatement() {}

func (n *Let) Token() lexer.Token { return n.token }
func (n *For) Token() lexer.Token { return n.token }
func (n *While) Token() lexer.Token { return n.token }
func (n *If) Token() lexer.Token { return n.token }
func (n *Throw) Token() lexer.Token { return n.token }
func (n *Try) Token() lexer.Token { return n.token }
func (n *ExpressionStatement) Token() lexer.Token { return n.token }
func (n *Return) Token() lexer.Token { return n.token }
func (n *Break) Token() lexer.Token { return n.token }
func (n *Continue) Token() lexer.Token { return n.token }
func (n *Import) Token() lexer.Token { return n.token }
func (n *Print) Token() lexer.Token { return n.token }