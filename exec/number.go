package exec

import (
	"math"
	"strconv"
)

func NewInt(n int64) *OwlObj {
	return NewNumber(n)
}

func NewFloat(n float64) *OwlObj {
	return NewNumber(n)
}

const (
	INT     = 0
	FLOAT   = 1
	UNKNOWN = 2
)

func getRawType(obj *OwlObj) byte {
	switch obj.Raw.(type) {
	case int64:
		return INT
	case float64:
		return FLOAT
	default:
		return UNKNOWN
	}
}

func NewNumber(v interface{}) *OwlObj {
	n := NewOwlObj()

	n.Raw = v
	n.SetDeepAttr("add", NewCallBridge(numberAdd))
	n.SetDeepAttr("sub", NewCallBridge(numberSub))
	n.SetDeepAttr("mul", NewCallBridge(numberMul))
	n.SetDeepAttr("div", NewCallBridge(numberDiv))
	n.SetDeepAttr("pow", NewCallBridge(numberPow))
	n.SetDeepAttr("mod", NewCallBridge(numberMod))
	n.SetDeepAttr("neg", NewCallBridge(numberNeg))
	n.SetDeepAttr("inc", NewCallBridge(numberInc))
	n.SetDeepAttr("dec", NewCallBridge(numberDec))
	n.SetDeepAttr("eq", NewCallBridge(numberEq))
	n.SetDeepAttr("ne", NewCallBridge(numberNe))
	n.SetDeepAttr("lt", NewCallBridge(numberLt))
	n.SetDeepAttr("le", NewCallBridge(numberLe))
	n.SetDeepAttr("gt", NewCallBridge(numberGt))
	n.SetDeepAttr("ge", NewCallBridge(numberGe))
	n.SetDeepAttr("str", NewCallBridge(numberStr))
	n.DeleteDeepAttr("iter")
	n.DeleteDeepAttr("index")
	n.DeleteDeepAttr("setIndex")
	n.DeleteDeepAttr("has")

	return n
}

func numberAdd(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewInt(a.Raw.(int64) + b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewFloat(a.Raw.(float64) + b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewFloat(float64(a.Raw.(int64)) + b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewFloat(a.Raw.(float64) + float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberSub(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewInt(a.Raw.(int64) - b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewFloat(a.Raw.(float64) - b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewFloat(float64(a.Raw.(int64)) - b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewFloat(a.Raw.(float64) - float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberMul(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewInt(a.Raw.(int64) * b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewFloat(a.Raw.(float64) * b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewFloat(float64(a.Raw.(int64)) * b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewFloat(a.Raw.(float64) * float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberDiv(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	var aFloat, bFloat float64

	if aType == INT {
		aFloat = float64(a.Raw.(int64))
	} else if aType == FLOAT {
		aFloat = a.Raw.(float64)
	} else {
		return nil, false
	}

	if bType == INT {
		bFloat = float64(b.Raw.(int64))
	} else if bType == FLOAT {
		bFloat = b.Raw.(float64)
	} else {
		return nil, false
	}

	return NewFloat(aFloat / bFloat), true
}

func numberPow(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewInt(int64(math.Pow(float64(a.Raw.(int64)), float64(b.Raw.(int64))))), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewFloat(math.Pow(a.Raw.(float64), b.Raw.(float64))), true
	} else if aType == INT && bType == FLOAT {
		return NewFloat(math.Pow(float64(a.Raw.(int64)), b.Raw.(float64))), true
	} else if aType == FLOAT && bType == INT {
		return NewFloat(math.Pow(a.Raw.(float64), float64(b.Raw.(int64)))), true
	}

	return nil, false
}

func betterModInt(a int64, b int64) int64 {
	v := a % b

	if v < 0 {
		v += b
	}

	return v
}

func betterModFloat(a float64, b float64) float64 {
	v := math.Mod(a, b)

	if v < 0 {
		v += b
	}

	return v
}

func numberMod(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		v := betterModInt(a.Raw.(int64), b.Raw.(int64))
		return NewInt(v), true
	} else if aType == FLOAT && bType == FLOAT {
		v := betterModFloat(a.Raw.(float64), b.Raw.(float64))
		return NewFloat(v), true
	} else if aType == INT && bType == FLOAT {
		v := betterModFloat(float64(a.Raw.(int64)), b.Raw.(float64))
		return NewFloat(v), true
	} else if aType == FLOAT && bType == INT {
		v := betterModFloat(a.Raw.(float64), float64(b.Raw.(int64)))
		return NewFloat(v), true
	}

	return nil, false
}

func numberNeg(args []*OwlObj) (*OwlObj, bool) {
	aType := getRawType(args[0])

	if aType == INT {
		return NewInt(-args[0].Raw.(int64)), true
	} else if aType == FLOAT {
		return NewFloat(-args[0].Raw.(float64)), true
	}

	return nil, false
}

func numberInc(args []*OwlObj) (*OwlObj, bool) {
	aType := getRawType(args[0])

	if aType == INT {
		return NewInt(args[0].Raw.(int64) + 1), true
	} else if aType == FLOAT {
		return NewFloat(args[0].Raw.(float64) + 1), true
	}

	return nil, false
}

func numberDec(args []*OwlObj) (*OwlObj, bool) {
	aType := getRawType(args[0])

	if aType == INT {
		return NewInt(args[0].Raw.(int64) - 1), true
	} else if aType == FLOAT {
		return NewFloat(args[0].Raw.(float64) - 1), true
	}

	return nil, false
}

func numberEq(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewBool(a.Raw.(int64) == b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewBool(a.Raw.(float64) == b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewBool(float64(a.Raw.(int64)) == b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewBool(a.Raw.(float64) == float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberNe(args []*OwlObj) (*OwlObj, bool) {
	val, ok := numberEq(args)

	if ok {
		return val.Not()
	}

	return nil, false
}

func numberLt(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewBool(a.Raw.(int64) < b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewBool(a.Raw.(float64) < b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewBool(float64(a.Raw.(int64)) < b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewBool(a.Raw.(float64) < float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberLe(args []*OwlObj) (*OwlObj, bool) {
	a := args[1]
	b := args[2]
	aType := getRawType(a)
	bType := getRawType(b)

	if aType == INT && bType == INT {
		return NewBool(a.Raw.(int64) <= b.Raw.(int64)), true
	} else if aType == FLOAT && bType == FLOAT {
		return NewBool(a.Raw.(float64) <= b.Raw.(float64)), true
	} else if aType == INT && bType == FLOAT {
		return NewBool(float64(a.Raw.(int64)) <= b.Raw.(float64)), true
	} else if aType == FLOAT && bType == INT {
		return NewBool(a.Raw.(float64) <= float64(b.Raw.(int64))), true
	}

	return nil, false
}

func numberGt(args []*OwlObj) (*OwlObj, bool) {
	val, ok := numberLe(args)

	if ok {
		return val.Not()
	}

	return nil, false
}

func numberGe(args []*OwlObj) (*OwlObj, bool) {
	val, ok := numberLt(args)

	if ok {
		return val.Not()
	}

	return nil, false
}

func numberStr(args []*OwlObj) (*OwlObj, bool) {
	aType := getRawType(args[0])

	if aType == INT {
		return NewString(strconv.FormatInt(args[0].Raw.(int64), 10)), true
	} else if aType == FLOAT {
		return NewString(strconv.FormatFloat(args[0].Raw.(float64), 'f', -1, 64)), true
	}

	return nil, false
}
