/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   09.05.2017

 * Check grammar.txt for the specification.
 * Each rule, R, defined in the grammar, becomes a method with the same name.
 * Alternatives (a1 | a2 | aN) become an if-elif-else statement
 * An optional grouping (…)* becomes a while statement that can loop over zero or more times
 * Each token reference T becomes a call to the method accept: accept(T).

 * If any changes in the grammar have to be made, first change it in grammar.txt and then edit this
 * file.
**/

package gython

// Start with simple statements
// Handle identifier startng with a number
import (
	log "github.com/Sirupsen/logrus"
)

type parser struct {
	tokens       chan Token
	currentToken Token
	nextToken    Token
	log_         bool
}

// advance modifies two pointers of the parser construct.
// It makes the current token as next token and get the next token
// from the channel.
// TODO: We might not need current token variable at all.
func (p *parser) advance() {
	select {
	case token, ok := <-p.tokens:
		if ok {
			if p.log_ {
				log.Infoln(token)
			}
			p.currentToken = p.nextToken
			p.nextToken = token
		} else {
			log.Infoln("Channel empty. Finished Parsing without any error.")
		}
	}
}

// accept is used to successfully consumes the token.
// Analogous to "eat" function.
func (p *parser) accept(t TokenType) bool {
	if p.nextToken.Type_ == t {
		p.advance()
		return true
	}
	return false
}

// expect expects the next token to be of type t
// if found true it advances else fails.
func (p *parser) expect(t TokenType) bool {
	if p.accept(t) {
		return true
	}
	log.Fatalln("Parsing Failed. Got unexpected token. %v", p.nextToken.Val)
	return false
}

// atom is terminal production.
func (p *parser) atom() {
	p.expect(TokenNumber)

}

// termExpr is the production for handling multiplication and
// division.
func (p *parser) termExpr() {
	p.atom()
	for p.nextToken.Type_ == TokenStar || p.nextToken.Type_ == TokenSlash {
		if p.nextToken.Type_ == TokenStar {
			// A call to expect eats up the token
			p.expect(TokenStar)
		} else if p.nextToken.Type_ == TokenSlash {
			// Eat the token
			p.expect(TokenSlash)
		}
		p.atom()
	}

}

// factExpr is the production for handling sum and subtraction.
func (p *parser) factExpr() {
	p.termExpr()
	for p.nextToken.Type_ == TokenPlus || p.nextToken.Type_ == TokenMinus {
		if p.nextToken.Type_ == TokenPlus {
			// A call to expect eats up the token
			p.expect(TokenPlus)

		} else if p.nextToken.Type_ == TokenMinus {
			// Eat the token
			p.expect(TokenMinus)
		} else {
			log.Fatalln("Parsing Failed. Got unexpected token. %v", p.nextToken.Val)
		}
		p.termExpr()
	}

}

// start is the starting production.
func (p *parser) start() {
	if p.accept(TokenName) {
		p.expect(TokenEqual)
		if p.nextToken.Type_ == TokenName {
			p.expect(TokenName)
		} else if p.nextToken.Type_ == TokenString {
			p.expect(TokenString)
		} else if p.nextToken.Type_ == TokenNumber {
			// recursive call to expression
			p.factExpr()
		} else {
			log.Fatalln("Parsing Failed. Got unexpected token. %v", p.nextToken.Val)
		}
	}
}

// Parser initialises the parser object.
func ParseEngine(input string, lg bool) {
	tokenChan := LexEngine(input)
	p := parser{tokens: tokenChan, log_: lg}
	p.advance()
	p.start()
}
