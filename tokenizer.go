/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   10.03.2017
 **/

package gpython

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"unicode/utf8"
)

//##########################//
//    TYPE AND CONST DEFS   //
//##########################//

// Token encapsulates a token using a type variable and
// its Value.
type Token struct {
	Type_ TokenType // his is a wrapper over the default int
	Val   string    // Val stores the Value of the token
}

// String cooks a pretty string for logging.
func (t Token) String() string {
	return fmt.Sprintf("<Type : %v Value : %v>", t.Type_, t.Val)
}

// Lexer objects stores the input as a string
// keeps two pointers to manage the tokens, "Start" and "Pos"
// Tokens is a channel in which the Lexer pushes.
type Lexer struct {
	input  string     // the string being scanner
	start  int        // the start state
	pos    int        // current position in input
	tokens chan Token // a channel for piping the token
}

const EOF = -1

//##########################//
//     Lexer Definition     //
//##########################//

// lex creates a new lex object
func lex(input string) chan Token {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.Run()
	return l.tokens
}

// Run is a wrapper over the main call
func (l *Lexer) Run() {
	// Another method is to remove this and reinitiate by
	// simple switch case
	// Reference - https://github.com/golang/go/blob/master/src/go/scanner/scanner.go#L598-L761
	defer close(l.tokens) // close the channel.
	initState(l)
}

// Emit keeps on pushing to the channel
func (l *Lexer) Emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos // update the start pointer
}

// Peek is for looking up the next rune not consuming it.
func (l *Lexer) Peek() rune {
	r := l.Next()
	// Nasty little bug!
	if r == EOF {
		return r
	}
	l.Backup()
	return r
}

// Next is for consuming the token
func (l *Lexer) Next() rune {
	// check if its end of file
	if l.pos >= len(l.input) {
		return EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:]) // throws out the next rune in the input string
	l.pos += 1
	return r

}

// Backup is to take a step back
func (l *Lexer) Backup() {
	l.pos -= 1
}

//##########################//
//      Baby Helpers        //
//##########################//

func isWhiteSpace(ch rune) bool {
	return ch == ' '
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

//##########################//
//   "Consume" Functions    //
//##########################//

// consumeSpace eats up all the white space
func consumeSpace(l *Lexer) {
	for isWhiteSpace(l.Peek()) {
		l.Next()
	}
	l.start = l.pos
}

// consumeIdentifier scans th identifiers and modifies the pointers
func consumeIdentifier(l *Lexer) {
	var ident string
	for p := l.Peek(); isLetter(p) || isDigit(p); {
		ident += string(p)
		l.Next()
		p = l.Peek()
	}
	if strings.Compare(ident, "print") == 0 {
		l.Emit(TokenPrint)
	} else if strings.Compare(ident, "while") == 0 {
		l.Emit(TokenWhile)
	} else {
		l.Emit(TokenName)
	}
}

// consumeNumber scans the numbers, currenlty only integers supported
func consumeNumber(l *Lexer) {
	// currently only integers
	for isDigit(l.Peek()) {
		l.Next()
	}
	l.Emit(TokenNumber)
}

// consumeLEGR is for consuming <, >, <=, >=
func consumeLEGR(l *Lexer, tok rune, tokenLR, tokenLGEqual, tokenLRShift TokenType) {
	/*
		General function for handling both '<' and '>'
		related operators
		Left/Right shift is overriding Less/Greater
	*/
	nt := l.Peek()
	if nt == '=' {
		l.Next()
		l.Emit(tokenLGEqual)
	} else if nt == tok { // for handling shift operators
		l.Next()
		l.Emit(tokenLRShift)
	} else {
		l.Emit(tokenLR)
	}
}

// consumeString consumes the string enclosed on '' or ""
func consumeString(l *Lexer) {
	// check if the string is starting from ' or "
	cur := l.Peek()
	var quoteType, rn rune
	if cur == '"' {
		quoteType = '"'
	} else if cur == '\'' {
		quoteType = '\''
	}
	// advance one rune as we already know what it is
	l.Next()
	for l.Peek() != quoteType {
		rn = l.Next()
		if rn == EOF {
			log.Error("Lexing Failed. Error in the string token.")
			l.Emit(TokenError)
			return // without this the function will become catastrophic
		}
	}
	// advance to include the quote in the string
	l.Next()
	l.Emit(TokenString)
}

// consumeGen is for consuming general lexemes
func consumeGen(l *Lexer) {
	ch := l.Peek()
	l.Next()
	switch ch {
	case '=':
		l.Emit(TokenEqual)
	case ',':
		l.Emit(TokenComma)
	case '{':
		l.Emit(TokenLpar)
	case '}':
		l.Emit(TokenRpar)
	case '+':
		l.Emit(TokenPlus)
	case '-':
		l.Emit(TokenMinus)
	case '*':
		l.Emit(TokenStar)
	case '%':
		l.Emit(TokenPercent)
	case '/':
		l.Emit(TokenSlash)
	case '\t':
		l.Emit(TokenIndent)
	case '\n':
		l.Emit(TokenNewLine)
	case '<':
		consumeLEGR(l, '<', TokenLess, TokenLessEqual, TokenLeftShift)
	case '>':
		consumeLEGR(l, '>', TokenGreater, TokenGreaterEqual, TokenRightShift)
	default:
		log.Errorf("Lexing Failed. Unexpected token. %v", ch)
		l.Emit(TokenError)
	}
}

//##########################//
//     The Core Engine      //
//##########################//

// initState initialises the stuff
func initState(l *Lexer) {
	var ch rune
	for ch != EOF {
		switch ch = l.Peek(); {
		case isWhiteSpace(ch):
			consumeSpace(l) // consume the white space
		case isLetter(ch):
			consumeIdentifier(l)
		case '0' <= ch && ch <= '9':
			consumeNumber(l)
		case ch == '"' || ch == '\'':
			consumeString(l)
		case ch == EOF:
			break
		default:
			consumeGen(l)
		}
	}
}

// LexEngineTest is for testing the code. Returns a list of Token (fetching from the channel)
func LexEngineTest(prog string) []Token {
	token_chan := lex(prog)
	tokenList := []Token{}
	out := false
	for {
		select {
		case token, ok := <-token_chan:
			if ok {
				tokenList = append(tokenList, token)
			} else {
				out = true
			}
		}
		if out {
			break
		}
	}
	return tokenList
}

// Lexer is the core of the tokenization. It takes the program as a string as input.
// and makes a call to "lex"(Unexported). lex fires up a goroutine pushing the tokens in a
// channel which'll be used by the parser concurrently.
// PARAMS:: prog - A string carrying the program
// RETURNS:: token_chan - Channel in which the Lexer pushes.
func LexEngine(prog string) chan Token {
	token_chan := lex(prog)
	return token_chan
}
