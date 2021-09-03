package lexical

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/gmvbr/box/pkg/ast"
)

var Keywords = map[string]rune{
	"package":    '0',
	"export":     '1',
	"class":      '2',
	"interface":  '3',
	"implements": '4',
	"private":    '5',
	"public":     '6',
	"int":        '7',
	"string":     '8',
	"bool":       '9',
}

var Operators = map[rune]rune{
	0x2B: '0', // +
	0x2D: '1', // -
	0x2A: '2', // *
	0x2F: '3', // /
	0x3A: '4', // :
}

var Separators = map[rune]rune{
	0x7B: '0', // {
	0x7D: '1', // }
	0x3B: '2', // ;
}

type Analyzer struct {
	Column      int
	ColumnStart int
	Line        int
	Current     rune
	Builder     *strings.Builder
}

func resolveWord(analyzer *Analyzer) {
	if analyzer.Builder.Len() > 0 {
		k := analyzer.Builder.String()
		if _, ok := Keywords[k]; ok {
			ast.ReadToken(&ast.Token{
				Type:   ast.Keyword,
				Value:  strings.ToLower(k),
				Line:   analyzer.Line,
				Column: analyzer.ColumnStart,
			})
		} else {
			ast.ReadToken(&ast.Token{
				Type:   ast.Identifier,
				Value:  k,
				Line:   analyzer.Line,
				Column: analyzer.ColumnStart,
			})
		}
	}
	analyzer.Builder.Reset()
}

func resolveSeparator(analyzer *Analyzer) bool {
	if _, is_separator := Separators[analyzer.Current]; is_separator {
		resolveWord(analyzer)
		ast.ReadToken(&ast.Token{
			Type:   ast.Separator,
			Value:  string(analyzer.Current),
			Line:   analyzer.Line,
			Column: analyzer.Column,
		})
		return true
	}
	return false
}

func resolveOperator(analyzer *Analyzer) bool {
	if _, is_operator := Operators[analyzer.Current]; is_operator {
		resolveWord(analyzer)
		ast.ReadToken(&ast.Token{
			Type:   ast.Operator,
			Value:  string(analyzer.Current),
			Line:   analyzer.Line,
			Column: analyzer.Column,
		})
		return true
	}
	return false
}

func resolveLetter(analyzer *Analyzer) bool {
	if unicode.IsLetter(analyzer.Current) {
		if analyzer.Builder.Len() == 0 {
			analyzer.ColumnStart = analyzer.Column
		}
		analyzer.Builder.WriteRune(analyzer.Current)
		return true
	}
	return false
}

func resolveSpace(analyzer *Analyzer) bool {
	if unicode.IsSpace(analyzer.Current) {
		resolveWord(analyzer)
		if analyzer.Current == '\n' {
			analyzer.Column = 0
			analyzer.Line++
		}
		return true
	}
	return false
}

func resolveNumer(analyzer *Analyzer) bool {
	if unicode.IsNumber(analyzer.Current) {
		if analyzer.Builder.Len() != 0 {
			analyzer.Builder.WriteRune(analyzer.Current)
		} else {
			resolveWord(analyzer)
		}
		return true
	}
	return false
}

func ParseFile(input string) {
	input_file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := input_file.Close(); err != nil {
			panic(err)
		}
	}()

	var sb strings.Builder
	analyzer := &Analyzer{
		Column:      0,
		ColumnStart: 0,
		Line:        1,
		Builder:     &sb,
	}
	reader := bufio.NewReader(input_file)
	for {
		if current, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
				break
			}
		} else {
			analyzer.Current = current
			analyzer.Column++

			if resolveSeparator(analyzer) {
				continue
			} else if resolveOperator(analyzer) {
				continue
			} else if resolveLetter(analyzer) {
				continue
			} else if resolveSpace(analyzer) {
				continue
			} else if resolveNumer(analyzer) {
				continue
			}
		}
	}
	resolveWord(analyzer)
}
