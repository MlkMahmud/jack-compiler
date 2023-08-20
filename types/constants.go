package types


type SubroutineKind int

const (
	Constructor SubroutineKind = iota
	Function
	Method
)

func (s SubroutineKind) String() string {
	return []string{"constructor", "function", "method"}[s]
}

var subroutineKindMap = map[string]SubroutineKind{
	"constructor": Constructor,
	"function": Function,
	"method": Method,
}

func GetSubroutineKind(str string) (SubroutineKind, bool) {
	subroutineKind, found := subroutineKindMap[str]
	return subroutineKind, found
}

type VarKind int

const (
	Field VarKind = iota
	Static
	Var
)

func (v VarKind) String() string {
	return []string{"field", "static", "var"}[v]
}

var varKindMap = map[string]VarKind{
	"field": Field,
	"static": Static,
	"var": Var,
}

func GetVarKind(str string) (VarKind, bool) {
	kind, found := varKindMap[str]
	return kind, found 
}

type BinaryOperator int

const (
	Addition BinaryOperator = iota
	Subraction
	Multiplication
	Division
	LessThan
	GreaterThan
	Equals
)

func (op BinaryOperator) String() string {
	return []string{"+", "-", "*", "/", "<", ">", "="}[op]
}

var binaryOperatorMap = map[string]BinaryOperator{
	"+": Addition,
	"-": Subraction,
	"*": Multiplication,
	"/": Division,
	"<": LessThan,
	">": GreaterThan,
	"=": Equals,
}

func GetBinaryOperator(str string) (BinaryOperator, bool) {
	op, found := binaryOperatorMap[str]
	return op, found
}

type LogicalOperator int

const (
	And LogicalOperator = iota
	Or
)

func (op LogicalOperator) String() string {
	return []string{"&", "|"}[op]
}

var logicalOperatorMap = map[string]LogicalOperator{
	"&": And,
	"|": Or,
}

func GetLogicalOperator(str string) (LogicalOperator, bool) {
	op, found := logicalOperatorMap[str]
	return op, found
}

type UnaryOperator int

const (
	ArithmeticNegation UnaryOperator = iota
	BooleanNegation
)

func (op UnaryOperator) String() string {
	return []string{"-", "~"}[op]
}

var unaryOperatorMap = map[string]UnaryOperator{
	"-": ArithmeticNegation,
	"~": BooleanNegation,
}

func GetUnaryOperator(str string) (UnaryOperator, bool) {
	op, found := unaryOperatorMap[str]
	return op, found
}

type LiteralType int

const (
	BooleanLiteral LiteralType = iota
	IntegerLiteral
	NullLiteral
	StringLiteral
	ThisLiteral
)

var literalTypesMap = map[string]LiteralType{
	"true": BooleanLiteral,
	"false": BooleanLiteral,
	"null": NullLiteral,
	"string": StringLiteral,
	"this": ThisLiteral,
}

func GetLiteralType(str string) (LiteralType, bool) {
	literal, ok := literalTypesMap[str]
	return literal, ok 
}

var KEYWORDS = map[string]bool{
	"class": true, "constructor": true, "method": true,
	"function": true, "int": true, "boolean": true,
	"char": true, "void": true, "var": true, "static": true,
	"field": true, "let": true, "do": true, "if": true,
	"else": true, "while": true, "return": true, "true": true,
	"false": true, "null": true, "this": true,
}

var SYMBOLS = map[string]bool{
	"(": true, ")": true, "{": true, "}": true, "[": true, "]": true, "/": true,
	"-": true, "+": true, "*": true, ",": true, ".": true, "=": true, ";": true,
	"&": true, "|": true, "<": true, ">": true, "~": true,
}
