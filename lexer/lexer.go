package lexer

import (
	"regexp"
)

type Lexer struct {
	input    string
	position int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}

type TokenType = string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	File    string
}

type TokenMatcher struct {
	Type    TokenType
	Matcher *regexp.Regexp
}

var reMap = []TokenMatcher{
	{"COMMENT", regexp.MustCompile(`//.*`)},
	{"NEWLINE", regexp.MustCompile(`\r?\n`)},

	{"IF", regexp.MustCompile(`if`)},
	{"ELSE", regexp.MustCompile(`else`)},
	{"FOR", regexp.MustCompile(`for`)},
	{"IN", regexp.MustCompile(`in`)},
	{"HAS", regexp.MustCompile(`has`)},
	{"RETURN", regexp.MustCompile(`return`)},
	{"LET", regexp.MustCompile(`let`)},
	{"WHILE", regexp.MustCompile(`while`)},
	{"CONTINUE", regexp.MustCompile(`continue`)},
	{"BREAK", regexp.MustCompile(`break`)},
	{"IMPORT", regexp.MustCompile(`import`)},
	{"PRINT", regexp.MustCompile(`print`)},
	{"NULL", regexp.MustCompile(`null`)},
	{"WHEN", regexp.MustCompile(`when`)},

	{"ARROW", regexp.MustCompile(`=>`)},

	{"COMPARE", regexp.MustCompile(`==|!=|<=|>=|<|>`)},
	{"ASSIGN", regexp.MustCompile(`([+-/*&|]|)=`)},

	{"AND", regexp.MustCompile(`and`)},
	{"OR", regexp.MustCompile(`or`)},
	{"NOT", regexp.MustCompile(`!|not`)},

	{"LPAREN", regexp.MustCompile(`\(`)},
	{"RPAREN", regexp.MustCompile(`\)`)},
	{"LBRACE", regexp.MustCompile(`\{`)},
	{"RBRACE", regexp.MustCompile(`\}`)},
	{"LBRACKET", regexp.MustCompile(`\[`)},
	{"RBRACKET", regexp.MustCompile(`\]`)},
	{"QUESTIONLPAREN", regexp.MustCompile(`\?\(`)},

	{"COMMA", regexp.MustCompile(`,`)},

	{"QUESTIONDOT", regexp.MustCompile(`\?\.`)},
	{"QUESTIONDOUBLECOLON", regexp.MustCompile(`\?\:\:`)},
	{"INCDEC", regexp.MustCompile(`\+\+|--`)},
	{"MINUS", regexp.MustCompile(`-`)},
	{"PLUS", regexp.MustCompile(`\+`)},
	{"SLASH", regexp.MustCompile(`/`)},
	{"DOUBLESTAR", regexp.MustCompile(`\*\*`)},
	{"STAR", regexp.MustCompile(`\*`)},
	{"DOUBLEQUESTION", regexp.MustCompile(`\?\?`)},
	{"PERCENT", regexp.MustCompile(`%`)},
	{"QUESTION", regexp.MustCompile(`\?`)},
	{"DOUBLECOLON", regexp.MustCompile(`\:\:`)},
	{"COLON", regexp.MustCompile(`\:`)},
	{"PIPE", regexp.MustCompile(`\|`)},
	{"TRIPLEDOT", regexp.MustCompile(`\.\.\.`)},
	{"DOT", regexp.MustCompile(`\.`)},

	{"STRING", regexp.MustCompile(`"([^\\"\n]|\\.)*"|'([^\\'\n]|\\.)*'`)},
	{"BOOL", regexp.MustCompile(`true|false`)},
	{"NUMBER", regexp.MustCompile(`[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`)},
	{"NAME", regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)},
	{"EOF", regexp.MustCompile(`$`)},
}

var WhiteSpaceMatcher = TokenMatcher{"WS", regexp.MustCompile(`[ \t]+`)}

func (l *Lexer) NextToken(line, column int, file string) (Token, int, int) {
	column += l.skipWhitespace()

	if l.position >= len(l.input) {
		return Token{Type: "EOF", Literal: "", Line: line, Column: column, File: file}, line, column
	}

	longest := 0
	var longest_tok Token

	for _, v := range reMap {
		match := v.Matcher.FindStringIndex(l.input[l.position:])

		if match == nil || match[0] != 0 {
			continue
		}

		lit := l.input[l.position : l.position+match[1]]

		if len(lit) > longest {
			longest = len(lit)
			longest_tok = Token{Type: v.Type, Literal: lit, Line: line, Column: column, File: file}
		}
	}

	if longest == 0 {
		tok := Token{Type: "ILLEGAL", Literal: l.input[l.position : l.position+1], Line: line, Column: column, File: file}
		l.position += 1
		column += 1
		return tok, line, column
	}

	l.position += longest

	if longest_tok.Type == "NEWLINE" {
		line++
		column = 1
	} else {
		column += longest
	}

	return longest_tok, line, column
}

func (l *Lexer) Tokenize(fileName string) []Token {
	tokens := []Token{}
	line := 1
	column := 1

	for {
		var tok Token
		tok, line, column = l.NextToken(line, column, fileName)
		tokens = append(tokens, tok)

		if tok.Type == "EOF" {
			break
		}
	}

	return tokens
}

func (l *Lexer) skipWhitespace() int {
	start := l.position

	for {
		match := WhiteSpaceMatcher.Matcher.FindStringIndex(l.input[l.position:])

		if match == nil || match[0] != 0 {
			break
		}

		l.position += len(WhiteSpaceMatcher.Matcher.FindString(l.input[l.position:]))
	}

	return l.position - start
}
