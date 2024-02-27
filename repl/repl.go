package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/evaluator"
)

const PROMPT = "PIG>> "

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		lex := lexer.New(line)

		// for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
		// 	fmt.Fprintf(out, "%+v\n", tok)
		// }

		parser := parser.New(lex)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, fmt.Sprintf("Parser errors %d:\n", len(errors)))

	for _, msg := range errors {
		io.WriteString(out, "\t" + msg + "\n")
	}
}
