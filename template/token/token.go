// Package token defines constants representing the lexical tokens.
package token

type Token int

const (
	ILLEGAL Token = iota
	EOF

	STRING // "text"
	INT    // 123
	BOOL   // true
	IDENT  // vars

	ADD  // +
	CALL // }}:\n

	LPAREN    // (
	RPAREN    // )
	LBRACK    // [
	RBRACK    // ]
	LDBRACE   // {{
	RDBRACE   // }}
	COMMA     // ,
	PERIOD    // .
	LARROW    // <-
	LINEBREAK // end of a larrow expression argument
)

// String returns t as string.
func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "illegal"
	case EOF:
		return "EOF"
	case STRING:
		return "string"
	case INT:
		return "int"
	case BOOL:
		return "bool"
	case IDENT:
		return "ident"
	case ADD:
		return "add"
	case LPAREN:
		return "lparen"
	case RPAREN:
		return "rparen"
	case LBRACK:
		return "lbrack"
	case RBRACK:
		return "rbrack"
	case LDBRACE:
		return "ldbrace"
	case RDBRACE:
		return "rdbrace"
	case COMMA:
		return "comma"
	case PERIOD:
		return "period"
	case LARROW:
		return "<-"
	case LINEBREAK:
		return "line break"
	default:
		return "unknown"
	}
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence.
const (
	LowestPrec  = 0 // non-operators
	HighestPrec = 3
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
func (t Token) Precedence() int {
	switch t {
	case ADD, LARROW, LDBRACE, STRING:
		return 1
	default:
		return LowestPrec
	}
}
