// Engine is a wrapper around the lexer, parser, and executor. It is used to
// load and execute Owl programs.

package exec

import (
	"os"
	"path/filepath"

	"github.com/AnthonyEdvalson/owl/lexer"
	"github.com/AnthonyEdvalson/owl/parser"
)

type OwlParams struct {
	Path    string
	Program *parser.Program
	Globals map[string]*OwlObj
}

func LoadProgramFromPath(path string) (bool, *OwlParams, []parser.ParserError) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return false, nil, nil
	}

	ret, parseErr := LoadProgram(string(bytes), path)
	if parseErr != nil {
		return false, nil, parseErr
	}

	return true, ret, nil
}

func LoadProgram(contents string, fileName string) (*OwlParams, []parser.ParserError) {
	l := lexer.NewLexer(contents)
	tok := l.Tokenize(filepath.Base(fileName))
	p := parser.NewParser(tok)
	program := p.Parse()

	if len(p.Errors) > 0 {
		return nil, p.Errors
	}

	return &OwlParams{
		Path:    fileName,
		Program: program,
		Globals: map[string]*OwlObj{},
	}, nil
}

func ExecuteProgram(params *OwlParams) (*OwlObj, *TreeExecutor) {
	t := NewTreeExecutor(params.Path)
	return t.ExecProgram(params.Program, params.Globals), t
}
