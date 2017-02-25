package main 

import (
	"fmt"
	"unicode/utf8"
)
// const token_names = []string[
// 	   "NAME", 
// 	   "NUMBER", 
// 	   "STRING", 
// 	   "LPAR",
//     "RPAR",
//     "COMMA",
//     "PLUS",
//     "MINUS",
//     "STAR",
//     "SLASH",
//     "VBAR",
//     "AMPER",
//     "LESS",	
//     "GREATER",
//     "EQUAL",
//     "DOT",
//     "PERCENT",
//     "EQEQUAL",
//     "NOTEQUAL",
//     "LESSEQUAL",
//     "GREATEREQUAL",
//     "TILDE",
//     "CIRCUMFLEX",
//     "LEFTSHIFT",
//     "RIGHTSHIFT",
//     "DOUBLESTAR",
//     "PLUSEQUAL",
//     "MINEQUAL",
//     "STAREQUAL",
//     "SLASHEQUAL",
//     "PERCENTEQUAL",
//     "AMPEREQUAL",
//     "VBAREQUAL",
//     "CIRCUMFLEXEQUAL",
//     "LEFTSHIFTEQUAL",
//     "RIGHTSHIFTEQUAL",
//     "DOUBLESTAREQUAL",
//     "DOUBLESLASH",
//     "DOUBLESLASHEQUAL",
//     "OP",
//     "<ERRORTOKEN>",
//     "<N_TOKENS>",
//     "BACKTICK"
// ]

//##########################//
//    TYPE AND CONST DEFS   //
//##########################//

// no fucking idea why
type tokenType int 

// this is to avoid lots of switch statements.
// like this - https://blog.gopheracademy.com/advent-2014/parsers-lexers/
type stateFunc func(*lexer) stateFunc

// definition of a token 
type token struct {
	typ tokenType // this is a wrapper over the default int
	val string	
}

// mapping integers to the token types
type lexer struct {
	name   string 	  // for error reports 
	input  string 	  // the string being scanner 
	start  int    	  // the start state 
	pos    int    	  // current position in input 
	width  int        // width of last rune 
    tokens chan token // a channel for piping the token
}

const (
	//TODO : pick from top level comment to implement more - line 165
	tokenError tokenType = iota 
	tokenNumber
	tokenString 
	tokenSpace 
	tokenName
	tokenComma 

)


const eof = -1

//##########################//
//     Lexer Definition     //
//##########################//

func lex(name, input string) (*lexer, chan token){
	l := &lexer{
		name   : name, 
		input  : input,
		tokens : make(chan token), 
	}
	go l.run()
	return l, l.tokens
}

func (l *lexer) run() {
	// Another method is to remove this and reinitiate by 
	// simple switch case 
	// Reference - https://github.com/golang/go/blob/master/src/go/scanner/scanner.go#L598-L761
	for state := initState; state != nil; {
		state = state(l)
	} 
	close(l.tokens) // close the channel.
}

// keeps on pushing to the channel 
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start : l.pos]}
	l.start = l.pos // update the start pointer 
}

func (l *lexer) peek() rune {
	r := l.next()
	if (r == eof) {
		return eof
	}
	l.backup()
	return r
}

func (l *lexer) next() rune {
	// check if its end of file 
	if l.pos >= len(l.input) { 
		l.width = 0
		return eof 
	}	
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:]) // throws out the next rune in the input string 
	l.width = 1 // updates the current width TODO : this is fucked too
	l.pos += 1
	return r

}

func (l *lexer) backup() {
	l.pos -= l.width
}

//##########################//
//      THE REAL sHiT       //
//##########################//

// pretty printing 
func (t token) String() string {
	return fmt.Sprintf("<Type : %v Value : %v>\n", t.typ, t.val)
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' 
}

func consumeSpace(l *lexer) stateFunc {

	for isWhiteSpace(l.peek()) {
		l.next()
	}
	l.emit(tokenSpace) // put the space type in the channel
	
	return initState
}

func scanIdentifier(l *lexer) stateFunc {
	// an indentifier can contain a letter as well a digit
	for isLetter(l.peek()) || isDigit(l.peek()){
		l.next()
	}
	l.emit(tokenName)	
	return initState
} 

func scanNumber(l *lexer) stateFunc {
	// currently only integers 
	for isDigit(l.peek()){
		l.next()
	}
	l.emit(tokenNumber)	
	return initState
}

func consumeGen(l *lexer) stateFunc{
	l.next()
	switch ch := l.peek(); ch{
		case ',':			
			l.emit(tokenComma)	
	}	
	return initState
}

func initState(l *lexer) stateFunc { 	
	sf := initState
	switch ch := l.peek(); {
		case isWhiteSpace(ch):
			sf = consumeSpace(l) // consume the white space 
		case isLetter(ch):
			sf = scanIdentifier(l)	
		case '0' <= ch && ch <= '9':
			sf = scanNumber(l)		
		case ch == eof:	
			return nil
		default:
			sf = consumeGen(l)	
			
		// TODO make it recognize other things like names/numbers/letters 	
	}
	return sf // this is fucking hard coded 
}


func main() {
	_, c := lex("test", "mehul 34 , ") // currently recognizing white space 
	for i := range c {
		fmt.Printf("%v\n", i)	
	}
	
}
