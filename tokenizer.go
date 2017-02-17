package main 

import "fmt"
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
	tokenError tokenType = iota 
	tokenNumber
	tokenString 
	tokenName
)


const eof = -1

//##########################//
//      THE REAL sHiT       //
//##########################//

// pretty printing 
func (t token) String() string {
	switch t.typ {
		case tokenEOF:
			return "EOF"
		case tokenError:
			return t.val 		
	}
	return fmt.S
}

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
	return l, l.items
}

func (l *lexer) run() {
	for state := initState; state != nil; {
		state = state(l)
	} 
	close(l.items) // close the channel.
}

// keeps on pushing to the channel 
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start : l.pos]}
	l.start = l.pos 
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) next() rune {
	// check if its end of file 
	if l.pos > len(l.input) { 
		l.width = 0
		return eof 
	}	

}
func initState(l *lexer) stateFunc{ // checks what state it is and returns the corresponding function 
	switch ch := l.peek(); ch {
		case isWhiteSpace(ch):
			// white space state 
		case isLetter(ch):
			// scan word state	

	}
}




func main() {
	fmt.Printf("%v", itemNumber)
}