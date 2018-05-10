/*

Filename: pretty-printer.go

Description:
 
	Takes any valid JSON as input (from standard input) 
	and outputs (to standard output) HTML that transforms the
	JSON in a neat and consistent way.

To run, type this command:

	$ go run a2.go input.json > output.html

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode"
)

const (
	// kinds of tokens
	LEFT_BRACE  = iota
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	COLON
	OBJECT_COMMA
	ARRAY_COMMA
	TRUE
	FALSE
	NULL
	STRING
	ESCAPE
	NUMBER
	LETTER
	UNKNOWN
)

var kindName = map[int]string{
	LEFT_BRACE:    "   LEFT_BRACE",
	RIGHT_BRACE:   "  RIGHT_BRACE",
	LEFT_BRACKET:  " LEFT_BRACKET",
	RIGHT_BRACKET: "RIGHT_BRACKET",
	COLON:         "        COLON",
	OBJECT_COMMA:  " OBJECT_COMMA",
	ARRAY_COMMA:   "  ARRAY_COMMA",
	TRUE:          "         TRUE",
	FALSE:         "        FALSE",
	NULL:          "         NULL",
	STRING:        "       STRING",
	ESCAPE:        "       ESCAPE",
	NUMBER:        "       NUMBER",
	LETTER:        "    LETTER(S)",
	UNKNOWN:       "      UNKNOWN",
}

var kindColor = map[int]string{
	LEFT_BRACE:    "rgb(60, 63, 163)",
	RIGHT_BRACE:   "rgb(60, 63, 163)",
	LEFT_BRACKET:  "rgb(232, 71, 162)",
	RIGHT_BRACKET: "rgb(232, 71, 162)",
	COLON:         "rgb(75, 193, 183)",
	OBJECT_COMMA:  "rgb(255, 177, 76)",
	ARRAY_COMMA:   "rgb(255, 177, 76)",
	TRUE:          "rgb(198, 76, 255)",
	FALSE:         "rgb(198, 76, 255)",
	NULL:          "rgb(198, 76, 255)",
	STRING:        "rgb(232, 6, 6)",
	ESCAPE:        "rgb(71, 155, 85)",
	NUMBER:        "rgb(76, 135, 255)",
	LETTER:        "rgb(63, 79, 74)",
	UNKNOWN:       "rgb(63, 79, 74)",
}

type Token struct {
	kind   int
	lexeme string
}

// for testing; method to display tokens
func (t Token) String() string {
	return kindName[t.kind] + " " + t.lexeme
}

func isValidNumberCode(c byte) bool {
	return unicode.IsDigit(rune(c)) || (c == '.' || c == 'e' || c == 'E' || c == '-' || c == '+')
}

// runs through string character-by-character, dividing it into valid JSON tokens,
// which are returned in the order they appear, on result
func scan(s string) (result []Token) {
	n := len(s)
	insideArray := false
	for i := 0; i < n; {
		switch {
		case s[i] == ' ':
			i++
		case s[i] == '{':
			result = append(result, Token{LEFT_BRACE, "{"})
			i++
		case s[i] == '}':
			result = append(result, Token{RIGHT_BRACE, "}"})
			i++
		case s[i] == '[':
			insideArray = !insideArray
			result = append(result, Token{LEFT_BRACKET, "["})
			i++
		case s[i] == ']':
			insideArray = !insideArray
			result = append(result, Token{RIGHT_BRACKET, "]"})
			i++
		case s[i] == ',' && !insideArray:
			result = append(result, Token{OBJECT_COMMA, ","})
			i++
		case s[i] == ',' && insideArray:
			result = append(result, Token{ARRAY_COMMA, ","})
			i++
		case s[i] == ':':
			result = append(result, Token{COLON, ":"})
			i++
		case s[i] == '"':
			// STRING
			start := i
			i++
			for i < n && s[i] != '"' {
				// if there is an escape character in the string, return string
				// that has already been processed, process and return the escape
				// character, and then continue to process the string
				if s[i] == '\\' {
					result = append(result, Token{STRING, s[start:i]})
					start = i
					i++
					switch s[i] {
					case 'u':
						i += 5
						result = append(result, Token{ESCAPE, s[start:i]})
					default:
						i++
						result = append(result, Token{ESCAPE, s[start:i]})
					}
					start = i
				} else {
					i++
				}
			}
			i++
			result = append(result, Token{STRING, s[start:i]})
		case unicode.IsLetter(rune(s[i])):
			start := i
			switch {
			case i+3 < n && (s[i] == 't' || s[i] == 'n'):
				if s[i:i+4] == "true" {
					i += 4
					result = append(result, Token{TRUE, s[start:i]})
					break
				} else if s[i:i+4] == "null" {
					i += 4
					result = append(result, Token{NULL, s[start:i]})
					break
				}
				fallthrough
			case i+4 < n && s[i] == 'f':
				if s[i:i+5] == "false" {
					i += 5
					result = append(result, Token{FALSE, s[start:i]})
					break
				}
				fallthrough
			default:
				for i < n && unicode.IsLetter(rune(s[i])) {
					i++
				}
				result = append(result, Token{LETTER, s[start:i]})
			}
		case s[i] == '-' || unicode.IsDigit(rune(s[i])):
			start := i
			for i < n && isValidNumberCode(s[i]) {
				i++
			}
			result = append(result, Token{NUMBER, s[start:i]})
		default:
			result = append(result, Token{UNKNOWN, string(s[i])})
			i++
		}

	}
	return result
}

func colorizedPrint(tokens []Token) {
	fmt.Printf("<span style=\"font-family:monospace; white-space:pre\">")
	var braceCount int
	for _, t := range tokens {
		switch t.kind {
		case LEFT_BRACE:
			fmt.Printf("\n")
			printIndent(braceCount)
			printSpanTags(kindColor[t.kind], t.lexeme)
			braceCount++
			fmt.Printf("\n")
			printIndent(braceCount)
		case RIGHT_BRACE:
			fmt.Printf("\n")
			braceCount--
			printIndent(braceCount)
			printSpanTags(kindColor[t.kind], t.lexeme)
			fmt.Printf("\n")
			printIndent(braceCount)
		case COLON:
			fmt.Printf(" ")
			printSpanTags(kindColor[t.kind], t.lexeme)
			fmt.Printf(" ")
		case OBJECT_COMMA:
			printSpanTags(kindColor[t.kind], t.lexeme)
			fmt.Printf("\n")
			printIndent(braceCount)
		case ARRAY_COMMA:
			printSpanTags(kindColor[t.kind], t.lexeme)
			fmt.Printf(" ")
		case STRING:
			for i := 0; i < len(t.lexeme); i++ {
				switch t.lexeme[i] {
				case '<':
					printSpanTags(kindColor[t.kind], "&lt;")
				case '>':
					printSpanTags(kindColor[t.kind], "&gt;")
				case '&':
					printSpanTags(kindColor[t.kind], "&amp;")
				case '"':
					printSpanTags(kindColor[t.kind], "&quot;")
				case '\'':
					printSpanTags(kindColor[t.kind], "&apos;")
				default:
					printSpanTags(kindColor[t.kind], string(t.lexeme[i]))
				}
			}
		case UNKNOWN:
			// do nothing
		default:
			printSpanTags(kindColor[t.kind], t.lexeme)
		}
	}
	fmt.Printf("</span>")
}

func printSpanTags(color string, lexeme string) {
	fmt.Printf("<span style=\"color:%s\">%s</span>", color, lexeme)
}

func printIndent(num int) {
	for i := 0; i < num; i++ {
		fmt.Printf("  ")
	}
}

func main() {
	// retrieve file name as input from command line
	fileName := os.Args[1]

	// convert file into a slice of bytes
	allBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("incorrect file name!")
	}

	// convert bytes to a string
	json := string(allBytes)

	// scan json file and convert into a list of tokens
	tokens := scan(json)
/*
	// for testing, print list of tokens
	for _, t := range tokens {
		fmt.Println(t)
	}
*/
	colorizedPrint(tokens)
}

