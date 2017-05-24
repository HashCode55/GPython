/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   11.05.2017
 **/

package gpython

type TokenType int

const (
	//TODO : pick from top level comment to implement more - line 165
	TokenError TokenType = iota
	TokenNumber
	TokenString
	TokenSpace
	TokenName
	TokenComma
	TokenLpar
	TokenRpar
	TokenPlus
	TokenMinus
	TokenLess
	TokenGreater
	TokenEqual
	TokenLessEqual
	TokenGreaterEqual
	TokenNotEqual
	TokenEqEqual
	TokenStar
	TokenPercent
	TokenSlash
	TokenVbar
	TokenAmper
	TokenLeftShift
	TokenRightShift
	TokenIndent
	TokenNewLine
	TokenWhile // Keywords
	TokenIf
	TokenPrint
	TokenFor
)
