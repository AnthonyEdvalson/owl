package parser

import (
	"testing"

	"github.com/AnthonyEdvalson/owl/lexer"
)

func compareTrees(t *testing.T, expected string, actual Node) {

	if expected != actual.ToString() {
		t.Errorf("\nExpected %s\r\n\r\nGot %s", expected, actual.ToString())
	}
}

func parse(t *testing.T, s string) Node {
	l := lexer.NewLexer(s)
	tokens := l.Tokenize("parser_test.hoot")
	p := NewParser(tokens)
	prog := p.Parse()

	for _, error := range p.Errors {
		t.Errorf("%d:%d:  %s\r\n", error.Token.Line, error.Token.Column, error.Message)
	}

	return prog
}

func TestName(t *testing.T) {
	input := `x`

	expected := "x"

	actual := parse(t, input)
	compareTrees(t, expected, actual)
}

func TestConst(t *testing.T) {
	input := "12\nfalse\n14.5\n\"string\""

	expected := "12\nfalse\n14.5\n\"string\""

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestLetStatement(t *testing.T) {
	input := `let x = 5`
	expected := "let x = 5"

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestUnaryOp(t *testing.T) {
	input := "!true"

	expected := "(!true)"

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestBinaryOp(t *testing.T) {
	input := []string{
		"1 + 2",
		"3 - 4",
		"5 * 6",
		"7 / 8",
		"1 == 2",
		"1 != 2",
		"1 < 2",
		"1 > 2",
		"1 <= 2",
		"1 >= 2",
		"1 % 2",
		"x has 1",
		"3 ** 4",
		"3 ?? 4",
	}

	expected := []string{
		"(1 + 2)",
		"(3 - 4)",
		"(5 * 6)",
		"(7 / 8)",
		"(1 == 2)",
		"(1 != 2)",
		"(1 < 2)",
		"(1 > 2)",
		"(1 <= 2)",
		"(1 >= 2)",
		"(1 % 2)",
		"(x has 1)",
		"(3 ** 4)",
		"(3 ?? 4)",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestPrecedence(t *testing.T) {
	input := []string{
		"-a * b",
		"!-a",
		"a+b+c",
		"a+b-c",
		"a*b*c",
		"a*b/c",
		"a-b/c",
		"a+b*c+d/e-f",
		"3+4 \n -5 * 5",
		"5>4 and 3<4",
		"3+4*5 == 3*1+4*5",
		"3+4*5 == 0 or 3*1+4*5 > 1",
		"1+(2+3)+4",
		"(5 + 5)*2",
		"2/(5+5)",
		"-(5 + 5)",
		"!(true == true)",
	}

	expected := []string{
		"((-a) * b)",
		"(!(-a))",
		"((a + b) + c)",
		"((a + b) - c)",
		"((a * b) * c)",
		"((a * b) / c)",
		"(a - (b / c))",
		"(((a + (b * c)) + (d / e)) - f)",
		"(3 + 4)\n((-5) * 5)",
		"((5 > 4) and (3 < 4))",
		"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		"(((3 + (4 * 5)) == 0) or (((3 * 1) + (4 * 5)) > 1))",
		"((1 + (2 + 3)) + 4)",
		"((5 + 5) * 2)",
		"(2 / (5 + 5))",
		"(-(5 + 5))",
		"(!(true == true))",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestIfExpression(t *testing.T) {
	input := []string{
		"a ? b : c",
		"a ? b : c ? d : e",
		"a ? (b ? c + 1 : d) : e",
		"a + b ? c * d : e - f",
	}

	expected := []string{
		"(a ? b : c)",
		"((a ? b : c) ? d : e)",
		"(a ? (b ? (c + 1) : d) : e)",
		"((a + b) ? (c * d) : (e - f))",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestList(t *testing.T) {
	input := "let x = 1, [2], 3, [1, 2, 1 + 2, []]"

	expected := "let x = [1, [2], 3, [1, 2, (1 + 2), []]]"

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestForStatement(t *testing.T) {
	input := []string{
		"for x in 1,2,3 {\n\tlet y = x\n\tcontinue\n}",
		"let x = 0 \n for i in 2, 2, 2, 10, 50, 20, 10 { x += i\nif i > 20 {break} } \nreturn x",
		"let x = 0 \n for i in 4, 2, 5, 4, 9, 8 {\n\tif i % 2 == 0 { continue }\nx++\n}\nreturn x",
		"for data, html in projects {\n\tapp.Get(\"/project/\" + data.name, (req) => html)\n}",
	}

	expected := []string{
		"for x in [1, 2, 3] {\nlet y = x\ncontinue\n}",
		"let x = 0\nfor i in [2, 2, 2, 10, 50, 20, 10] {\nx += i\nif (i > 20) {\nbreak\n}\n}\nreturn x",
		"let x = 0\nfor i in [4, 2, 5, 4, 9, 8] {\nif ((i % 2) == 0) {\ncontinue\n}\nx++\n}\nreturn x",
		"for data, html in projects {\napp.Get([(\"/project/\" + data.name), (req) => {\nreturn html\n}])\n}",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestWhileStatement(t *testing.T) {
	input := "while x < 5 {\n\tlet y = x\nbreak\n}"

	expected := "while (x < 5) {\nlet y = x\nbreak\n}"

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestIfStatement(t *testing.T) {
	input := []string{
		"if x < 5 {\n\tlet y = x\n} \n\nelse { let y = x * 2 }",
		"if (a < 10) { a = 12 }",
		"if (a < 10) { a = 12 } \nelse { a = 13 }",
		"if (a < 10) { a = 12 } else { a = 13 }",
		"if (a < 10) { a = 12 } else if (a < 20) { a = 13 }",
		"if (a < 10) { a = 12 } else if (a < 20) { a = 13 } else { a = 14}",
		"if (a < 10)\n\n{\n\na = 12\n\n}\n\nelse if (a < 20)\n\n{\n\na = 13\n\n}\n\nelse\n\n{\n\na = 14\n\n}",
	}

	expected := []string{
		"if (x < 5) {\nlet y = x\n}\nelse {\nlet y = (x * 2)\n}",
		"if (a < 10) {\na = 12\n}",
		"if (a < 10) {\na = 12\n}\nelse {\na = 13\n}",
		"if (a < 10) {\na = 12\n}\nelse {\na = 13\n}",
		"if (a < 10) {\na = 12\n}\nelse {\nif (a < 20) {\na = 13\n}\n}",
		"if (a < 10) {\na = 12\n}\nelse {\nif (a < 20) {\na = 13\n}\nelse {\na = 14\n}\n}",
		"if (a < 10) {\na = 12\n}\nelse {\nif (a < 20) {\na = 13\n}\nelse {\na = 14\n}\n}",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestFunctionDef(t *testing.T) {
	input := []string{
		"() => {}",
		"(a) => {}",
		"a => {}",
		"() => 3",
		"(a) => a",
		"(a, b) => a + b",
		"(a) => {\n\tlet x = a\n\treturn x\n}",
		"(a, b) => {\n\tlet x = a\n\tlet y = b\n\treturn x + y\n}",
	}

	expected := []string{
		"(<>) => {\n}",
		"(a) => {\n}",
		"(a) => {\n}",
		"(<>) => {\nreturn 3\n}",
		"(a) => {\nreturn a\n}",
		"(a, b) => {\nreturn (a + b)\n}",
		"(a) => {\nlet x = a\nreturn x\n}",
		"(a, b) => {\nlet x = a\nlet y = b\nreturn (x + y)\n}",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestFunctionCall(t *testing.T) {
	input := []string{
		"f()",
		"f(1)",
		"f(1, 2)",
		"f(1 + 2)",
		"f(1 + 2, 3)",
		"(f + g)(1 + 2)",
		"(f + g)(1 + 2, 3)",
		"a + add(b * c) + d",
		"add(1, 2 * 3, add(6, 7 * 8))",
		"add(a + b + c * d / f + g)",
		"app.get('/', (req) => {\n    return home\n})",
	}

	expected := []string{
		"f()",
		"f(1)",
		"f([1, 2])",
		"f((1 + 2))",
		"f([(1 + 2), 3])",
		"(f + g)((1 + 2))",
		"(f + g)([(1 + 2), 3])",
		"((a + add((b * c))) + d)",
		"add([1, (2 * 3), add([6, (7 * 8)])])",
		"add((((a + b) + ((c * d) / f)) + g))",
		"app.get([\"/\", (req) => {\nreturn home\n}])",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestAttribute(t *testing.T) {
	input := []string{
		"a.b",
		"a.b.c",
		"(a + b).c[0]",
		"a[b].c",
		"a[b][c].d",
		"a::b",
		"a.b::c",
	}

	expected := []string{
		"a.b",
		"a.b.c",
		"(a + b).c[0]",
		"a[b].c",
		"a[b][c].d",
		"a::b",
		"a.b::c",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestAssignStatement(t *testing.T) {
	input := []string{
		"x = 1",
		"x, y = 1, 2",
		"x[1] = 1",
		"x.attr = 3",
		"x.attr[1].val, y = 2",
		"x::deepattr = 3",
		"x += 1",
		"a, b /= 15",
		"a, ...b = 1, 2, 3",
		//"{a, b} = {a: 1, b: 2}",
		//"{a, ...rest} = {a: 1, b: 2, c: 3}",
	}

	expected := []string{
		"x = 1",
		"x, y = [1, 2]",
		"x[1] = 1",
		"x.attr = 3",
		"x.attr[1].val, y = 2",
		"x::deepattr = 3",
		"x += 1",
		"a, b /= 15",
		"a, ...b = [1, 2, 3]",
		//"{a, b} = {a: 1, b: 2}",
		//"{a, ...rest} = {a: 1, b: 2, c: 3}",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestIndex(t *testing.T) {
	input := []string{
		"x[1]",
		"x[1][2]",

		"x[1 + 2]",
		"x[1][2 + 3]",

		"(true ? a : b)[1]",
		"(true ? a : b)[4 - p]",

		"x[1:2]",
		"x[-3:-1]",
		"x[:3]",
		"x[3:]",
	}

	expected := []string{
		"x[1]",
		"x[1][2]",

		"x[(1 + 2)]",
		"x[1][(2 + 3)]",

		"(true ? a : b)[1]",
		"(true ? a : b)[(4 - p)]",

		"x[1:2]",
		"x[(-3):(-1)]",
		"x[:3]",
		"x[3:]",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestIncDec(t *testing.T) {
	input := []string{
		"x++",
		"x--",
		"5 + x++ * 3",
		"4 * x-- - 3",
	}

	expected := []string{
		"x++",
		"x--",
		"(5 + (x++ * 3))",
		"((4 * x--) - 3)",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestMap(t *testing.T) {
	input := []string{
		"{}",
		"{a: 1}",
		"{a: 1, b: 2}",
		"{inc: (a) => a + 1}",
		"{a: 1 + 2, b: 2 + 3, c: 3 + 4}",
		"v = () => {\nreturn {\n    a: b.c,\n    \n    d: e.f\n}\n}",
	}

	expected := []string{
		"{\n}",
		"{\na: 1\n}",
		"{\na: 1,\nb: 2\n}",
		"{\ninc: (a) => {\nreturn (a + 1)\n}\n}",
		"{\na: (1 + 2),\nb: (2 + 3),\nc: (3 + 4)\n}",
		"v = (<>) => {\nreturn {\na: b.c,\nd: e.f\n}\n}",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestImport(t *testing.T) {
	input := []string{
		"import 'foo'",
	}

	expected := []string{
		"import foo",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestPrint(t *testing.T) {
	input := []string{
		"print 1",
		"print 1, 2",
		"print 5 * 4",
	}

	expected := []string{
		"print 1",
		"print [1, 2]",
		"print (5 * 4)",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestStringExpression(t *testing.T) {
	input := []string{
		"\"\"",
		"\"foo\"",
		"\"foo\\nbar\"",
		"\"foo\\tbar\"",
		"\"foo\\nbar\\nbaz\"",
		"\"foo\\tbar\\tbaz\"",
	}

	expected := []string{
		"\"\"",
		"\"foo\"",
		"\"foo\nbar\"",
		"\"foo\tbar\"",
		"\"foo\nbar\nbaz\"",
		"\"foo\tbar\tbaz\"",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestNull(t *testing.T) {
	input := []string{
		"null",
	}

	expected := []string{
		"null",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestNullAccess(t *testing.T) {
	input := []string{
		"a.b",
		"a?.b",
		"a?::b",
		"a?.b.c",
		"a?.b::c",
		"a.b?::c.d",
	}

	expected := []string{
		"a.b",
		"a?.b",
		"a?::b",
		"a?.b?.c",
		"a?.b?::c",
		"a.b?::c?.d",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestNullCall(t *testing.T) {
	input := "a.b?()"

	expected := "a.b?()"

	actual := parse(t, input)

	compareTrees(t, expected, actual)
}

func TestSpread(t *testing.T) {
	input := []string{
		"[...a]",
		"[a, ...b]",
	}

	expected := []string{
		"[...a]",
		"[a, ...b]",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestOverload(t *testing.T) {
	input := []string{
		"a = () => 5 \n | (a) => a",
		"a = () => 3 | (a) => a | (a, b) => a + b",
	}

	expected := []string{
		"a = <(<>) => {\nreturn 5\n} | (a) => {\nreturn a\n}>",
		"a = <(<>) => {\nreturn 3\n} | (a) => {\nreturn a\n} | (a, b) => {\nreturn (a + b)\n}>",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}

func TestFunctionCaseMatching(t *testing.T) {
	input := []string{
		"a = (0) => 3",
		"a = (0) => 3 | \n (a) => a",
		"a = (0, b) => 0 | (1, b) => b | (-1, b) => -b",
		"a = when a > 3, b == 4 (a, b) => 2 | (a, b) => a + b",
	}

	expected := []string{
		"a = when ($0 == 0)($0) => {\nreturn 3\n}",
		"a = <when ($0 == 0)($0) => {\nreturn 3\n} | (a) => {\nreturn a\n}>",
		"a = <when ($0 == 0)($0, b) => {\nreturn 0\n} | when ($0 == 1)($0, b) => {\nreturn b\n} | when ($0 == (-1))($0, b) => {\nreturn (-b)\n}>",
		"a = <when ((a > 3) and (b == 4))(a, b) => {\nreturn 2\n} | (a, b) => {\nreturn (a + b)\n}>",
	}

	for i := 0; i < len(input); i++ {
		compareTrees(t, expected[i], parse(t, input[i]))
	}
}
