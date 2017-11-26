/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   09.05.2017

 * Check grammar.txt for the specification.
 * Each rule, R, defined in the grammar, becomes a method with the same name.
 * Alternatives (a1 | a2 | aN) become an if-elif-else statement
 * An optional grouping (…)* becomes a while statement that can loop over zero or more times
 * Each token reference T becomes a call to the method Accept: Accept(T).

 * If any changes in the grammar have to be made, first change it in grammar.txt and then edit this
 * file.
**/

package gpython

// Handle identifier startng with a number
import (
	"strconv"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// Parser does what its name suggests. It eats up tokens and builds an AST.
type Parser struct {
	Tokens       chan Token // Tokens is the channel that is passed by the lexer
	CurrentToken Token      // CurrentToken stores the the current token, for the building of AST
	NextToken    Token      // NextToken stores the next token for lookup (LL(1))
	log_         bool       // log_ is just for testing
}

// AST is just the top level abstraction of AST nodes.
type AST struct {
}

// Node inherits from the AST, why this? See next line
// TODO: Currently token is taking extra space for internal nodes/
// A simple solution is to create multiple types of nodes
type Node struct {
	AST         // Annymous field
	left  *Node // Left node
	token Token // Store the token
	right *Node // Right node
}

// Traverse is a breadth first traversal over the AST for testing and debugging.
func (ast *AST) Traverse(root *Node) []string {
	tokenList := []string{}
	queue := []*Node{}	
	queue = append(queue, root)
	for len(queue) != 0 {
		// Pop from queue
		node := queue[0]		
		queue = queue[1:]
		log.Info(node.token)
		tokenList = append(tokenList, node.token.Val)
		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}
	return tokenList
}


func (ast *AST) EvaluateTree(root *Node) float64 {
	if root == nil {
		return 0;
	}

	if root.left == nil && root.right == nil {
		value, _ := strconv.ParseFloat(root.token.Val, 64)
		return value
	}

	leftValue := ast.EvaluateTree(root.left)
	rightValue := ast.EvaluateTree(root.right)

	if root.token.Type_ == TokenPlus {
		return leftValue + rightValue
	}
	if root.token.Type_ == TokenMinus {
		return leftValue - rightValue
	}
	if root.token.Type_ == TokenStar {
		return leftValue * rightValue
	}
	if root.token.Type_ == TokenSlash {
		return leftValue / rightValue
	}
	return -1
}

// Advance modifies two pointers of the Parser construct.
// It makes the current token as next token and get the next token
// from the channel.
// TODO: We might not need current token variable at all.
func (p *Parser) Advance() {
	select {
	case token, ok := <-p.Tokens:
		if ok {
			if p.log_ {
				log.Infoln(token)
			}
			p.CurrentToken = p.NextToken
			p.NextToken = token
		} else {
			p.CurrentToken = p.NextToken
			// Replace the next token with an empty token
			p.NextToken = Token{}
			log.Infoln("Channel empty. Finished Parsing.\n")
		}
	}
}

// Accept is used to successfully consumes the token.
// Analogous to "eat" function.
func (p *Parser) Accept(t TokenType) bool {
	if p.NextToken.Type_ == t {
		p.Advance()
		return true
	}
	return false
}

// Expect expects the next token to be of type t
// if found true it Advances else fails.
func (p *Parser) Expect(t TokenType) error {
	if p.Accept(t) {
		return nil
	}
	// log.Fatalln("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
	err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
	return err
}

// atom is terminal production, returns a Node for AST
func (p *Parser) atom() (*Node, error) {
	// kingsguard 
	if p.NextToken.Type_ == TokenRpar {		
		err := fmt.Errorf("Parsing Failed. Right parenthesis existing without left or empty parenthesis used with operator. %v", p.NextToken.Val)
		return nil, err
	}

	// Handling parenthesis 
	if p.NextToken.Type_ == TokenLpar {		
		err := p.Expect(TokenLpar)
		node, err := p.factExpr()
		if err != nil {							
			return nil, err
		}
		err = p.Expect(TokenRpar)
		if err != nil {
			err := fmt.Errorf("Parsing Failed. Right parenthesis missing. %v", p.NextToken.Val)
			return nil, err
		}
		return node, nil 
	}    
	
	// TokenName incorporated - Puneet 
	if (p.NextToken.Type_ == TokenNumber) {
		p.Expect(TokenNumber)
	} else if (p.NextToken.Type_ == TokenName) {
		p.Expect(TokenName)
	} else {
		err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
		return nil, err
	}	
	return &Node{left: nil, token: p.CurrentToken, right: nil}, nil
}

// termExpr is the production for handling multiplication and
// division.
func (p *Parser) termExpr() (*Node, error) {
	node, err := p.atom()
	if err != nil {
		return nil, err
	}
	for p.NextToken.Type_ == TokenStar || p.NextToken.Type_ == TokenSlash {
		if p.NextToken.Type_ == TokenStar {
			// A call to Expect eats up the token
			// we do not have to check for the error here
			p.Expect(TokenStar)
		} else if p.NextToken.Type_ == TokenSlash {
			// Eat the token
			p.Expect(TokenSlash)
		} else {
			err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
			return nil, err
		}
		// make the AST node
		curtok := p.CurrentToken
		rightNode, err := p.atom()
		
		if err != nil {
			return nil, err
		}
		node = &Node{left: node, token: curtok, right: rightNode}
	}
	return node, nil
}

// factExpr is the production for handling sum and subtraction.
func (p *Parser) factExpr() (*Node, error) {	
	node, err := p.termExpr()
	if err != nil {
		return nil, err
	}
	for p.NextToken.Type_ == TokenPlus || p.NextToken.Type_ == TokenMinus {
		if p.NextToken.Type_ == TokenPlus {
			// A call to Expect eats up the token
			p.Expect(TokenPlus)

		} else if p.NextToken.Type_ == TokenMinus {
			// Eat the token
			p.Expect(TokenMinus)
		} else {
			err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
			return nil, err
		}
		// make the AST node
		curtok := p.CurrentToken
		rightNode, err := p.termExpr()
		if err != nil {
			return nil, err
		}
		node = &Node{left: node, token: curtok, right: rightNode}
	}
	return node, nil
}

// start is the starting production.
// This is uglyyyyy.
func (p *Parser) start() (*Node, error) {
	var node *Node 
	if p.Accept(TokenName) {
		node = &Node{left: nil, token: p.CurrentToken, right: nil}
		// Expect '=' token
		err := p.Expect(TokenEqual)
		if err != nil {
			return nil, err
		}
		// Build the root as '=' and continue
		node = &Node{left: node, token: p.CurrentToken, right: nil}
		
		// Else if Condition modified - Puneet 
		if p.NextToken.Type_ == TokenString {
			p.Expect(TokenString)
			node.right = &Node{left: nil, token: p.CurrentToken, right: nil}
		} else if p.NextToken.Type_ == TokenNumber || p.NextToken.Type_ == TokenName ||
			 p.NextToken.Type_ == TokenLpar || p.NextToken.Type_ == TokenRpar {
			// recursive call to expression
			node.right, err = p.factExpr()
			if err != nil {
				return nil, err
			}
		} else {
			err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
			return nil, err
		}
	} else {
		err := fmt.Errorf("Parsing Failed. Bad Syntax. %v", p.NextToken.Val)
		return nil, err
	}
	return node, nil
}

// ParseEngine is the driver od the parser, it makes a call to the LexEngine,
// gets the token channel and builds the AST.
// PARAMS:: input program and logging flag
// TODO: Top level lexer call
func ParseEngine(input string, lg bool) (*Node, error) {
	tokenChan := LexEngine(input)
	p := Parser{Tokens: tokenChan, log_: lg}
	p.Advance()
	// get the ast
	ast, err := p.start()
	if err != nil {
		log.Info(err)
		return nil, err
	}
	if lg {
		log.Info("Traversing the AST... [debugging]")
		if ast == nil {
			return nil, nil
		}
		ast.Traverse(ast)
	}

	fmt.Println(ast.EvaluateTree(ast.right))

	return ast, nil
}
