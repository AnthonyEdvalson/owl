package exec

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AnthonyEdvalson/owl/lexer"
	"github.com/AnthonyEdvalson/owl/parser"
)

type Frame map[string]*OwlObj

type TreeExecutor struct {
	Frames      []Frame
	currentPath string
}

func NewTreeExecutor(path string) *TreeExecutor {
	t := &TreeExecutor{
		Frames:      []Frame{},
		currentPath: path,
	}

	return t
}

type RunState struct {
	State  int
	Return *OwlObj
}

const (
	RUN = iota
	RETURN
	BREAK
	CONTINUE
)

func (t *TreeExecutor) set(name string, value *OwlObj) {
	(*t.bottomFrame())[name] = value
}

func (t *TreeExecutor) get(name string, token lexer.Token) *OwlObj {
	for i := len(t.Frames) - 1; i >= 0; i-- {
		f := t.Frames[i]
		if value, ok := f[name]; ok {
			return value
		}
	}

	t.panic("Unable to find variable '"+name+"'", token)
	return nil
}

func (t *TreeExecutor) panic(msg string, token lexer.Token) {
	panic(fmt.Sprintf("%s:%d:%d:  %s\r\n", token.File, token.Line, token.Column, msg))
}

func (t *TreeExecutor) resetStack() {
	t.Frames = []Frame{}
	t.pushFrame()
	t.set("this", NewNull())
}

func (t *TreeExecutor) pushFrame() {
	t.Frames = append(t.Frames, Frame{})
}

func (t *TreeExecutor) pushFrameContext(frame *Frame) {
	t.Frames = append(t.Frames, *frame)
}

func (t *TreeExecutor) bottomFrame() *Frame {
	return &t.Frames[len(t.Frames)-1]
}

func (t *TreeExecutor) popFrame() {
	t.Frames = t.Frames[:len(t.Frames)-1]
}

func (t *TreeExecutor) Assign(assign parser.Assign, value *OwlObj) {
	switch a := assign.(type) {
	case *parser.AssignName:
		t.set(a.Name, value)
	case *parser.AssignList:
		values, ok := value.AsList()
		if !ok {
			values = []*OwlObj{value}
		}

		// Find spread, if it exists
		spreadIndex := -1
		beforeSpread := []parser.Assign{}
		afterSpread := []parser.Assign{}
		for i, part := range a.Parts {
			_, isSpread := part.(*parser.AssignSpread)
			if isSpread {
				if spreadIndex != -1 {
					t.panic("Multiple spreads in assignment", a.Token())
				}
				spreadIndex = i
			} else if spreadIndex == -1 {
				beforeSpread = append(beforeSpread, part)
			} else {
				afterSpread = append(afterSpread, part)
			}
		}

		// Check that the number of values is correct
		if spreadIndex == -1 && len(values) != len(a.Parts) {
			t.panic("Expected "+fmt.Sprintf("%d", len(a.Parts))+" values, got "+fmt.Sprintf("%d", len(values)), a.Token())
		} else if spreadIndex != -1 && len(values) < len(beforeSpread)+len(afterSpread) {
			t.panic("Expected at least "+fmt.Sprintf("%d", len(beforeSpread)+len(afterSpread))+" values, got "+fmt.Sprintf("%d", len(values)), a.Token())
		}

		// Assign values
		for i, part := range beforeSpread {
			t.Assign(part, values[i])
		}
		if spreadIndex != -1 {
			t.Assign(a.Parts[spreadIndex], NewList(values[len(beforeSpread):len(values)-len(afterSpread)]))
		}
		for i, part := range afterSpread {
			t.Assign(part, values[len(values)-len(afterSpread)+i])
		}
	case *parser.AssignAttribute:
		target := t.EvalExpression(a.Target)
		if a.IsCoalesce && target.IsNullish() {
			break
		}
		if a.IsDeep {
			target.SetDeepAttr(a.Attribute, value)
		} else {
			target.SetAttr(a.Attribute, value)
		}
	case *parser.AssignIndex:
		target := t.getFromAssign(a.Target)
		index := t.EvalExpression(a.Index)
		target.SetIndex(index, value)
	/*TODO case *parser.AssignMap:*/
	case *parser.AssignSpread:
		// Convert value to list if it isn't already
		_, isList := value.AsList()
		if isList {
			t.Assign(a.Target, value)
		} else {
			t.Assign(a.Target, NewList([]*OwlObj{value}))
		}
	case *parser.AssignNull:
		if value != nil {
			t.panic("Expected nil, got "+value.TrueStr(), a.Token())
		}
	default:
		n, ok := assign.(parser.Node)

		if !ok {
			t.panic("Unknown assign type", lexer.Token{})
		}

		t.panic("Unknown assign type "+n.ToString(), n.Token())
	}
}

func (t *TreeExecutor) getFromAssign(assign parser.Assign) *OwlObj {
	switch a := assign.(type) {
	case *parser.AssignName:
		return t.get(a.Name, a.Token())
	case *parser.AssignAttribute:
		var val *OwlObj
		ok := false

		if a.IsDeep {
			val, ok = t.EvalExpression(a.Target).GetDeepAttr(a.Attribute)
		} else {
			val, ok = t.EvalExpression(a.Target).GetAttr(a.Attribute)
		}

		if !ok {
			t.panic("Unable to find attribute "+a.Attribute, a.Token())
		}
		return val
	default:
		t.panic("Unable to get "+a.ToString()+" in assignment", a.Token())
		return nil
	}
}

func ExecuteFile(path string) (*OwlObj, *TreeExecutor, []parser.ParserError) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return Execute(string(bytes), path, map[string]*OwlObj{})
}

func Execute(s string, path string, globals map[string]*OwlObj) (*OwlObj, *TreeExecutor, []parser.ParserError) {
	l := lexer.NewLexer(s)
	tok := l.Tokenize(filepath.Base(path))
	p := parser.NewParser(tok)
	program := p.Parse()

	if len(p.Errors) > 0 {
		return nil, nil, p.Errors
	}

	e := NewTreeExecutor(path)
	for k, v := range globals {
		e.set(k, v)
	}
	o := e.ExecProgram(program)
	return o, e, nil
}

func (t *TreeExecutor) ExecProgram(program *parser.Program) *OwlObj {
	t.resetStack()
	s := t.ExecBlock(program.Body)

	return s.Return
}

func (t *TreeExecutor) ExecBlock(block []parser.Statement) RunState {
	for _, stmt := range block {
		s := t.execStatement(stmt)

		if s.State != RUN {
			return s
		}
	}

	return RunState{RUN, nil}
}

func (t *TreeExecutor) execStatement(stmt parser.Statement) RunState {
	switch stmt := stmt.(type) {
	case *parser.ExpressionStatement:
		t.EvalExpression(stmt.Value)
		return RunState{RUN, nil}
	case *parser.Let:
		t.execLetStatement(stmt)
		return RunState{RUN, nil}
	case *parser.Return:
		value := t.execReturnStatement(stmt)
		return RunState{RETURN, value}
	case *parser.If:
		return t.execIfStatement(stmt)
	case *parser.While:
		return t.execWhileStatement(stmt)
	case *parser.For:
		return t.execForStatement(stmt)
	case *parser.Break:
		return RunState{BREAK, nil}
	case *parser.Continue:
		return RunState{CONTINUE, nil}
	case *parser.Import:
		return t.execImportStatement(stmt)
	case *parser.Print:
		return t.execPrintStatement(stmt)
	default:
		t.panic("Unable to process statement '"+stmt.ToString()+"'", stmt.Token())
		return RunState{-1, nil}
	}
}

func (t *TreeExecutor) execLetStatement(l *parser.Let) {
	val := t.EvalExpression(l.Value)
	t.Assign(l.Target, val)
}

func (t *TreeExecutor) execReturnStatement(stmt *parser.Return) *OwlObj {
	v := t.EvalExpression(stmt.Value)
	return v
}

func (t *TreeExecutor) execForStatement(f *parser.For) RunState {
	iter := t.EvalExpression(f.Iter)

	list, ok := iter.AsList()

	if !ok {
		t.panic("For loop iter is not a list", f.Token())
	}

	for _, item := range list {
		t.Assign(f.Target, item)
		state := t.ExecBlock(f.Body)

		switch state.State {
		case CONTINUE:
			continue
		case BREAK:
			return RunState{RUN, nil} // TODO add broken block? Similar to Python for else
		case RETURN:
			return state
		}
	}

	return RunState{RUN, nil}
}

func (t *TreeExecutor) execIfStatement(i *parser.If) RunState {
	if t.EvalExpression(i.Test).IsTruthy() {
		return t.ExecBlock(i.Body)
	} else if i.Else != nil {
		return t.ExecBlock(i.Else)
	}

	return RunState{RUN, nil}
}

func (t *TreeExecutor) execWhileStatement(w *parser.While) RunState {
	for t.EvalExpression(w.Test).IsTruthy() {
		state := t.ExecBlock(w.Body)

		switch state.State {
		case CONTINUE:
			continue
		case BREAK:
			return RunState{RUN, nil}
		case RETURN:
			return state
		}
	}

	return RunState{RUN, nil}
}

// execImportStatement executes an import AST node. Import names have the
// following formats:
//  1. A relative path starting with ./ or ../ will be resolved relative to the
//     current file.
//  2. An absolute path (starting with /) will be resolved to the given path.
//  3. A name beginning with no slashes or dots will be resolved to the
//     standard library.
func (t *TreeExecutor) execImportStatement(i *parser.Import) RunState {
	module, alias := NewModule(i.Name, t.currentPath)

	t.set(alias, module)

	return RunState{RUN, nil}
}

func (t *TreeExecutor) execPrintStatement(i *parser.Print) RunState {
	v := t.EvalExpression(i.Value)

	//TODO may need to replace \n with \r\n
	fmt.Println(v.TrueStr())

	return RunState{RUN, nil}
}

func (t *TreeExecutor) EvalExpression(expr parser.Expression) *OwlObj {
	switch expr := expr.(type) {
	case *parser.Const:
		return t.evalConst(expr)
	case *parser.BinOp:
		return t.evalBinOp(expr)
	case *parser.UnaryOp:
		return t.evalUnaryOp(expr)
	case *parser.FunctionDef:
		return t.evalFunctionDef(expr)
	case *parser.Overload:
		return t.evalOverload(expr)
	case *parser.FunctionCall:
		return t.evalFunctionCall(expr)
	case *parser.IfExpression:
		return t.evalIfExpression(expr)
	case *parser.AssignExpression:
		return t.evalAssignExpression(expr)
	case *parser.Map:
		return t.evalMap(expr)
	/*case *parser.Set:
	return t.evalSet(expr)*/
	case *parser.Index:
		return t.evalIndex(expr)
	case *parser.Slice:
		return t.evalSlice(expr)
	case *parser.Attribute:
		return t.evalAttribute(expr)
	case *parser.Name:
		return t.evalName(expr)
	case *parser.List:
		return t.evalCommaOp(expr)
	case *parser.IncDec:
		return t.evalIncDec(expr)
	case *parser.Null:
		return t.evalNull(expr)
	case nil:
		return nil
	}

	t.panic("Unable to evaluate expression '"+expr.ToString()+"'", expr.Token())
	return nil
}

func (t *TreeExecutor) evalConst(c *parser.Const) *OwlObj {
	switch c.Value.(type) {
	case string:
		return NewString(c.Value.(string))
	case int64:
		return NewInt(c.Value.(int64))
	case float64:
		return NewFloat(c.Value.(float64))
	case bool:
		return NewBool(c.Value.(bool))
	}

	t.panic("Unable to evaluate constant '"+c.ToString()+"'", c.Token())
	return nil
}

func (t *TreeExecutor) evalNull(c *parser.Null) *OwlObj {
	return NewNull()
}

func (t *TreeExecutor) evalName(n *parser.Name) *OwlObj {
	return t.get(n.Name, n.Token())
}

func (t *TreeExecutor) evalCommaOp(c *parser.List) *OwlObj {
	var values []*OwlObj

	for _, expr := range c.Parts {
		spread, isSpread := expr.(*parser.Spread)
		if isSpread {
			value := t.EvalExpression(spread.Target)
			valueList, ok := value.AsList()
			if !ok {
				t.panic("Spread value is not a list", spread.Token())
			}
			values = append(values, valueList...)
		} else {
			values = append(values, t.EvalExpression(expr))
		}
	}

	return NewList(values)
}

func (t *TreeExecutor) evalAssignExpression(a *parser.AssignExpression) *OwlObj {
	val := t.EvalExpression(a.Value)
	ok := false

	// TODO: this can break if RHS has not implemented the operation
	// Would be better for the parser to replace a += b with a = a + b, then it can use BinOp implementation.
	switch a.Op {
	case "=":
		ok = true
	case "+=":
		val, ok = t.getFromAssign(a.Target).Add(t.getFromAssign(a.Target), val)
	case "-=":
		val, ok = t.getFromAssign(a.Target).Sub(t.getFromAssign(a.Target), val)
	case "*=":
		val, ok = t.getFromAssign(a.Target).Mul(t.getFromAssign(a.Target), val)
	case "/=":
		val, ok = t.getFromAssign(a.Target).Div(t.getFromAssign(a.Target), val)
	case "%=":
		val, ok = t.getFromAssign(a.Target).Mod(t.getFromAssign(a.Target), val)
	case "&=":
		val, ok = t.getFromAssign(a.Target).And(t.getFromAssign(a.Target), val)
	case "|=":
		val, ok = t.getFromAssign(a.Target).Or(t.getFromAssign(a.Target), val)
	}

	if !ok {
		t.panic("Unable to assign value '"+val.TrueStr()+"' to '"+a.Target.ToString()+"'", a.Token())
	}

	t.Assign(a.Target, val)
	return val
}

func (t *TreeExecutor) evalBinOp(b *parser.BinOp) *OwlObj {
	left := t.EvalExpression(b.Left)
	var val *OwlObj
	ok := false

	if b.Op == "and" || b.Op == "or" || b.Op == "??" {
		lazyRight := NewCallBridge(func(_ []*OwlObj) (*OwlObj, bool) { return t.EvalExpression(b.Right), true })

		if b.Op == "and" {
			val, ok = left.And(left, lazyRight)
		} else if b.Op == "or" {
			val, ok = left.Or(left, lazyRight)
		} else if b.Op == "??" {
			val, ok = left.Coalesce(left, lazyRight)
		}
	} else {
		right := t.EvalExpression(b.Right)

		switch b.Op {
		case "+":
			val, ok = left.Add(left, right)
			if !ok {
				val, ok = right.Add(left, right)
			}
		case "-":
			val, ok = left.Sub(left, right)
		case "*":
			val, ok = left.Mul(left, right)
			if !ok {
				val, ok = right.Mul(left, right)
			}
		case "/":
			val, ok = left.Div(left, right)
		case "**":
			val, ok = left.Pow(left, right)
		case "%":
			val, ok = left.Mod(left, right)
		case "==":
			val, ok = left.Eq(left, right)
			if !ok {
				val, ok = right.Eq(left, right)
			}
		case "!=":
			val, ok = left.Ne(left, right)
			if !ok {
				val, ok = right.Ne(left, right)
			}
		case ">":
			val, ok = left.Gt(left, right)
			if !ok {
				val, ok = right.Le(left, right)
			}
		case "<":
			val, ok = left.Lt(left, right)
			if !ok {
				val, ok = right.Ge(left, right)
			}
		case ">=":
			val, ok = left.Ge(left, right)
			if !ok {
				val, ok = right.Lt(left, right)
			}
		case "<=":
			val, ok = left.Le(left, right)
			if !ok {
				val, ok = right.Gt(left, right)
			}
		case "has":
			val, ok = left.Has(right)
		default:
			t.panic("Unknown binary operator '"+b.ToString()+"'", b.Token())
			return nil
		}

	}

	if !ok {
		t.panic("Unable to evaluate binary operator '"+b.ToString()+"'", b.Token())
	}

	return val
}

func (t *TreeExecutor) evalUnaryOp(u *parser.UnaryOp) *OwlObj {
	v := t.EvalExpression(u.Value)
	var val *OwlObj
	ok := false

	switch u.Op {
	case "!", "not":
		val, ok = v.Not()
	case "-":
		val, ok = v.Neg()
	default:
		t.panic("Unknown unary operator '"+u.ToString()+"'", u.Token())
		return nil
	}

	if !ok {
		t.panic("Unable to evaluate unary operator '"+u.ToString()+"'", u.Token())
	}

	return val
}

func (t *TreeExecutor) evalIncDec(i *parser.IncDec) *OwlObj {
	assign := i.Target
	v := t.getFromAssign(assign)
	var newV *OwlObj
	ok := false

	if i.Op == "++" {
		newV, ok = v.Inc()
	} else {
		newV, ok = v.Dec()
	}

	if !ok {
		t.panic("Unable to evaluate increment/decrement '"+i.ToString()+"'", i.Token())
	}

	t.Assign(assign, newV)

	return newV
}

func (t *TreeExecutor) evalAttribute(a *parser.Attribute) *OwlObj {
	target := t.EvalExpression(a.Target)
	var val *OwlObj
	var ok bool

	if a.IsCoalesce && target.IsNullish() {
		val = target
		ok = true
	} else {
		if a.IsDeep {
			val, ok = target.GetDeepAttr(a.Attribute)
		} else {
			val, ok = target.GetAttr(a.Attribute)
		}
	}

	if !ok {
		t.panic("Unable to evaluate attribute '"+a.ToString()+"'", a.Token())
	}

	return val
}

func (t *TreeExecutor) evalIndex(i *parser.Index) *OwlObj {
	target := t.EvalExpression(i.Target)
	index := t.EvalExpression(i.Index)
	val, ok := target.Index(index)

	if !ok {
		t.panic("Unable to evaluate index '"+i.ToString()+"', "+i.Target.ToString()+" does not have index "+index.TrueStr(), i.Token())
	}

	return val
}

func (t *TreeExecutor) evalSlice(i *parser.Slice) *OwlObj {
	target := t.EvalExpression(i.Target)
	var start, end *OwlObj

	if i.Start != nil {
		start = t.EvalExpression(i.Start)
	}

	if i.End != nil {
		end = t.EvalExpression(i.End)
	}

	val, ok := target.Slice(start, end)

	if !ok {
		t.panic("Unable to evaluate slice '"+i.ToString()+"'", i.Token())
	}

	return val
}

func (t *TreeExecutor) evalFunctionDef(f *parser.FunctionDef) *OwlObj {
	fn := NewFunc(t, f, t.bottomFrame())
	return fn
}

func (t *TreeExecutor) evalOverload(o *parser.Overload) *OwlObj {
	for i := 0; i < len(o.Cases)-1; i++ {
		o.Cases[i].Else = &o.Cases[i+1]
	}
	return t.evalFunctionDef(&o.Cases[0])
}

func (t *TreeExecutor) evalFunctionCall(c *parser.FunctionCall) *OwlObj {
	fn := t.EvalExpression(c.Target)
	if c.IsCoalesce && fn.IsNullish() {
		return fn
	}
	arg := t.EvalExpression(c.Arg)
	val, ok := fn.Call(arg)

	if !ok {
		msg := val.TrueStr()
		t.panic("Unable to evaluate function call '"+c.ToString()+"', "+msg, c.Token())
	}

	return val
}

func (t *TreeExecutor) evalMap(m *parser.Map) *OwlObj {
	o := NewOwlObj()

	for i, attr := range m.Keys {
		val := t.EvalExpression(m.Values[i])
		o.SetAttr(attr, val)
	}

	return o
}

func (t *TreeExecutor) evalIfExpression(i *parser.IfExpression) *OwlObj {
	if t.EvalExpression(i.Test).IsTruthy() {
		return t.EvalExpression(i.IfTrue)
	} else {
		return t.EvalExpression(i.IfFalse)
	}
}
