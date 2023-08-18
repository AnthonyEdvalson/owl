package parser

import (
	"owl/lexer"
	"strings"
)

/*

start = Program(body []statement)

statement = Let(target assign, value expr)
          | AssignStatementtarget assign, value expr)
		  | For(target expr, iter expr, body []statement)
		  | While(test expr, body []statement)
		  | If(test expr, body []statement, else []statement)
		  | Throw(value expr)
		  | Try(body []statement, catch []statement, finally []statement)
		  | Expression(value expr)
		  | Return(value expr)
		  | Break()
		  | Continue()

expression = BinOp(left expr, op binop, right expr)
           | UnaryOp(op unaryop, right expr)
		   | FunctionCall(target expr, args []expr)
		   | IfExpression(test expr, iftrue expr, iffalse expr)
		   | Map(keys []expr, values []expr)
		   | Set(values []expr)
		   | List(values []expr)
		   | Constant(value constant)
		   | Function(args []string, body []statement)
		   | Attrubute(target expr, attr string)
		   | Index(target expr, index expr)
		   | Name(name string)

binop = Add | Sub | Mul | Div | FloorDiv | Mod | Pow | Eq | Neq | Lt | Gt | Le | Ge | And | Or
unaryop = Not | Negate
constant = Number | String | Boolean
*/

type Node interface {
	ToString() string
	Token() lexer.Token
}

type Program struct {
	Body []Statement
}

func (p *Program) ToString() string {
	var b strings.Builder
	
	for _, s := range p.Body {
		b.WriteString(s.ToString())
	}

	return strings.TrimSuffix(b.String(), "\n")
}

func (p *Program) Token() lexer.Token {
	return lexer.Token{Type: "PROGRAM"}
}