package gpython

import "testing"

func TestParseSuccess(t *testing.T) {
	_, err := ParseEngine("hello = 123 - 22", false)
	if err != nil {
		t.Error("Parsing failed for fine syntax.")
	}
	_, err = ParseEngine("hello = 123 - 22 + ab / 2 * 6", false)
	if err != nil {
		t.Error("Parsing failed for fine syntax.")
	}
}

func TestParserFail(t *testing.T) {
	_, err := ParseEngine("hello = 123 +* aas", false)
	if err == nil {
		t.Error("Parser not failing when bad syntax.")
	}

	_, err = ParseEngine("ll = 12 + 12 + \"asd\"", false)
	if err == nil {
		t.Error("Parser not failing when bad syntax.")
	}
}

func TestParseString(t *testing.T) {
	_, err := ParseEngine("hello = \"hashcode\"", false)
	if err != nil {
		t.Error("Couldn't parse string allocation.")
	}
}

func TestParenthesis(t *testing.T) {
	// TestParenthesisPass
	_, err := ParseEngine("hello = ((((3))))", false)
	if err != nil {
		t.Error("Parser failed for fine syntax")
	}

	_, err = ParseEngine("hello = (3) + ((4))", false)
	if err != nil {
		t.Error("Parser failed for fine syntax")
	}

	_, err = ParseEngine("hello = (3) + (4 + (5 * 6 / (7))) * (9 - (6))", false)
	if err != nil {
		t.Error("Parser failed for fine syntax")
	}


	// TestParenthesisFail
	_, err = ParseEngine("hello = ()", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = 2 + 3 + ()", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}
	
	_, err = ParseEngine("() = a + 2", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("(hello) = a + 2", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = (()(((()", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = (2 + 3 + (a + 6)", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = )", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = )(", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = (a + ) 2", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = a (+) 2", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}

	_, err = ParseEngine("hello = (a + 2)(b + 3)", false)
	if err == nil {
		t.Error("Parser not failing for bad syntax")
	}
}

func TestAST(t *testing.T) {
	astNode, _ := ParseEngine("hello = 123 - 22", false)
	tokenList := astNode.Traverse(astNode)
	// Hard matching the AST tokens
	hardTokens := []string{"=", "hello", "-", "123", "22"}
	for i, token := range tokenList {
		if token != hardTokens[i] {
			t.Error("Error in constructing the abstract syntax tree")
		}
	}

	astNode, _ = ParseEngine("hello = \"hashcode\"", false)
	tokenList = astNode.Traverse(astNode)
	// Hard matching the AST tokens
	hardTokens = []string{"=", "hello", "\"hashcode\""}
	for i, token := range tokenList {
		if token != hardTokens[i] {
			t.Error("Error in constructing the abstract syntax tree")
		}
	}
}
