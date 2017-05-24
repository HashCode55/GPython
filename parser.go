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

// AST is just the top level abstraction of further nodes.
type AST struct {
}

// Node inherits from the AST and builts the AST.
// TODO: Currently token is taking extra space for internal nodes/
type Node struct {
	AST         // Annymous field showing it inherits from AST
	left  *Node // Left node
	token Token // Store the token
	right *Node // Right node
}

// traverse is a breadth first traversal over the AST for testing and debugging.
func (ast *AST) traverse(root *Node) {
	// TODO: implement this
	queue := []*Node{}
	queue = append(queue, root)
	for len(queue) != 0 {
		// Pop from queue
		node := queue[0]
		queue = queue[1:]
		log.Info(node.token)
		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}
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
			log.Infoln("Channel empty. Finished Parsing without any error.\n")
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
	log.Fatalln("Parsing Failed. Bad Syntax. %v", p.nextToken.Val)
	return false
}

// atom is terminal production, returns a Node for AST
func (p *parser) atom() *Node {
	p.expect(TokenNumber)
	return &Node{left: nil, token: p.currentToken, right: nil}
}

// termExpr is the production for handling multiplication and
// division.
func (p *parser) termExpr() *Node {
	node := p.atom()
	for p.nextToken.Type_ == TokenStar || p.nextToken.Type_ == TokenSlash {
		if p.nextToken.Type_ == TokenStar {
			// A call to expect eats up the token
			p.expect(TokenStar)
		} else if p.nextToken.Type_ == TokenSlash {
			// Eat the token
			p.expect(TokenSlash)
		} else {
			log.Fatalln("Parsing Failed. Bad Syntax. %v", p.nextToken.Val)
		}
		// make the AST node
		node = &Node{left: node, token: p.currentToken, right: p.atom()}
	}
	return node
}

// factExpr is the production for handling sum and subtraction.
func (p *parser) factExpr() *Node {
	node := p.termExpr()
	for p.nextToken.Type_ == TokenPlus || p.nextToken.Type_ == TokenMinus {
		if p.nextToken.Type_ == TokenPlus {
			// A call to expect eats up the token
			p.expect(TokenPlus)

		} else if p.nextToken.Type_ == TokenMinus {
			// Eat the token
			p.expect(TokenMinus)
		} else {
			log.Fatalln("Parsing Failed. Bad Syntax. %v", p.nextToken.Val)
		}
		// make the AST node
		node = &Node{left: node, token: p.currentToken, right: p.termExpr()}
	}
	return node
}

// start is the starting production.
func (p *parser) start() *Node {
	var node *Node
	if p.accept(TokenName) {
		node = &Node{left: nil, token: p.currentToken, right: nil}
		// expect '=' token
		p.expect(TokenEqual)
		// Built the root as '=' and continue
		node = &Node{left: node, token: p.currentToken, right: nil}
		if p.nextToken.Type_ == TokenName {
			p.expect(TokenName)
			node.right = &Node{left: nil, token: p.currentToken, right: nil}
		} else if p.nextToken.Type_ == TokenString {
			p.expect(TokenString)
			node.right = &Node{left: nil, token: p.currentToken, right: nil}
		} else if p.nextToken.Type_ == TokenNumber {
			// recursive call to expression
			node.right = p.factExpr()
		} else {
			log.Fatalln("Parsing Failed. Bad Syntax. %v", p.nextToken.Val)
		}
	}
	return node
}

// Parser initialises the parser object.
func ParseEngine(input string, lg bool) {
	tokenChan := LexEngine(input)
	p := parser{tokens: tokenChan, log_: lg}
	p.advance()
	// get the ast
	ast := p.start()
	log.Info("Traversing the AST... [debugging]")
	ast.traverse(ast)
}
