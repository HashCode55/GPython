package tokens

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

type TokenType int 
const (
	//TODO : pick from top level comment to implement more - line 165
	TokenError TokenType = iota 
	TokenNumber
	TokenString 
	Token 
	TokenSpace 
	TokenName
	TokenComma 
	Token
)