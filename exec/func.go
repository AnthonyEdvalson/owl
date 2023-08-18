package exec

import (
	"owl/parser"
	"strings"
)

type FuncData struct {
	Body []parser.Statement
	Exec *TreeExecutor
	Arg parser.Assign
	This *OwlObj
	Env *Frame
}

func NewFunc(exec *TreeExecutor, argAssign parser.Assign, body []parser.Statement, frame *Frame) *OwlObj {
	f := NewOwlObj()
	f.SetDeepAttr("str",  NewCallBridge(funcStr))

	ctx := Frame{}
	for k, v := range *frame {
		ctx[k] = v
	}

	data := &FuncData{body, exec, argAssign, nil, &ctx}
	f.Bind = func(this *OwlObj) {data.This = this}
	f.Raw = data

	f.BridgeCall = func (a *OwlObj) (*OwlObj, bool) { return funcCall(data, a) }

	return f
}

func funcCall(f *FuncData, arg *OwlObj) (*OwlObj, bool) {
	f.Exec.pushFrameContext(f.Env)
	f.Exec.pushFrame()
	f.Exec.Assign(f.Arg, arg)
	f.Exec.set("this", f.This)
	state := f.Exec.ExecBlock(f.Body)
	f.Exec.popFrame()
	f.Exec.popFrame()
	return state.Return, true
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