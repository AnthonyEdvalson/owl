package parser

import (
	"fmt"
	"owl/lexer"
	"strconv"
	"strings"
)

// ======================================================================================
//
//                                    Parser Definition
//
// ======================================================================================

// Types

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type ParserError struct {
	Message string
	Token   lexer.Token
}

type Parser struct {
	position int
	input    []lexer.Token

	Errors []ParserError

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

// Constants

// Precedences loosely based on Python, C++, and C# precedence
// https://docs.python.org/3/reference/expressions.html
// https://en.cppreference.com/w/cpp/language/operator_precedence
// https://docs.microsoft.com/en-us/dotnet/csharp/language-reference/operators/
// Parser will do high precedence operations before low

var precedences = map[lexer.TokenType]int{
	"LPAREN":              HIGH,
	"QUESTIONLPAREN":      HIGH,
	"LBRACKET":            HIGH,
	"LBRACE":              HIGH,
	"INCDEC":              HIGH,
	"CALL":                HIGH,
	"ATTRIBUTE":           HIGH,
	"INDEX":               HIGH,
	"DOT":                 HIGH,
	"DOUBLECOLON":         HIGH,
	"QUESTIONDOT":         HIGH,
	"QUESTIONDOUBLECOLON": HIGH,
	"DOUBLESTAR":          POWER,
	"PERCENT":             MULDIV,
	"SLASH":               MULDIV,
	"STAR":                MULDIV,
	"PLUS":                ADDSUB,
	"MINUS":               ADDSUB,
	"COMPARE":             COMPARE,
	"HAS":                 COMPARE,
	"AND":                 AND,
	"OR":                  OR,
	"DOUBLEQUESTION":      COALESCE, // TODO: null coalesce is right associative
	"QUESTION":            IFEXP,
	"COLON":               IFEXP,
	"ARROW":               ARROW,
	"TRIPLEDOT":           ARROW,
	"PIPE":                OVERLOAD,
	"COMMA":               COMMA,
	"ASSIGN":              ASSIGN,
}

const (
	LOW = iota
	ASSIGN
	COMMA
	OVERLOAD
	ARROW
	IFEXP
	COALESCE
	OR
	AND
	COMPARE
	ADDSUB
	MULDIV
	PREFIX
	POWER
	HIGH
)

// Parser

func NewParser(input []lexer.Token) *Parser {
	p := &Parser{position: 0, input: input}
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix("NAME", func() Expression { return p.parseName() })
	p.registerPrefix("NULL", p.parseNull)
	p.registerPrefix("NUMBER", p.parseNumber)
	p.registerPrefix("STRING", p.parseString)
	p.registerPrefix("BOOL", p.parseBool)
	p.registerPrefix("NOT", p.parsePrefix)
	p.registerPrefix("MINUS", p.parsePrefix)
	p.registerPrefix("LPAREN", p.parseParens)
	p.registerPrefix("LBRACE", p.parseBrace)
	p.registerPrefix("LBRACKET", p.parseBracket)
	p.registerPrefix("TRIPLEDOT", p.parseSpread)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix("PLUS", p.parseBinOp)
	p.registerInfix("MINUS", p.parseBinOp)
	p.registerInfix("SLASH", p.parseBinOp)
	p.registerInfix("STAR", p.parseBinOp)
	p.registerInfix("DOUBLESTAR", p.parseBinOp)
	p.registerInfix("PERCENT", p.parseBinOp)
	p.registerInfix("COMPARE", p.parseBinOp)
	p.registerInfix("HAS", p.parseBinOp)
	p.registerInfix("AND", p.parseBinOp)
	p.registerInfix("OR", p.parseBinOp)
	p.registerInfix("DOUBLEQUESTION", p.parseBinOp)
	p.registerInfix("COMMA", p.parseComma)
	p.registerInfix("QUESTION", p.parseIfExpression)
	p.registerInfix("ARROW", p.parseArrow)
	p.registerInfix("DOT", p.parseAttribute)
	p.registerInfix("DOUBLECOLON", p.parseAttribute)
	p.registerInfix("QUESTIONDOT", p.parseAttribute)
	p.registerInfix("QUESTIONDOUBLECOLON", p.parseAttribute)
	p.registerInfix("ASSIGN", p.parseAssignExpression)
	p.registerInfix("LBRACKET", p.parseIndex)
	p.registerInfix("LPAREN", p.parseCall)
	p.registerInfix("QUESTIONLPAREN", p.parseCall)
	p.registerInfix("INCDEC", p.parseIncDec)
	p.registerInfix("PIPE", p.parseOverload)

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) current() lexer.Token {
	i := p.position

	if i >= len(p.input) {
		i = len(p.input) - 1
	}

	return p.input[i]
}

func (p *Parser) next() {
	p.position++
}

func (p *Parser) consume(t lexer.TokenType) {
	if p.current().Type != t {
		p.error(fmt.Sprintf("Expected token %s, got %s", t, p.current().Type), p.current())
	}

	p.next()
}

func (p *Parser) consumeAny(t lexer.TokenType) {
	for p.current().Type == t {
		p.next()
	}
}

func (p *Parser) error(msg string, token lexer.Token) {
	p.Errors = append(p.Errors, ParserError{Message: msg, Token: token})

	if len(p.Errors) > 100 {
		panic("Too many errors")
	}
}

func (p *Parser) Parse() *Program {
	lexer := p.parseProgram()
	return lexer
}

func (p *Parser) parseProgram() *Program {
	program := &Program{}
	program.Body = p.parseBlock(false)

	return program
}

func (p *Parser) parseAssign() Assign {
	exp := p.parseExpression(ASSIGN)
	return p.expressionToAssign(exp)
}

func (p *Parser) expressionToAssign(expr Expression) Assign {
	if expr == nil {
		a := &AssignNull{}
		a.token = p.current()
		return a
	}

	switch e := expr.(type) {
	case *Name:
		a := &AssignName{}
		a.Name = e.Name
		a.token = e.token
		return a

	case *Index:
		a := &AssignIndex{}
		a.Target = p.expressionToAssign(e.Target)
		a.Index = e.Index
		a.token = e.token
		return a

	case *Attribute:
		a := &AssignAttribute{}
		a.Target = e.Target
		a.Attribute = e.Attribute
		a.IsDeep = e.IsDeep
		a.IsCoalesce = e.IsCoalesce
		a.token = e.token
		return a

	case *List:
		a := &AssignList{}
		a.Parts = make([]Assign, len(e.Parts))
		a.token = e.token

		for i, part := range e.Parts {
			a.Parts[i] = p.expressionToAssign(part)
		}

		return a

	/*case *Map:
	a := &AssignMap{}
	a.KeyAssign = expr.(*Map)
	return a*/

	case *Spread:
		a := &AssignSpread{}
		a.Target = p.expressionToAssign(e.Target)
		a.token = e.token
		return a

	default:
		p.error(fmt.Sprintf("Cannot use %T in assignment", expr), p.current())
		return nil
	}
}

// ======================================================================================
//
//                                    Statement Parsing
//
// ======================================================================================

func (p *Parser) parseStatement() Statement {
	switch p.current().Type {
	case "LET":
		return p.parseLet()
	case "FOR":
		return p.parseFor()
	case "WHILE":
		return p.parseWhile()
	case "IF":
		return p.parseIf()
	//case "THROW":
	//	return p.parseThrow()
	//case "TRY":
	//	return p.parseTry()
	case "RETURN":
		return p.parseReturn()
	case "BREAK":
		return p.parseBreak()
	case "CONTINUE":
		return p.parseContinue()
	case "IMPORT":
		return p.parseImport()
	case "PRINT":
		return p.parsePrint()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlock(braces bool) []Statement {
	var block []Statement

	p.consumeAny("NEWLINE")

	if braces {
		p.consume("LBRACE")
	}

	for p.current().Type != "RBRACE" && p.current().Type != "EOF" {
		if p.current().Type == "NEWLINE" {
			p.next()
			continue
		}

		i := p.position
		stmt := p.parseStatement()
		block = append(block, stmt)

		if i == p.position {
			p.error("Parser did not progress on statement "+p.current().Type+", '"+p.current().Literal+"'", p.current())
			break
		}
	}

	if braces {
		p.consume("RBRACE")
	}

	p.consumeAny("NEWLINE")

	return block
}

func (p *Parser) parseLet() *Let {
	l := &Let{}
	l.token = p.current()

	p.consume("LET")

	target := p.parseAssign()
	p.consume("ASSIGN")
	value := p.parseExpression(LOW)

	l.Target = target
	l.Value = value

	return l
}

func (p *Parser) parseName() *Name {
	n := &Name{}
	n.token = p.current()

	name := p.current().Literal
	p.consume("NAME")

	n.Name = name

	return n
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{}
	stmt.token = p.current()
	stmt.Value = p.parseExpression(LOW)

	for p.current().Type == "NEWLINE" {
		p.next()
	}

	return stmt
}

func (p *Parser) parseFor() *For {
	f := &For{}
	f.token = p.current()

	p.consume("FOR")

	f.Target = p.parseAssign()
	p.consume("IN")
	f.Iter = p.parseExpression(LOW)

	f.Body = p.parseBlock(true)

	return f
}

func (p *Parser) parseWhile() *While {
	w := &While{}
	w.token = p.current()

	p.consume("WHILE")

	w.Test = p.parseExpression(LOW)
	w.Body = p.parseBlock(true)

	return w
}

func (p *Parser) parseIf() *If {
	i := &If{}
	i.token = p.current()

	p.consume("IF")

	i.Test = p.parseExpression(LOW)
	p.consumeAny("NEWLINE")
	i.Body = p.parseBlock(true)
	p.consumeAny("NEWLINE")

	if p.current().Type == "ELSE" {
		p.consume("ELSE")

		if p.current().Type == "IF" {
			i.Else = []Statement{p.parseIf()}
		} else {
			i.Else = p.parseBlock(true)
		}
	}

	return i
}

func (p *Parser) parseReturn() *Return {
	r := &Return{}
	r.token = p.current()

	p.consume("RETURN")

	r.Value = p.parseExpression(LOW)

	return r
}

func (p *Parser) parseBreak() *Break {
	b := &Break{}
	b.token = p.current()

	p.consume("BREAK")

	return b
}

func (p *Parser) parseContinue() *Continue {
	c := &Continue{}
	c.token = p.current()

	p.consume("CONTINUE")

	return c
}

func (p *Parser) parseImport() *Import {
	i := &Import{}
	i.token = p.current()

	p.consume("IMPORT")
	i.Name = p.parseString().(*Const).Value.(string)

	return i
}

func (p *Parser) parsePrint() *Print {
	print := &Print{}
	print.token = p.current()

	p.consume("PRINT")
	print.Value = p.parseExpression(LOW)

	return print
}

// Utility

func Precedence(tokenType lexer.TokenType) int {
	if p, ok := precedences[tokenType]; ok {
		return p
	}

	return LOW
}

// ======================================================================================
//
//                                    Expression Parsing
//
// ======================================================================================

func (p *Parser) parseExpression(precedence int) Expression {
	t := p.current()
	prefix := p.prefixParseFns[t.Type]

	if prefix == nil {
		p.error("Unexpected token "+t.Type+" '"+t.Literal+"' is not a registered prefix or infix operator", t)
		return nil
	}

	leftExp := prefix()

	for precedence < Precedence(p.current().Type) {
		t = p.current()

		infix := p.infixParseFns[t.Type]

		if infix == nil {
			return leftExp
		}

		leftExp = infix(leftExp)

		if leftExp == nil {
			break
		}
	}

	return leftExp
}

func (p *Parser) parsePrefix() Expression {
	unop := &UnaryOp{}
	unop.token = p.current()

	unop.Op = p.current().Literal
	p.next()

	unop.Value = p.parseExpression(PREFIX)

	return unop
}

func (p *Parser) parseParens() Expression {
	// Question parens are used for coalesce calls
	if p.current().Type == "QUESTIONLPAREN" {
		p.consume("QUESTIONLPAREN")
	} else {
		p.consume("LPAREN")
	}

	if p.current().Type == "RPAREN" {
		p.next()
		return nil
	}

	inner := p.parseExpression(LOW)

	p.consume("RPAREN")
	return inner
}

func (p *Parser) parseBracket() Expression {
	p.consume("LBRACKET")
	p.consumeAny("NEWLINE")

	tok := p.current()
	var inner Expression

	if p.current().Type != "RBRACKET" {
		inner = p.parseExpression(LOW)
	}

	p.consumeAny("NEWLINE")
	p.consume("RBRACKET")

	switch t := inner.(type) {
	case *List:
		t.token = tok
		return t
	case nil:
		return &List{token: tok}
	default:
		return &List{token: tok, Parts: []Expression{inner}}
	}
}

func (p *Parser) parseBrace() Expression {
	m := &Map{}
	m.token = p.current()

	m.Keys = make([]string, 0)
	m.Values = make([]Expression, 0)

	p.consume("LBRACE")

	for p.current().Type != "RBRACE" && p.current().Type != "EOF" {
		p.consumeAny("NEWLINE")

		var name string

		if p.current().Type == "NAME" {
			name = p.parseName().Name
		} else if p.current().Type == "STRING" {
			name = p.parseString().(*Const).Value.(string)
		} else {
			p.error("Expected NAME or STRING, got "+p.current().Type, p.current())
		}

		p.consume("COLON")
		value := p.parseExpression(COMMA)
		p.consumeAny("COMMA")
		p.consumeAny("NEWLINE")

		m.Keys = append(m.Keys, name)
		m.Values = append(m.Values, value)
	}

	p.consume("RBRACE")

	return m
}

func (p *Parser) parseBinOp(left Expression) Expression {
	bop := &BinOp{}
	bop.token = p.current()

	bop.Left = left
	bop.Op = p.current().Literal

	precedence := Precedence(p.current().Type)
	p.next()
	bop.Right = p.parseExpression(precedence)

	return bop
}

func (p *Parser) parseComma(left Expression) Expression {
	c := &List{}
	c.token = p.current()

	parts := []Expression{left}

	for p.current().Type == "COMMA" {
		p.consumeAny("NEWLINE")
		p.consume("COMMA")
		p.consumeAny("NEWLINE")
		part := p.parseExpression(COMMA)

		if part != nil {
			parts = append(parts, part)
		}
	}

	c.Parts = parts

	return c
}

func (p *Parser) parseAssignExpression(left Expression) Expression {
	a := &AssignExpression{}
	a.token = p.current()

	a.Target = p.expressionToAssign(left)
	a.Op = p.current().Literal
	p.consume("ASSIGN")
	a.Value = p.parseExpression(ASSIGN)
	return a
}

func (p *Parser) parseIndex(left Expression) Expression {
	token := p.current()
	target := left

	var index Expression

	p.consume("LBRACKET")
	if p.current().Type != "COLON" {
		index = p.parseExpression(LOW)
		if p.current().Type != "COLON" {
			p.consume("RBRACKET")
			return &Index{target, index, token}
		}
	}

	p.consume("COLON")
	start := index
	var end Expression

	if p.current().Type != "RBRACKET" {
		end = p.parseExpression(LOW)
	}

	p.consume("RBRACKET")

	return &Slice{target, start, end, token}
}

func (p *Parser) parseIfExpression(left Expression) Expression {
	ie := &IfExpression{}
	ie.token = p.current()

	ie.Test = left

	p.consume("QUESTION")

	ie.IfTrue = p.parseExpression(IFEXP)

	p.consume("COLON")

	ie.IfFalse = p.parseExpression(IFEXP)

	return ie
}

func (p *Parser) parseArrow(left Expression) Expression {
	fd := &FunctionDef{}
	fd.token = p.current()

	fd.Arg = p.expressionToAssign(left)

	p.consume("ARROW")

	if p.current().Type == "LBRACE" {
		fd.Body = p.parseBlock(true)
	} else {
		ret := &Return{Value: p.parseExpression(ARROW)}

		fd.Body = []Statement{ret}
	}
	p.consumeAny("NEWLINE")

	return fd
}

func (p *Parser) parseCall(left Expression) Expression {
	c := &FunctionCall{}
	c.token = p.current()

	c.IsCoalesce = p.current().Type == "QUESTIONLPAREN"
	c.Target = left
	c.Arg = p.parseParens()

	return c
}

func (p *Parser) parseIncDec(left Expression) Expression {
	incDec := &IncDec{}
	incDec.token = p.current()

	incDec.Op = p.current().Literal
	incDec.Target = p.expressionToAssign(left)

	p.consume("INCDEC")

	return incDec
}

func (p *Parser) parseAttribute(left Expression) Expression {
	d := &Attribute{}
	d.token = p.current()

	var isDeep = false
	var isCoalesce = false

	if p.current().Type == "DOT" {
		p.consume("DOT")
	} else if p.current().Type == "DOUBLECOLON" {
		p.consume("DOUBLECOLON")
		isDeep = true
	} else if p.current().Type == "QUESTIONDOT" {
		p.consume("QUESTIONDOT")
		isCoalesce = true
	} else if p.current().Type == "QUESTIONDOUBLECOLON" {
		p.consume("QUESTIONDOUBLECOLON")
		isCoalesce = true
		isDeep = true
	} else {
		p.error("Unexpected token "+p.current().Type+" '"+p.current().Literal+"' is not a valid attribute operator", p.current())
		return nil
	}

	// Propagate coalesce
	attr, isAttr := left.(*Attribute)
	if isAttr && attr.IsCoalesce {
		isCoalesce = true
	}

	d.Target = left
	d.Attribute = p.current().Literal
	d.IsCoalesce = isCoalesce
	d.IsDeep = isDeep
	p.consume("NAME")
	return d
}

// Constants

func (p *Parser) parseNull() Expression {
	n := &Null{}
	n.token = p.current()
	p.consume("NULL")
	return n
}

func (p *Parser) parseNumber() Expression {
	c := &Const{}
	c.token = p.current()

	s := p.current().Literal
	p.consume("NUMBER")

	var v interface{}

	fVal, fErr := strconv.ParseFloat(s, 64)
	iVal, iErr := strconv.ParseInt(s, 0, 64)

	if iErr == nil {
		v = iVal
	} else if fErr == nil {
		v = fVal
	} else {
		msg := fmt.Sprintf("could not parse %q as integer", s)
		p.error(msg, p.current())
		return nil
	}

	c.Value = v
	return c
}

func (p *Parser) parseString() Expression {
	c := &Const{}
	c.token = p.current()

	s := p.current().Literal
	p.consume("STRING")

	s = strings.Replace(s, `\n`, "\n", -1)
	s = strings.Replace(s, `\t`, "\t", -1)
	s = strings.Replace(s, `\r`, "\r", -1)
	s = strings.Replace(s, `\f`, "\f", -1)
	s = strings.Replace(s, `\v`, "\v", -1)
	s = strings.Replace(s, `\\`, "\\", -1)
	s = strings.Replace(s, `\"`, "\"", -1)
	s = strings.Replace(s, `\'`, "'", -1)
	s = strings.Replace(s, `\x`, "\\x", -1)
	s = strings.Replace(s, `\u`, "\\u", -1)
	s = strings.Replace(s, `\U`, "\\U", -1)

	c.Value = s[1 : len(s)-1]
	return c
}

func (p *Parser) parseBool() Expression {
	c := &Const{}
	c.token = p.current()

	s := p.current().Literal
	p.consume("BOOL")

	var v bool = false

	if s[0] == 't' {
		v = true
	}

	c.Value = v
	return c
}

func (p *Parser) parseSpread() Expression {
	s := &Spread{}
	s.token = p.current()

	p.consume("TRIPLEDOT")
	s.Target = p.parseExpression(ARROW)

	return s
}

func (p *Parser) parseOverload(left Expression) Expression {
	o := &Overload{}
	o.token = p.current()

	_, isFunc := left.(*FunctionDef)
	if !isFunc {
		p.error("Left side of | operator must be a function definition", p.current())
		return nil
	}

	cases := []FunctionDef{*left.(*FunctionDef)}

	for p.current().Type == "PIPE" {
		p.consumeAny("NEWLINE")
		p.consume("PIPE")
		p.consumeAny("NEWLINE")
		c := p.parseExpression(OVERLOAD)
		if c != nil {
			f, isFunc := c.(*FunctionDef)
			if !isFunc {
				p.error("Right side of | operator must be a function definition", p.current())
				return nil
			}
			cases = append(cases, *f)
		}
	}

	o.Cases = cases
	return o
}
