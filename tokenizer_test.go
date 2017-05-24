package gopython

import (
	"testing"
)

func TestTokensLength(t *testing.T) {
	// Test the number of tokens lexed
	// should be equal to 5 in this case
	tokens := LexEngineTest("2 * 2 * 2")
	if len(tokens) != 5 {
		t.Errorf("5 tokens were expected, got %v", len(tokens))
	}
	tokens = LexEngineTest("hello = nothello")
	if len(tokens) != 3 {
		t.Errorf("3 tokens were expected, got %v", len(tokens))
	}
	tokens = LexEngineTest("hello = \"hello\"")
	if len(tokens) != 3 {
		t.Errorf("3 tokens were expected, got %v", len(tokens))
	}
}

func TestTokenizerFail(t *testing.T) {
	// Test if the lexer is failing at invalid input or not
	// '?' is invalid here
	tokens := LexEngineTest("< ? 2 * 2")
	flag := false
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type_ == 0 {
			flag = true
		}
	}
	if !flag {
		t.Errorf("Lexer not crashing at invalid input")
	}
}

func TestLGTokens(t *testing.T) {
	// Test for comparison operators
	// 6 operators
	tokens := LexEngineTest("<>>><<<=>=")
	if len(tokens) != 6 {
		t.Errorf("6 tokens were expected, got %v", len(tokens))
	}
}

func TestStringFail(t *testing.T) {
	tokens := LexEngineTest("hello = \"sajdhb")
	flag := false
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type_ == 0 {
			flag = true
		}
	}
	if !flag {
		t.Errorf("Lexer not crashing at invalid input")
	}
	tokens = LexEngineTest("hello = \"sajdhb' ")
	flag = false
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type_ == 0 {
			flag = true
		}
	}
	if !flag {
		t.Errorf("Lexer not crashing at invalid input")
	}
}
