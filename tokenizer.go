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
)


const eof = -1

//##########################//
//      THE REAL sHiT       //
//##########################//

// pretty printing 
// func (t token) String() string {
// 	switch t.typ {
// 		case tokenEOF:
// 			return "EOF"
// 		case tokenError:
// 			return t.val 		
// 	}
// }

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}


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
	//l.width = Pos(w) // updates the current width TODO : this is fucked too
	l.pos += 1
	return r

}

func (l *lexer) backup() {
	l.pos -= l.width
}

func initState(l *lexer) stateFunc { 	
	switch ch := l.peek(); {
		case isWhiteSpace(ch):
			consumeSpace(l) // consume the white space 
		// TODO make it recognize other things like names/numbers/letters 	
	}
	return nil // this is fucking hard coded 
}

func consumeSpace(l *lexer) stateFunc {
	for isWhiteSpace(l.peek()) {
		l.next()
	}
	l.emit(tokenSpace) // put the space type in the channel
	return initState
}

func main() {
	_, c := lex("test", "    ") // currently recognizing white space 
	fmt.Printf("%v\n", -c)
}
