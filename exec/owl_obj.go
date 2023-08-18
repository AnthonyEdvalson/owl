package exec

import (
	"sort"
	"strings"
)

type OwlObj struct {
	Attr       map[string]*OwlObj
	DeepAttr   map[string]*OwlObj
	BridgeCall func(arg *OwlObj) (*OwlObj, bool)
	Raw        interface{}
	Bind       func(*OwlObj)
}

func NewOwlObj() *OwlObj {
	obj := OwlObj{}
	obj.Attr = make(map[string]*OwlObj)
	obj.DeepAttr = make(map[string]*OwlObj)
	obj.SetDeepAttr("str", NewCallBridge(objStr))
	obj.SetDeepAttr("index", NewCallBridge(objIndex))
	obj.SetDeepAttr("setIndex", NewCallBridge(objSetIndex))
	obj.SetDeepAttr("iter", NewCallBridge(objIter))
	obj.SetDeepAttr("has", NewCallBridge(objHas))
	obj.SetDeepAttr("coalesce", NewCallBridge(objCoalesce))

	return &obj
}

func objStr(args []*OwlObj) (*OwlObj, bool) {
	b := strings.Builder{}

	b.WriteString("{")

	for k, v := range args[0].Attr {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v.TrueStr())
		b.WriteString(", ")
	}

	b.WriteString("}")

	return NewString(b.String()), true
}

func objIndex(args []*OwlObj) (*OwlObj, bool) {
	this := args[0]
	index := args[1].TrueStr()

	return this.GetAttr(index)
}

func objHas(args []*OwlObj) (*OwlObj, bool) {
	this := args[0]
	index := args[1].TrueStr()

	_, ok := this.GetAttr(index)

	return NewBool(ok), true
}

func objCoalesce(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	if a.IsNullish() {
		b, bOk := args[2].Call(nil)
		if bOk {
			return b, true
		}
	}
	return a, true
}

func objSetIndex(args []*OwlObj) (*OwlObj, bool) {
	this := args[0]
	index := args[1].TrueStr()

	this.SetAttr(index, args[2])

	return NewNull(), true
}

func objIter(args []*OwlObj) (*OwlObj, bool) {
	items := make([]*OwlObj, 0)

	for k, v := range args[0].Attr {
		items = append(items, NewList([]*OwlObj{NewString(k), v}))
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Raw.([]*OwlObj)[0].Raw.(string) < items[j].Raw.([]*OwlObj)[0].Raw.(string)
	})

	return NewList(items), true
}

func (o *OwlObj) GetAttr(name string) (*OwlObj, bool) {
	v, ok := o.Attr[name]
	return v, ok
}

func (o *OwlObj) GetDeepAttr(name string) (*OwlObj, bool) {
	v, ok := o.DeepAttr[name]
	return v, ok
}

func (o *OwlObj) SetAttr(name string, value *OwlObj) {
	if value.Bind != nil {
		value.Bind(o)
	}

	o.Attr[name] = value
}

func (o *OwlObj) SetAttrs(values map[string]*OwlObj) {
	for k, v := range values {
		o.SetAttr(k, v)
	}
}

func (o *OwlObj) SetDeepAttr(name string, value *OwlObj) {
	if value.Bind != nil {
		value.Bind(o)
	}

	o.DeepAttr[name] = value
}

func (o *OwlObj) DeleteDeepAttr(name string) {
	if o.DeepAttr[name].Bind != nil {
		o.DeepAttr[name].Bind(nil)
	}

	delete(o.DeepAttr, name)
}

func (o *OwlObj) DeleteAttr(name string) {
	if o.Attr[name].Bind != nil {
		o.Attr[name].Bind(nil)
	}

	delete(o.Attr, name)
}

func (o *OwlObj) IsTruthy() bool {
	// TODO cast to bool and retun Raw.(bool)
	switch t := o.Raw.(type) {
	case bool:
		return t
	case int64:
		return t != 0
	case float64:
		return t != 0.0
	case string:
		return t != ""
	case []*OwlObj:
		return len(t) != 0
	case map[string]*OwlObj:
		return len(t) != 0 // Is this ever used?
	case *OwlObj:
		return t != nil
	case nil:
		return false
	default:
		return true
	}
}

func (o *OwlObj) IsNullish() bool {
	null, ok := o.GetDeepAttr("null")
	return ok && null.IsTruthy()
}

func (o *OwlObj) DeepCall(name string, arg *OwlObj) (*OwlObj, bool) {
	v, ok := o.GetDeepAttr(name)

	if !ok {
		return nil, false
	}

	return v.Call(arg)
}

func (o *OwlObj) Call(arg *OwlObj) (*OwlObj, bool) {
	if o.BridgeCall != nil {
		return o.BridgeCall(arg)
	}

	return o.DeepCall("call", arg)
}

func (this *OwlObj) Add(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("add", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Sub(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("sub", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Mul(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("mul", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Div(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("div", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Mod(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("mod", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Pow(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("pow", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Eq(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("eq", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Ne(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("ne", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Lt(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("lt", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Le(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("le", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Gt(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("gt", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Ge(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("ge", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Or(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("or", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) And(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("and", NewList([]*OwlObj{left, right}))
}

func (this *OwlObj) Coalesce(left *OwlObj, right *OwlObj) (*OwlObj, bool) {
	return this.DeepCall("coalesce", NewList([]*OwlObj{left, right}))
}

func (o *OwlObj) Not() (*OwlObj, bool) {
	return o.DeepCall("not", nil)
}

func (left *OwlObj) Has(right *OwlObj) (*OwlObj, bool) {
	return left.DeepCall("has", right)
}

func (o *OwlObj) Neg() (*OwlObj, bool) {
	return o.DeepCall("neg", nil)
}

func (o *OwlObj) Inc() (*OwlObj, bool) {
	return o.DeepCall("inc", nil)
}

func (o *OwlObj) Dec() (*OwlObj, bool) {
	return o.DeepCall("dec", nil)
}

func (o *OwlObj) Str() (*OwlObj, bool) {
	return o.DeepCall("str", nil)
}

func (o *OwlObj) Index(i *OwlObj) (*OwlObj, bool) {
	return o.DeepCall("index", i)
}

func (o *OwlObj) SetIndex(i *OwlObj, v *OwlObj) (*OwlObj, bool) {
	return o.DeepCall("setIndex", NewList([]*OwlObj{i, v}))
}

func (o *OwlObj) Slice(start *OwlObj, end *OwlObj) (*OwlObj, bool) {
	return o.DeepCall("slice", NewList([]*OwlObj{start, end}))
}

func (o *OwlObj) Iter() (*OwlObj, bool) {
	return o.DeepCall("iter", nil)
}

func (o *OwlObj) TrueStr() string {
	v, ok := o.Str()

	if !ok {
		return ""
	}

	return v.Raw.(string)
}

func (o *OwlObj) TrueInt() (int64, bool) {
	v, ok := o.Raw.(int64)
	return v, ok
}

func (o *OwlObj) TrueFloat() (float64, bool) {
	vf, ok := o.Raw.(float64)

	if !ok {
		vi, ok := o.Raw.(int64)

		if ok {
			vf = float64(vi)
		}
	}

	return vf, ok
}

func (o *OwlObj) TrueBool() (bool, bool) {
	v, ok := o.Raw.(bool)
	return v, ok
}

func (o *OwlObj) TrueList() ([]*OwlObj, bool) {
	v, ok := o.Raw.([]*OwlObj)
	return v, ok
}

func (o *OwlObj) AsList() ([]*OwlObj, bool) {
	// Gets a list of objects from either TrueList() or Iter()
	list, ok := o.TrueList()

	if !ok {
		var iter *OwlObj
		iter, ok = o.Iter()
		if ok {
			list, ok = iter.TrueList()
		}
	}

	return list, ok
}
