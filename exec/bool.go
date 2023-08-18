package exec

var cache = [2]*OwlObj{}

func NewBool(v bool) *OwlObj {
	if v && cache[1] != nil {
		return cache[1]
	} else if !v && cache[0] != nil {
		return cache[0]
	}

	b := NewOwlObj()
	b.SetDeepAttr("and", NewCallBridge(boolAnd))
	b.SetDeepAttr("or", NewCallBridge(boolOr))
	b.SetDeepAttr("not", NewCallBridge(boolNot))
	b.SetDeepAttr("eq", NewCallBridge(boolEq))
	b.SetDeepAttr("ne", NewCallBridge(boolNe))
	b.SetDeepAttr("str", NewCallBridge(boolStr))
	b.DeleteDeepAttr("iter")
	b.DeleteDeepAttr("index")
	b.DeleteDeepAttr("setIndex")
	b.DeleteDeepAttr("has")

	//b.BridgeVal = v
	b.Raw = v

	if v {
		cache[1] = b
	} else {
		cache[0] = b
	}

	return b
}

func boolAnd(args []*OwlObj) (*OwlObj, bool) {
	v1 := args[1].IsTruthy()

	if !(v1) {
		return NewBool(false), true
	}

	v2o, ok := args[2].Call(nil)

	if !ok {
		return nil, false
	}

	v2 := v2o.IsTruthy()

	return NewBool(v2), true
}

func boolOr(args []*OwlObj) (*OwlObj, bool) {
	v1 := args[1].IsTruthy()

	if v1 {
		return NewBool(v1), true
	}

	v2o, ok := args[2].Call(nil)

	if !ok {
		return nil, false
	}

	v2 := v2o.IsTruthy()

	return NewBool(v2), true
}

func boolNot(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(!args[0].IsTruthy()), true
}

func boolEq(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].Raw == args[2].Raw), true
}

func boolNe(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].Raw != args[2].Raw), true
}

func boolStr(args []*OwlObj) (*OwlObj, bool) {
	if args[0].Raw.(bool) {
		return NewString("true"), true
	} else {
		return NewString("false"), true
	}
}
