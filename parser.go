/*
Author - hashcode55 (Mehul Ahuja)

Check grammar.txt for the specification.
Each rule, R, defined in the grammar, becomes a method with the same name.
Alternatives (a1 | a2 | aN) become an if-elif-else statement
An optional grouping (…)* becomes a while statement that can loop over zero or more times
Each token reference T becomes a call to the method accept: accept(T).

If any changes in the grammar have to be made, first change it in grammar.txt and then edit this
file.
*/
package gython

// Start with simple statements
// Handle identifier startng with a number
import (
	"fmt"
)

type parser struct {
	tokens       chan Token
	currentToken Token
	nextToken    Token
}

// advance modifies two pointers of the parser construct.
// It makes the current token as next token and get the next token
// from the channel.
// TODO: We might not need current token variable at all.
func (p *parser) advance() {
	select {
	case token, ok := <-p.tokens:
		if ok {
			fmt.Println(token)
			p.currentToken = p.nextToken
			p.nextToken = token
		} else {
			fmt.Println("Token chan over.")
		}
	}
}

// accept is used to successfully consumes the token.
// Analogous to "eat" function.
func (p *parser) accept(t TokenType) bool {
	if p.nextToken.type_ == t {
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
	fmt.Println("Parsing Failed. Got unexpected token.")
	return false
}

// atom is terminal production.
func (p *parser) atom() {
	p.expect(TokenNumber)
}

// term_expr is the production for handling multiplication and
// division.
func (p *parser) term_expr() {
	p.atom()
	for p.nextToken.type_ == TokenStar || p.nextToken.type_ == TokenSlash {
		if p.nextToken.type_ == TokenStar {
			// A call to expect eats up the token
			p.expect(TokenStar)

		} else if p.nextToken.type_ == TokenSlash {
			// Eat the token
			p.expect(TokenSlash)
		}
		p.atom()
	}

}

// fact_expr is the production for handling sum and subtraction.
func (p *parser) fact_expr() {
	p.term_expr()
	for p.nextToken.type_ == TokenPlus || p.nextToken.type_ == TokenMinus {
		if p.nextToken.type_ == TokenPlus {
			// A call to expect eats up the token
			p.expect(TokenPlus)

		} else if p.nextToken.type_ == TokenMinus {
			// Eat the token
			p.expect(TokenMinus)
		}
		p.term_expr()
	}

}

// start is the starting production.
func (p *parser) start() {
	if p.accept(TokenName) {
		p.expect(TokenEqual)
		if p.nextToken.type_ == TokenName {
			p.expect(TokenName)
		} else if p.nextToken.type_ == TokenString {
			p.expect(TokenString)
		}
		// recursive call to expression
		p.fact_expr()
	}
}

// Parser initialises the parser object.
func Parser(input string) {
	/*
		This is the exposed function.
	*/

	// Testing if the package works.
	token_chan := Lexer(input)
	p := parser{tokens: token_chan}
	p.advance()
	p.start()
}
