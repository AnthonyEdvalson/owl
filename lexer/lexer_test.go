package lexer

import (
	"fmt"
	"testing"
)

type ShortToken struct {
	Type    TokenType
	Literal string
}

func compareTokens(t *testing.T, expected []Token, actual []Token) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(actual))
		fmt.Println(actual)
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("Expected %#v, got %#v", expected[i], actual[i])
			fmt.Println(actual)
		}
	}
}

func compareShortTokens(t *testing.T, expected []ShortToken, actual []Token) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(actual))
		fmt.Println(actual)
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].Literal != actual[i].Literal || expected[i].Type != actual[i].Type {
			t.Errorf("Expected %#v, got %#v", expected[i], actual[i])
			fmt.Println(actual)
		}
	}
}

func tokenize(s string) []Token {
	l := NewLexer(s)
	return l.Tokenize("lexer_test.hoot")
}

func TestLet(t *testing.T) {
	tokens := tokenize("let x = 3")

	expected := []Token{
		{"LET", "let", 1, 1, "lexer_test.hoot"},
		{"NAME", "x", 1, 5, "lexer_test.hoot"},
		{"ASSIGN", "=", 1, 7, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 9, "lexer_test.hoot"},
		{"EOF", "", 1, 10, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestExpression(t *testing.T) {
	tokens := tokenize("inc = 3 + 4 % 3")

	expected := []Token{
		{"NAME", "inc", 1, 1, "lexer_test.hoot"},
		{"ASSIGN", "=", 1, 5, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 7, "lexer_test.hoot"},
		{"PLUS", "+", 1, 9, "lexer_test.hoot"},
		{"NUMBER", "4", 1, 11, "lexer_test.hoot"},
		{"PERCENT", "%", 1, 13, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 15, "lexer_test.hoot"},
		{"EOF", "", 1, 16, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestSimpleFunction(t *testing.T) {
	tokens := tokenize("x => { x + 3 }")

	expected := []Token{
		{"NAME", "x", 1, 1, "lexer_test.hoot"},
		{"ARROW", "=>", 1, 3, "lexer_test.hoot"},
		{"LBRACE", "{", 1, 6, "lexer_test.hoot"},
		{"NAME", "x", 1, 8, "lexer_test.hoot"},
		{"PLUS", "+", 1, 10, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 12, "lexer_test.hoot"},
		{"RBRACE", "}", 1, 14, "lexer_test.hoot"},
		{"EOF", "", 1, 15, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestRecursiveSimpleFunction(t *testing.T) {
	tokens := tokenize("x => x => x")

	expected := []Token{
		{"NAME", "x", 1, 1, "lexer_test.hoot"},
		{"ARROW", "=>", 1, 3, "lexer_test.hoot"},
		{"NAME", "x", 1, 6, "lexer_test.hoot"},
		{"ARROW", "=>", 1, 8, "lexer_test.hoot"},
		{"NAME", "x", 1, 11, "lexer_test.hoot"},
		{"EOF", "", 1, 12, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestWhiteSpace(t *testing.T) {
	tokens := tokenize("x\t=    3\ny\t \t = 4")

	expected := []Token{
		{"NAME", "x", 1, 1, "lexer_test.hoot"},
		{"ASSIGN", "=", 1, 3, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 8, "lexer_test.hoot"},
		{"NEWLINE", "\n", 1, 9, "lexer_test.hoot"},
		{"NAME", "y", 2, 1, "lexer_test.hoot"},
		{"ASSIGN", "=", 2, 6, "lexer_test.hoot"},
		{"NUMBER", "4", 2, 8, "lexer_test.hoot"},
		{"EOF", "", 2, 9, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestFunction(t *testing.T) {
	tokens := tokenize("let f = (a, b) => { let x = a ** 2\nlet y = b * 2\n return x * y }")

	expected := []Token{
		{"LET", "let", 1, 1, "lexer_test.hoot"},
		{"NAME", "f", 1, 5, "lexer_test.hoot"},
		{"ASSIGN", "=", 1, 7, "lexer_test.hoot"},
		{"LPAREN", "(", 1, 9, "lexer_test.hoot"},
		{"NAME", "a", 1, 10, "lexer_test.hoot"},
		{"COMMA", ",", 1, 11, "lexer_test.hoot"},
		{"NAME", "b", 1, 13, "lexer_test.hoot"},
		{"RPAREN", ")", 1, 14, "lexer_test.hoot"},
		{"ARROW", "=>", 1, 16, "lexer_test.hoot"},
		{"LBRACE", "{", 1, 19, "lexer_test.hoot"},
		{"LET", "let", 1, 21, "lexer_test.hoot"},
		{"NAME", "x", 1, 25, "lexer_test.hoot"},
		{"ASSIGN", "=", 1, 27, "lexer_test.hoot"},
		{"NAME", "a", 1, 29, "lexer_test.hoot"},
		{"DOUBLESTAR", "**", 1, 31, "lexer_test.hoot"},
		{"NUMBER", "2", 1, 34, "lexer_test.hoot"},
		{"NEWLINE", "\n", 1, 35, "lexer_test.hoot"},
		{"LET", "let", 2, 1, "lexer_test.hoot"},
		{"NAME", "y", 2, 5, "lexer_test.hoot"},
		{"ASSIGN", "=", 2, 7, "lexer_test.hoot"},
		{"NAME", "b", 2, 9, "lexer_test.hoot"},
		{"STAR", "*", 2, 11, "lexer_test.hoot"},
		{"NUMBER", "2", 2, 13, "lexer_test.hoot"},
		{"NEWLINE", "\n", 2, 14, "lexer_test.hoot"},
		{"RETURN", "return", 3, 2, "lexer_test.hoot"},
		{"NAME", "x", 3, 9, "lexer_test.hoot"},
		{"STAR", "*", 3, 11, "lexer_test.hoot"},
		{"NAME", "y", 3, 13, "lexer_test.hoot"},
		{"RBRACE", "}", 3, 15, "lexer_test.hoot"},
		{"EOF", "", 3, 16, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestMultipleStatements(t *testing.T) {
	tokens := tokenize("x = 3\ny = 4\nf = (a) => 'test' has a")

	expected := []ShortToken{
		{"NAME", "x"},
		{"ASSIGN", "="},
		{"NUMBER", "3"},
		{"NEWLINE", "\n"},
		{"NAME", "y"},
		{"ASSIGN", "="},
		{"NUMBER", "4"},
		{"NEWLINE", "\n"},
		{"NAME", "f"},
		{"ASSIGN", "="},
		{"LPAREN", "("},
		{"NAME", "a"},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"STRING", "'test'"},
		{"HAS", "has"},
		{"NAME", "a"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestPatternMatch(t *testing.T) {
	tokens := tokenize("let f = (\n(0) => 1,\n(x) => x * f(x - 1)\n)")

	expected := []ShortToken{
		{"LET", "let"},
		{"NAME", "f"},
		{"ASSIGN", "="},
		{"LPAREN", "("},
		{"NEWLINE", "\n"},
		{"LPAREN", "("},
		{"NUMBER", "0"},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"NUMBER", "1"},
		{"COMMA", ","},
		{"NEWLINE", "\n"},
		{"LPAREN", "("},
		{"NAME", "x"},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"NAME", "x"},
		{"STAR", "*"},
		{"NAME", "f"},
		{"LPAREN", "("},
		{"NAME", "x"},
		{"MINUS", "-"},
		{"NUMBER", "1"},
		{"RPAREN", ")"},
		{"NEWLINE", "\n"},
		{"RPAREN", ")"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestForLoop(t *testing.T) {
	tokens := tokenize("for x in 1, 2, 3 { x + 1 }")

	expected := []ShortToken{
		{"FOR", "for"},
		{"NAME", "x"},
		{"IN", "in"},
		{"NUMBER", "1"},
		{"COMMA", ","},
		{"NUMBER", "2"},
		{"COMMA", ","},
		{"NUMBER", "3"},
		{"LBRACE", "{"},
		{"NAME", "x"},
		{"PLUS", "+"},
		{"NUMBER", "1"},
		{"RBRACE", "}"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestWhileLoop(t *testing.T) {
	tokens := tokenize("while x < 3 { x + 1 }")

	expected := []ShortToken{
		{"WHILE", "while"},
		{"NAME", "x"},
		{"COMPARE", "<"},
		{"NUMBER", "3"},
		{"LBRACE", "{"},
		{"NAME", "x"},
		{"PLUS", "+"},
		{"NUMBER", "1"},
		{"RBRACE", "}"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestIfStatement(t *testing.T) {
	tokens := tokenize("if x < 3 { x + 1 } else { x - 1 }")

	expected := []ShortToken{
		{"IF", "if"},
		{"NAME", "x"},
		{"COMPARE", "<"},
		{"NUMBER", "3"},
		{"LBRACE", "{"},
		{"NAME", "x"},
		{"PLUS", "+"},
		{"NUMBER", "1"},
		{"RBRACE", "}"},
		{"ELSE", "else"},
		{"LBRACE", "{"},
		{"NAME", "x"},
		{"MINUS", "-"},
		{"NUMBER", "1"},
		{"RBRACE", "}"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestStringLiteral(t *testing.T) {
	tokens := tokenize("\"\" \"hello 'world\" 'hello \"world' \"hello\\n\\\"world\\\"\"")
	// TODO: Escape sequences
	expected := []ShortToken{
		{"STRING", "\"\""},
		{"STRING", "\"hello 'world\""},
		{"STRING", "'hello \"world'"},
		{"STRING", "\"hello\\n\\\"world\\\"\""},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestAttributeAccess(t *testing.T) {
	tokens := tokenize("x.y x::y")
	expected := []ShortToken{
		{"NAME", "x"},
		{"DOT", "."},
		{"NAME", "y"},
		{"NAME", "x"},
		{"DOUBLECOLON", "::"},
		{"NAME", "y"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestSpread(t *testing.T) {
	tokens := tokenize("...x")
	expected := []ShortToken{
		{"TRIPLEDOT", "..."},
		{"NAME", "x"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestContinue(t *testing.T) {
	tokens := tokenize("continue")
	expected := []ShortToken{
		{"CONTINUE", "continue"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestBreak(t *testing.T) {
	tokens := tokenize("break")
	expected := []ShortToken{
		{"BREAK", "break"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestImport(t *testing.T) {
	tokens := tokenize("import 'x'")
	expected := []ShortToken{
		{"IMPORT", "import"},
		{"STRING", "'x'"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestIllegal(t *testing.T) {
	tokens := tokenize("x @ y + 3")
	expected := []Token{
		{"NAME", "x", 1, 1, "lexer_test.hoot"},
		{"ILLEGAL", "@", 1, 3, "lexer_test.hoot"},
		{"NAME", "y", 1, 5, "lexer_test.hoot"},
		{"PLUS", "+", 1, 7, "lexer_test.hoot"},
		{"NUMBER", "3", 1, 9, "lexer_test.hoot"},
		{"EOF", "", 1, 10, "lexer_test.hoot"},
	}

	compareTokens(t, expected, tokens)
}

func TestPrint(t *testing.T) {
	tokens := tokenize("print x")
	expected := []ShortToken{
		{"PRINT", "print"},
		{"NAME", "x"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestNull(t *testing.T) {
	tokens := tokenize("null")
	expected := []ShortToken{
		{"NULL", "null"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestSlicing(t *testing.T) {
	tokens := tokenize("x[1:3]")
	expected := []ShortToken{
		{"NAME", "x"},
		{"LBRACKET", "["},
		{"NUMBER", "1"},
		{"COLON", ":"},
		{"NUMBER", "3"},
		{"RBRACKET", "]"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestNullCoalesce(t *testing.T) {
	tokens := tokenize("x ?? y")
	expected := []ShortToken{
		{"NAME", "x"},
		{"DOUBLEQUESTION", "??"},
		{"NAME", "y"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestNullAccess(t *testing.T) {
	tokens := tokenize("x?.y")
	expected := []ShortToken{
		{"NAME", "x"},
		{"QUESTIONDOT", "?."},
		{"NAME", "y"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestNullDeepAccess(t *testing.T) {
	tokens := tokenize("x?::y")
	expected := []ShortToken{
		{"NAME", "x"},
		{"QUESTIONDOUBLECOLON", "?::"},
		{"NAME", "y"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestNullCall(t *testing.T) {
	tokens := tokenize("x?(y)")
	expected := []ShortToken{
		{"NAME", "x"},
		{"QUESTIONLPAREN", "?("},
		{"NAME", "y"},
		{"RPAREN", ")"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestPipe(t *testing.T) {
	tokens := tokenize("() => 1 \n | () => 2")
	expected := []ShortToken{
		{"LPAREN", "("},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"NUMBER", "1"},
		{"NEWLINE", "\n"},
		{"PIPE", "|"},
		{"LPAREN", "("},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"NUMBER", "2"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}

func TestWhen(t *testing.T) {
	tokens := tokenize("x = when a < b (a, b) => 1")
	expected := []ShortToken{
		{"NAME", "x"},
		{"ASSIGN", "="},
		{"WHEN", "when"},
		{"NAME", "a"},
		{"COMPARE", "<"},
		{"NAME", "b"},
		{"LPAREN", "("},
		{"NAME", "a"},
		{"COMMA", ","},
		{"NAME", "b"},
		{"RPAREN", ")"},
		{"ARROW", "=>"},
		{"NUMBER", "1"},
		{"EOF", ""},
	}

	compareShortTokens(t, expected, tokens)
}
