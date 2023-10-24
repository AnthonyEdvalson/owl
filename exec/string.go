package exec

import (
	"regexp"
	"strings"
)

func NewString(v string) *OwlObj {
	s := NewOwlObj()

	s.SetDeepAttr("add", NewCallBridge(stringAdd))
	s.SetDeepAttr("eq", NewCallBridge(stringEq))
	s.SetDeepAttr("ne", NewCallBridge(stringNe))
	s.SetDeepAttr("gt", NewCallBridge(stringGt))
	s.SetDeepAttr("lt", NewCallBridge(stringLt))
	s.SetDeepAttr("ge", NewCallBridge(stringGe))
	s.SetDeepAttr("le", NewCallBridge(stringLe))
	s.SetDeepAttr("str", NewCallBridge(stringStr))
	s.SetDeepAttr("index", NewCallBridge(stringIndex))
	s.SetDeepAttr("has", NewCallBridge(stringHas))
	s.SetDeepAttr("slice", NewCallBridge(stringSlice))
	s.SetDeepAttr("iter", NewCallBridge(stringIter))

	s.SetAttr("Split", NewCallBridge(stringSplit))
	s.SetAttr("Len", NewCallBridge(stringLen))
	s.SetAttr("Replace", NewCallBridge(stringReplace))
	s.SetAttr("ReReplace", NewCallBridge(stringRegexReplace))
	s.SetAttr("ReIndex", NewCallBridge(stringRegexIndexOf))
	s.SetAttr("Index", NewCallBridge(stringIndexOf))
	s.SetAttr("Trim", NewCallBridge(stringTrim))
	s.Raw = v

	return s
}

func stringAdd(args []*OwlObj) (*OwlObj, bool) {
	return NewString(args[1].TrueStr() + args[2].TrueStr()), true
}

func stringEq(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() == args[2].TrueStr()), true
}

func stringNe(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() != args[2].TrueStr()), true
}

func stringLt(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() < args[2].TrueStr()), true
}

func stringGt(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() > args[2].TrueStr()), true
}

func stringLe(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() <= args[2].TrueStr()), true
}

func stringGe(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1].TrueStr() >= args[2].TrueStr()), true
}

func stringStr(args []*OwlObj) (*OwlObj, bool) {
	return args[0], true
}

func stringIndex(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	i, ok := args[1].TrueInt()

	if !ok {
		return nil, false
	}

	return NewString(string(s[int(i)])), ok
}

func stringHas(args []*OwlObj) (*OwlObj, bool) {
	this := args[0].TrueStr()
	sub := args[1].TrueStr()

	return NewBool(strings.Contains(this, sub)), true
}

func stringSplit(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	delim := args[1].TrueStr()

	var parts []string

	if len(args) == 2 {
		parts = strings.Split(s, delim)
	} else {
		n, ok := args[2].TrueInt()

		if !ok {
			return nil, false
		}

		parts = strings.SplitN(s, delim, int(n)+1)
	}

	objs := make([]*OwlObj, len(parts))
	for i, v := range parts {
		objs[i] = NewString(v)
	}

	return NewList(objs), true
}

func stringLen(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	return NewInt(int64(len(s))), true
}

func stringSlice(args []*OwlObj) (*OwlObj, bool) {
	list := args[0].TrueStr()
	start := args[1]
	end := args[2]

	var startIndex64, endIndex64 int64
	var startIndex, endIndex int = 0, len(list)
	var startOk, endOk bool

	if start != nil {
		startIndex64, startOk = start.TrueInt()
		if startOk {
			startIndex = mapIndex(startIndex64, len(list))
		}
	}

	if end != nil {
		endIndex64, endOk = end.TrueInt()
		if endOk {
			endIndex = mapIndex(endIndex64, len(list))
		}
	}

	if startIndex < 0 {
		startIndex = startIndex + len(list)
	}

	if endIndex > len(list) {
		endIndex = len(list)
	}

	if startIndex > endIndex {
		return NewString("test"), true
	}

	if startIndex >= len(list) || endIndex > len(list) {
		return nil, false
	}

	return NewString(list[startIndex:endIndex]), true
}

func stringIter(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()

	objs := make([]*OwlObj, len(s))
	for i, v := range s {
		objs[i] = NewString(string(v))
	}

	return NewList(objs), true
}

func stringReplace(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	old := args[1].TrueStr()
	new := args[2].TrueStr()

	return NewString(strings.Replace(s, old, new, -1)), true
}

func stringIndexOf(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	sub := args[1].TrueStr()

	return NewInt(int64(strings.Index(s, sub))), true
}

func stringRegexReplace(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	re := args[1].TrueStr()
	repl := args[2].TrueStr()

	exp := regexp.MustCompile(re)

	return NewString(exp.ReplaceAllString(s, repl)), true
}

func stringRegexIndexOf(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	re := args[1].TrueStr()

	exp := regexp.MustCompile(re)

	r := exp.FindIndex([]byte(s))
	var start int64 = -1
	var end int64 = -1
	if r != nil {
		start = int64(r[0])
		end = int64(r[1])
	}

	return NewList([]*OwlObj{NewInt(start), NewInt(end)}), true
}

func stringTrim(args []*OwlObj) (*OwlObj, bool) {
	s := args[0].TrueStr()
	cut := args[1].TrueStr()

	return NewString(strings.Trim(s, cut)), true
}
