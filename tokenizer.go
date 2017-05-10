package gython

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//##########################//
//    TYPE AND CONST DEFS   //
//##########################//

// TODO: Error handling

// this is to avoid lots of switch statements.
// like this - https://blog.gopheracademy.com/advent-2014/parsers-lexers/
//type stateFunc func(*lexer) stateFunc

// Token encapsulates a token using a type variable and
// its value.
type Token struct {
	type_ TokenType // this is a wrapper over the default int
	val   string
}

// mapping integers to the token types
type lexer struct {
	input  string     // the string being scanner
	start  int        // the start state
	pos    int        // current position in input
	tokens chan Token // a channel for piping the token
}

const EOF = -1

var DONE bool = false

//##########################//
//     Lexer Definition     //
//##########################//

func lex(input string) chan Token {
	l := &lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l.tokens
}

func (l *lexer) run() {
	// Another method is to remove this and reinitiate by
	// simple switch case
	// Reference - https://github.com/golang/go/blob/master/src/go/scanner/scanner.go#L598-L761
	for DONE == false {
		initState(l)
	}
	// for state := initState; state != nil; {
	// 	state = state(l)
	// }
	close(l.tokens) // close the channel.
}

// keeps on pushing to the channel
func (l *lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos // update the start pointer
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	if r == EOF {
		return EOF
	}
	return r
}

func (l *lexer) next() rune {
	// check if its end of file
	if l.pos >= len(l.input) {
		return EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:]) // throws out the next rune in the input string
	l.pos += 1
	return r

}

func (l *lexer) backup() {
	l.pos -= 1
}

//##########################//
//      THE REAL sHiT       //
//##########################//

// pretty printing
func (t Token) String() string {
	return fmt.Sprintf("<Type : %v Value : %v>\n", t.type_, t.val)
}

func isWhiteSpace(ch rune) bool {
	return ch == ' '
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func consumeSpace(l *lexer) {

	for isWhiteSpace(l.peek()) {
		l.next()
	}
	//l.emit(TokenSpace) // put the space type in the channel
}

func scanIdentifier(l *lexer) {
	// an indentifier can contain a letter as well a digit
	// Handle the keywords here
	var ident string
	for p := l.peek(); isLetter(p) || isDigit(p); {
		ident += string(p)
		l.next()
		p = l.peek()
	}
	if strings.Compare(ident, "print") == 0 {
		l.emit(TokenPrint)
	} else if strings.Compare(ident, "while") == 0 {
		l.emit(TokenWhile)
	} else {
		l.emit(TokenName)
	}
}

func scanNumber(l *lexer) {
	// currently only integers
	for isDigit(l.peek()) {
		l.next()
	}
	l.emit(TokenNumber)
}

func consumeLEGR(l *lexer, tok rune,
	tokenLR,
	tokenLGEqual,
	tokenLRShift TokenType) {
	/*
		General function for handling bot '<' and '>'
		related operators
		Left/Right shift is overriding Less/Greater
	*/
	nt := l.peek()
	if nt == '=' {
		l.next()
		l.emit(tokenLGEqual)
	} else if nt == tok { // for handling shift operators
		l.next()
		l.emit(tokenLRShift)
	} else {
		l.emit(tokenLR)
	}
}

func consumeGen(l *lexer) {
	ch := l.peek()
	l.next()
	switch ch {
	case '=':
		l.emit(TokenEqual)
	case ',':
		l.emit(TokenComma)
	case '{':
		l.emit(TokenLpar)
	case '}':
		l.emit(TokenRpar)
	case '+':
		l.emit(TokenPlus)
	case '-':
		l.emit(TokenMinus)
	case '*':
		l.emit(TokenStar)
	case '%':
		l.emit(TokenPercent)
	case '/':
		l.emit(TokenSlash)
	case '\t':
		l.emit(TokenIndent)
	case '\n':
		l.emit(TokenNewLine)
	case '<':
		consumeLEGR(l, '<', TokenLess, TokenLessEqual, TokenLeftShift)
	case '>':
		consumeLEGR(l, '>', TokenGreater, TokenGreaterEqual, TokenRightShift)
	}
}

func initState(l *lexer) {
	switch ch := l.peek(); {
	case isWhiteSpace(ch):
		consumeSpace(l) // consume the white space
	case isLetter(ch):
		scanIdentifier(l)
	case '0' <= ch && ch <= '9':
		scanNumber(l)
	case ch == EOF:
		DONE = true
	default:
		consumeGen(l)
	}
}

func Lexer_Test(prog string) {
	token_chan := lex(prog)
	out := false
	for {
		select {
		case token, ok := <-token_chan:
			if ok {
				fmt.Println(token)
			} else {
				out = true
			}
		}
		if out {
			break
		}
	}
}
func Lexer(prog string) chan Token {
	// recieves the program
	token_chan := lex(prog) // currently recognizing white space
	return token_chan
}
