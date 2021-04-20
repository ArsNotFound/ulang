// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"github.com/Ars2014/ulang/token"
)

const (
	NoState    = -1
	NumStates  = 93
	NumSymbols = 114
)

type Lexer struct {
	src     []byte
	pos     int
	line    int
	column  int
	Context token.Context
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:     src,
		pos:     0,
		line:    1,
		column:  1,
		Context: nil,
	}
	return lexer
}

// SourceContext is a simple instance of a token.Context which
// contains the name of the source file.
type SourceContext struct {
	Filepath string
}

func (s *SourceContext) Source() string {
	return s.Filepath
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	lexer := NewLexer(src)
	lexer.Context = &SourceContext{Filepath: fpath}
	return lexer, nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = &token.Token{}
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		tok.Pos.Context = l.Context
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn
	tok.Pos.Context = l.Context

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: ';'
1: 'f'
2: 'n'
3: 'i'
4: 'f'
5: 'e'
6: 'l'
7: 's'
8: 'e'
9: 'r'
10: 'e'
11: 't'
12: 'u'
13: 'r'
14: 'n'
15: 'f'
16: 'o'
17: 'r'
18: 'n'
19: 'u'
20: 'l'
21: 'l'
22: 't'
23: 'r'
24: 'u'
25: 'e'
26: 'f'
27: 'a'
28: 'l'
29: 's'
30: 'e'
31: '|'
32: '|'
33: '&'
34: '&'
35: '!'
36: '='
37: '|'
38: '^'
39: '&'
40: '.'
41: '.'
42: '{'
43: '}'
44: ','
45: ':'
46: '+'
47: '-'
48: '('
49: ')'
50: '!'
51: '~'
52: '['
53: ']'
54: '.'
55: '='
56: '='
57: '!'
58: '='
59: '<'
60: '<'
61: '='
62: '>'
63: '>'
64: '='
65: '~'
66: '<'
67: '<'
68: '>'
69: '>'
70: '*'
71: '/'
72: '%'
73: '/'
74: '/'
75: '\n'
76: '/'
77: '*'
78: '*'
79: '*'
80: '/'
81: '_'
82: '0'
83: '0'
84: 'x'
85: 'X'
86: 'e'
87: 'E'
88: '+'
89: '-'
90: '`'
91: '`'
92: '"'
93: '\'
94: '"'
95: '"'
96: '\'
97: 'n'
98: '\'
99: 'r'
100: '\'
101: 't'
102: ' '
103: '\n'
104: '\t'
105: '\r'
106: 'a'-'z'
107: 'A'-'Z'
108: '0'-'9'
109: '0'-'7'
110: 'a'-'f'
111: 'A'-'F'
112: '1'-'9'
113: .
*/
