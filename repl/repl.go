package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"owl/exec"
	"owl/lexer"
	"owl/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	wd, _ := os.Getwd()

	env := exec.NewTreeExecutor(wd)
	env.ExecProgram(&parser.Program{Body: []parser.Statement{}})

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		tryRun(line, env)
	}
}

func tryRun(line string, env *exec.TreeExecutor) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	l := lexer.NewLexer(line)
	toks := l.Tokenize("cmd.hoot")
	p := parser.NewParser(toks)
	program := p.Parse()

	for _, error := range p.Errors {
		fmt.Println(error)
	}

	state := env.ExecBlock(program.Body)

	if state.State == exec.RETURN {
		fmt.Printf("%s\r\n", state.Return.TrueStr())
	}
}
