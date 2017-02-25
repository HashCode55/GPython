package tokens


type TokenType int;

const (
	//TODO : pick from top level comment to implement more - line 165
	TokenError TokenType = iota 
	TokenNumber
	TokenSTring 
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
)
