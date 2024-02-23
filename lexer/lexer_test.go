package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	t.Run("single token test", func (t *testing.T) {
		input := `=+(){},;`

		tests := []struct {
			expectedType token.TokenType
			expectedLiteral string
		} {
			{token.ASSIGN, "="},
				{token.PLUS, "+"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.COMMA, ","},
				{token.SEMICOLON, ";"},
			}

		l := New(input)
		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - TokenType wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})

	t.Run("test snippet of monkey code", func(t *testing.T){
		input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);`

		tests := []struct{
			expectedType token.TokenType
			expectedLiteral string
		} {
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
		}

		l := New(input)
		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - TokenType wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
			}
		}

	})

	t.Run("adding new tokens", func(t *testing.T) {
		input := `!-/*5;
5 < 10 > 5;`

		tests := []struct{
			expectedType token.TokenType
			expectedLiteral string
		} {
			{token.BANG, "!"},
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.ASTERISK, "*"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.GT, ">"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
		}

		l := New(input)
		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - TokenType wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
			}
		}

	})

	t.Run("tokenize a proper statement", func(t *testing.T) {
		input := `
if (5 < 10) {
    return true;
} else {
    return false;
}`

		tests := []struct{
			expectedType token.TokenType
			expectedLiteral string
		} {
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},

		}

		l := New(input)
		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - TokenType wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})

	t.Run("2 character tokens", func(t *testing.T) {
		input := `== != >= <=`

		tests := []struct{
			expectedType token.TokenType
			expectedLiteral string
		} {
			{token.EQ, "=="},
			{token.NOT_EQ, "!="},
			{token.GE, ">="},
			{token.LE, "<="},

		}

		l := New(input)
		for i, tt := range tests {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - TokenType wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - Literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	})
}
