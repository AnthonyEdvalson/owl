package exec

import (
	"strings"

	"github.com/AnthonyEdvalson/owl/parser"
)

type FuncData struct {
	Body      []parser.Statement
	Exec      *TreeExecutor
	Arg       parser.Assign
	This      *OwlObj
	Env       *Frame
	Condition parser.Expression
	Else      *FuncData
}

func NewFunc(exec *TreeExecutor, def *parser.FunctionDef, frame *Frame) *OwlObj {
	f := NewOwlObj()
	f.SetDeepAttr("str", NewCallBridge(funcStr))

	data := funcDefToData(def, exec, frame)
	f.Bind = func(this *OwlObj) { data.This = this }
	f.Raw = data
	f.BridgeCall = func(a *OwlObj) (*OwlObj, bool) { return funcCall(data, a) }

	return f
}

func funcDefToData(def *parser.FunctionDef, exec *TreeExecutor, frame *Frame) *FuncData {
	ctx := Frame{}
	for k, v := range *frame {
		ctx[k] = v
	}

	var elseData *FuncData
	if def.Else != nil {
		elseData = funcDefToData(def.Else, exec, frame)
	}

	data := &FuncData{
		Body:      def.Body,
		Exec:      exec,
		Arg:       def.Arg,
		This:      nil,
		Env:       &ctx,
		Condition: def.Condition,
		Else:      elseData,
	}
	return data
}

func funcCall(f *FuncData, arg *OwlObj) (*OwlObj, bool) {
	t := f.Exec
	data := f
	for data != nil {
		t.pushFrameContext(data.Env)
		t.pushFrame()
		t.Assign(data.Arg, arg)
		t.set("this", data.This)
		if data.Condition == nil || data.Exec.EvalExpression(data.Condition).IsTruthy() {
			state := data.Exec.ExecBlock(data.Body)
			t.popFrame()
			t.popFrame()
			return state.Return, true
		}
		t.popFrame()
		t.popFrame()
		data = data.Else
	}
	return NewString("Unable to find a matching overload"), false
}

func funcStr(args []*OwlObj) (*OwlObj, bool) {
	f := args[0].Raw.(*FuncData)

	b := strings.Builder{}

	b.WriteString("{")

	for _, stmt := range f.Body {
		b.WriteString(stmt.ToString())
	}

	b.WriteString("}")

	return NewString(b.String()), true
}
