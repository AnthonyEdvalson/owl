package exec

import (
	"fmt"
	"os"
	"owl/lexer"
	"owl/parser"
	"testing"
)

// THINGS TO TEST
// Statements
// - let
// - return
// - if
// - while
// - for
// - break
// - continue
// - expression statement
//
// Expressions
// - augassign
// - assign
// - number
// - string
// - bool
// - variable
// - function call
// - function def
// - list
// - map
// - set
// - unary
// - binary
// - if expression
// - index
// - comma
// - objects
// - attributes
// - increment / decrement
//
// Assignments
// - simple
// - comma
// - index
// - list
// - map
// - attributes
//
// Misc
// - recursion
// - self, simulated objects
// - closures
// - scope
// - pattern matching functions
// - imports

func testTruthy(t *testing.T, actual *OwlObj, expected bool) {
	truthy := actual.IsTruthy()

	if truthy != expected {
		t.Errorf("Expected %t but got %t", expected, truthy)
	}
}

func testInt(t *testing.T, actual *OwlObj, expected int64) {

	if actual == nil {
		t.Errorf("Expected %d but got nil", expected)
		return
	}

	result, ok := actual.TrueInt()

	if !ok {
		t.Errorf("Result is not Integer. got=%T (%+v)", actual.Raw, actual.Raw)
		return
	}

	if result != expected {
		t.Errorf("Result has wrong value. expected %d, got %d", expected, result)
	}
}

func testIntArray(t *testing.T, evaluated *OwlObj, expected []int64) {
	result, ok := evaluated.TrueList()

	if !ok {
		t.Errorf("Result is not Integer array. got=%T (%+v)", evaluated.Raw, evaluated.Raw)
	}

	if len(result) != len(expected) {
		t.Errorf("Result has wrong number of elements. expected %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		testInt(t, v, expected[i])
	}
}

func testFloat(t *testing.T, actual *OwlObj, expected float64) {
	result, ok := actual.TrueFloat()

	if !ok {
		t.Errorf("Result is not Float. got=%T (%+v)", actual.Raw, actual.Raw)
	}

	if result != expected {
		t.Errorf("Result has wrong value. expected %f, got %f", expected, result)
	}
}

func testString(t *testing.T, evaluated *OwlObj, expected string) {
	result, ok := evaluated.Raw.(string)

	if !ok {
		t.Errorf("Result is not String. got=%T (%+v)", evaluated.Raw, evaluated.Raw)
	}

	if result != expected {
		t.Errorf("Result has wrong value. expected %q, got %q", expected, result)
	}
}

func eval(s string) *OwlObj {
	l := lexer.NewLexer(s)
	t := l.Tokenize("exec_test.hoot")
	p := parser.NewParser(t)
	program := p.Parse()

	for _, error := range p.Errors {
		fmt.Println(error)
	}

	wd, _ := os.Getwd()

	e := NewTreeExecutor(wd)
	o := e.ExecProgram(program)
	return o
}

func TestNils(t *testing.T) {

	wd, _ := os.Getwd()
	e := NewTreeExecutor(wd)

	v := e.EvalExpression(nil)
	if v != nil {
		t.Errorf("Expected nil but got %+v", v)
	}
}

func TestBoolExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return true", true},
		{"return false", false},
		{"return not true", false},
		{"return not false", true},
		{"return not not true", true},
		{"return not not false", false},
		{"return true or false", true},
		{"return false or false", false},
		{"return false and true", false},
		{"return true and true", true},
		{"return true and not false", true},
		{"return false and not true", false},
		{"return not true or not true", false},
		{"return not false and not false", true},
		{"return not (true or false)", false},
		{"return not (false and true)", true},
		{"return not (true and not false)", false},
		{"return not (false and not true)", true},
		{"return not (false and [1][10])", true},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testTruthy(t, evaluated, tt.expected)
	}
}

func TestCompareExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return 1 < 2", true},
		{"return 1 > 2", false},
		{"return 1 < 1", false},
		{"return 1 > 1", false},
		{"return 1 == 1", true},
		{"return 1 != 1", false},
		{"return 1 == 2", false},
		{"return 1 != 2", true},
		{"return true == true", true},
		{"return false == false", true},
		{"return true == false", false},
		{"return true != false", true},
		{"return false != true", true},
		{"return (1 < 2) == true", true},
		{"return (1 < 2) == false", false},
		{"return (1 > 2) == true", false},
		{"return (1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testTruthy(t, evaluated, tt.expected)
	}
}

func TestNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 5", 5},
		{"return 10", 10},
		{"return -5", -5},
		{"return -10", -10},
		{"return 5 + 5 + 5 + 5 - 10", 10},
		{"return 2 * 2 * 2 * 2 * 2", 32},
		{"return 2 ** 4", 16},
		{"return 10 % 3", 1},
		{"return -10 % 3", 2},
		{"return -50 + 100 + -50", 0},
		{"return 5 * 2 + 10", 20},
		{"return 5 + 2 * 10", 25},
		{"return 20 + 2 * -10", 0},
		{"return 2 * (5 + 10)", 30},
		{"return 3 * 3 * 3 + 10", 37},
		{"return 3 * (3 * 3) + 10", 37},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}

	testsFloat := []struct {
		input    string
		expected float64
	}{
		{"return 5.0", 5.0},
		{"return 10.0", 10.0},
		{"return -5.0", -5.0},
		{"return -10.0", -10.0},
		{"return 5 + 5 + 5 + 5 - 10.0", 10.0},
		{"return 2 * 2.0 * 2.0 * 2 * 2", 32.0},
		{"return 2 ** 0.5", 1.4142135623730951},
		{"return 1.5 ** 2", 2.25},
		{"return -50 + 100 + -50.0", 0.0},
		{"return 5.0 * 2 + 10", 20.0},
		{"return 50 / 2 * 2 + 10", 60.0},
		{"return 5.0 + 2 * 10", 25.0},
		{"return (5 + 10 * 2 + 15 / 3) * 2 + -10", 50.0},
		{"return 20 + 2.0 * -10.0", 0.0},
		{"return 50 / 2.0 * 2.0 + 10.0", 60.0},
	}

	for _, tt := range testsFloat {
		evaluated := eval(tt.input)
		testFloat(t, evaluated, tt.expected)
	}
}

func TestStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return \"hello\"", "hello"},
		{"return 'hello'", "hello"},
		{"return \"hello\" + \" \" + \"world\"", "hello world"},
		{"return 'hello' + \" \" + 'world'", "hello world"},
		{"return 'a' < 'b' ? 'TRUE' : 'FALSE'", "TRUE"},
		{"return 'c' < 'b' ? 'TRUE' : 'FALSE'", "FALSE"},
		{"return 'a' > 'b' ? 'TRUE' : 'FALSE'", "FALSE"},
		{"return 'a' <= 'a' ? 'TRUE' : 'FALSE'", "TRUE"},
		{"return 'a' >= 'a' ? 'TRUE' : 'FALSE'", "TRUE"},
		{"return 'a' < 'a' ? 'TRUE' : 'FALSE'", "FALSE"},
		{"return 'a' > 'a' ? 'TRUE' : 'FALSE'", "FALSE"},
		{"return 'a,b,c'.Split(',').Join('.')", "a.b.c"},
		{"return 'a,b,c'.Split(',', 1).Join('.')", "a.b,c"},
		{"return 'abcd'[0]", "a"},
		{"return 'abcd'[1]", "b"},
		{"return 'abcd'[2]", "c"},
		{"return 'abcd'[2:]", "cd"},
		{"return 'abcd'[:-1]", "abc"},
		{"return 'abcd'[1:-1]", "bc"},
		{"return 'bac' has 'a' ? 'y' : 'n'", "y"},
		{"return 'bxc' has 'a' ? 'y' : 'n'", "n"},
		{"return '123'.Len() == 3 ? 'y' : 'n'", "y"},
		{"return 'aabbaxxb'.ReReplace('a(x*)b', '_${1}_')", "a__b_xx_"},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testString(t, evaluated, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return true ? 10 : 5", 10},
		{"return false ? 10 : 5", 5},
		{"return 1 < 2 ? 10 : 5", 10},
		{"return 1 > 2 ? 10 : 5", 5},
		{"return true ? (false ? 0 : 1) : 2", 1},
		{"return false ? (false ? 0 : 1) : 2", 2},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestLet(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5 \n return a", 5},
		{"let a = 5 \n return a + 1", 6},
		{"let a = 5 \n let b = 10 \n return a + b", 15},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestAssignmentExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"x = { a: [1, 2, 3] } \n x.a[1] = 5 \n return x.a[1]", 5},
		{"let a = 5 \n a = 2 \n return a", 2},
		{"let a = 5 \n a = 2 \n return a + 1", 3},
		{"let a = 5 \n let b = 10 \n a = 2 \n b = 20 \n return a + b", 22},
		{"x = { a: 3 } \n k = 'a' \n x[k] = 4 \n return x.a", 4},
		{"a = 1 \n b = 2 \n a, b = b, a \n return a", 2},
		{"a = 1 \n b = 2 \n a, b = b, a \n return b", 1},
		{"a = { v: 0 } \n b = { v: 1 } \n a.v, b.v = [1, 2] \n return a.v + b.v", 3},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"x = { a: 3 } \n return x.a", 3},
		{"x = null \n return x?.a == null ? 1 : 0", 1},
		{"x = { a: 3 } \n return x?.a == null ? 1 : 0", 0},
		{"x = null \n return x?.a.b.c == null ? 1 : 0", 1},
		{"x = { a: null } \n return x?.a?.b.c == null ? 1 : 0", 1},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestCommaExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"return 1, 2, 3", []int64{1, 2, 3}},
		{"return 1, 2, 3, 4, 5", []int64{1, 2, 3, 4, 5}},
		{"return 1 + 1, 2 + 2, 5 * 5", []int64{2, 4, 25}},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testIntArray(t, evaluated, tt.expected)
	}
}

func TestListExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"return (2, 3, 1, 8, 2).Sort()", []int64{1, 2, 2, 3, 8}},
		{"return (1, 2, 3, 4, 5)[0], (1, 2, 3, 4, 5)[-1]", []int64{1, 5}},
		{"return (1, 2, 3, 4, 5)[-5], (1, 2, 3, 4, 5)[4]", []int64{1, 5}},
		{"return (1, 2, 3, 4, 5)[:-2]", []int64{1, 2, 3}},
		{"return (1, 2, 3, 4, 5)[-2:]", []int64{4, 5}},
		{"return (1, 2, 3, 4, 5)[1:-2]", []int64{2, 3}},
		{"return (1 + 1, 2 + 2, 5 * 5)", []int64{2, 4, 25}},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testIntArray(t, evaluated, tt.expected)
	}
}

func TestListOps(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"return [1, 2, 3].Map(v => v + 1)", []int64{2, 3, 4}},
		{"return [1, 2, 3].Filter(v => v == 1)", []int64{1}},
		{"return [[1, 2, 3].Reduce((a, b) => a + b, 0)]", []int64{6}},
		{"return [1, 2, 3].FlatMap(v => [v, v + 1])", []int64{1, 2, 2, 3, 3, 4}},
		{"return (1, 2).Add(3)", []int64{1, 2, 3}},
		{"return (1, 2, 3, 4).Reverse()", []int64{4, 3, 2, 1}},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testIntArray(t, evaluated, tt.expected)
	}
}

func TestFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return ((x) => x)(5)", 5}, // TODO allow removal of parenthesis for args: (x => x)(5)
		{"return ((x) => { return x })(10)", 10},
		{"return (x => x)(10)", 10},
		{"return ((x, y) => { return x + y })(10, 5)", 15},
		{"return ((x) => x + 1)(5)", 6},
		{"return ((x, y) => x + y)(5, 10)", 15},
		{"return ((x) => x)((x) => x)(5)", 5},
		{"return ((x) => x)((x) => x + 1)(5)", 6},
		//{"return ((a) => (b) => a + b)(3)(4)", 7},  // Currying doesn't work yet, might need to mess with associativity
		{"return ((a, b) => a() + b())(() => 3, () => 4)", 7},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5 \n if (a < 10) { a = 12 } \nreturn a", 12},
		{"let a = 10 \n if (a < 10) { a = 12 } \nreturn a", 10},
		{"let a = 5 \n if (a < 10) { a = 12 } \nelse { a = 13 } \nreturn a", 12},
		{"let a = 10 \n if (a < 10) { a = 12 } else { a = 13 } \nreturn a", 13},
		{"let a = 5 \n if (a < 10) { a = 12 } else if (a < 20) { a = 13 } \nreturn a", 12},
		{"let a = 10 \n if (a < 10) { a = 12 } else if (a < 20) { a = 13 } \nreturn a", 13},
		{"let a = 20 \n if (a < 10) { a = 12 } else if (a < 20) \n{ a = 13 } \nreturn a", 20},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 0 \n while (a < 10) { a++ } \nreturn a", 10},
		{"let a = 20 \n while (a > 10) { a-- } \nreturn a", 10},
		{"let a = 0 \n while (a < 1000) { \na++\nif a > 10 {\nbreak\n} } \nreturn a", 11},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestForStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 0 \n for i in (1, 2, 3) { x += i } \nreturn x", 6},
		{"let x = 0 \n y = '' \n for k, v in { a: 2, b: 3, c: 4 } { \n x += v \n y += k \n } \n return x", 9},
		//{"let x = 0 \n for i in [] { x += 1 } \nreturn a", 0},
		{"let x = 0 \n for i in (2, 2, 2, 10, 50, 20, 10) { x = x + i\nif i > 20 {break} } \nreturn x", 66},
		{"let x = 0 \n for i in (4, 2, 5, 4, 9, 8) {\n\tif i % 2 == 0 { continue }\nx = x + i\n}\nreturn x", 14},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = {a: 5, b: 2} \n return a.a", 5},
		{"let a = {a: 5, b: 2} \n return a.b", 2},
		{"let a = {a: 5, b: 2} \n return a.a + a.b", 7},
		{"let a = {} \n a.v = 6 \n return a.v", 6},
		{"let a = {a: 5, b: 2} \n a.a = 10 \n return a.a", 10},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestDeepAttribute(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = {v: 4} \n let b = {v: 12} \n a::add = (x, y) => x.v + y.v \n return a + b", 16},
		{"let a = {v: 4} \n let b = {v: 12} \n a::add = (x, y) => x.v + y.v \n return a::add(a, b)", 16},
		{"let a = 4 \n a::neg = () => this + 1 \n return -a", 5},
		{"let a = true \n a::neg = () => !this \n return -a == false ? 1 : 0", 1},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = (x) => x + 1 \n return a(5)", 6},
		{"let a = (x) => x + 1 \n return a(5) + 1", 7},
		{"let a = (x) => x + 1 \n return a(a(5))", 7},
		{"let a = (x) => x + 1 \n let b = (x) => x + 2 \n return a(b(5))", 8},
		{"let a = null \n return a?(5) == null ? 1 : 0", 1},
		{"let a = (x) => x + 1 \n return a?(5) == null ? 1 : 0", 0},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestNull(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return ((x) => x == null ? 1 : 0)(null)", 1},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestPrint(t *testing.T) {
	tests := []string{
		"print 1",
		"print 1, 2, 3",
		"print {inc: (a) => a + 1}",
		"print 'string'",
		"print 1.0",
		"print [{\"company\": \"Snap\",\"highlights\": [\"tallico\", \"obsidian\", \"owl\"]}]",
	}

	for _, tt := range tests {
		eval(tt)
	}
}

func TestClosureScope(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a = 2 \n f = n => a + n \n a = 1 \n return f(3)", 5},
		{"incFactory = n => (v => v + n) \n inc2 = incFactory(2) \n inc5 = incFactory(5) \n return inc2(inc5(3))", 10},
		{"apply = (f, a) => (b => f(a, b)) \n sum = (a, b) => a + b \n inc = apply(sum, 1) \n a = 12 \n b = 35 \n f = 6 \n apply = 0 \n sum = 4 \n return inc(1)", 2},
		{"m = {} \n for i in [0, 1, 2] { \n m[i] = v => v + i \n } \n return m[1](4) + m[2](7)", 14},
		{"a = (f, v) => { l = [] \n while v > 0 { \n l.Add(1) \n v-- \n } \n f((_) => 0, 5) \n return l } \n return a(a, 1).Len()", 1},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

func TestSpread(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"f = (a, b, c) => a + b + c \n l = [1, 2, 3] \n return f(l)", 6},
		{"f = (a, b, c) => a + b + c \n l = [1, 2, 3] \n return f(2, ...l[1:])", 7},
		{"f = (a, b, c, d) => a + b + c + d \n return f(1, ...[2, 4], 2)", 9},
		{"f = (a, ...b) => a + b[0] + b[1] \n return f(1, 2, 3)", 6},
		{"a, ...b = 1, 2, 3, 4 \n return a + b[0] + b[1] + b[2]", 10},
		{"a, ...b = 1 \n return b.Len()", 0},
		{"...a = 1 \n return a.Len()", 1},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testInt(t, evaluated, tt.expected)
	}
}

// Tests are broken, libs are located basedon the directory things were executed from. Current design assumes that libs
// Are in the /lib directory adjacent to the executable. This is not the case when running tests.
/*func TestImport(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"import 'math' \n return math.tau", 6.28318530717958647692},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testFloat(t, evaluated, tt.expected)
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"import 'json'\nreturn json.ToObject('[123, 0]')", []int64{123, 0}},
		{"import 'json'\nreturn json.ToObject('[123, 0, 456]')", []int64{123, 0, 456}},
	}

	for _, tt := range tests {
		evaluated := eval(tt.input)
		testIntArray(t, evaluated, tt.expected)
	}
}

func TestHTTP(t *testing.T) {
	tests := []string{
		"import 'http' \n app = http.NewApp() \n app.Get('/project/', (req) => 'Hello!')",
	}

	for _, tt := range tests {
		eval(tt)
	}
}*/
