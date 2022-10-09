package lib

import (
	"jack-compiler/utils"
)

var keywords = map[string]bool{
	"class": true, "constructor": true, "method": true, "function": true,
	"int": true, "boolean": true, "char": true, "void": true, "var": true,
	"static": true, "field": true, "let": true, "do": true, "if": true,
	"else": true, "while": true, "return": true, "true": true, 
	"false": true, "null": true, "this": true,
};

var symbols = map[string]bool{
	"(": true, ")": true, "{": true, "}": true, "[": true, "]": true,
	"/": true, "*": true, "+": true, "-": true, ",": true, ".": true,
	"=": true, ";": true, "&": true, "|": true, "<": true, ">": true,
	"~": true,
};


type Lexer struct {
	Source string
	Index int
}