package lexical

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

const (
	Keyword    = '0'
	Identifier = '1'
	Operator   = '2'
	Separator  = '3'
)

const space = ' '
const new_line = '\n'

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

// Operator: 2 Position Rune
func writeOperator(out *os.File, b *strings.Builder, p int, c rune) {
	b.Reset()
	b.WriteRune(Operator)
	b.WriteRune(space)
	b.WriteString(strconv.Itoa(p))
	b.WriteRune(space)
	b.WriteRune(c)
	b.WriteRune(new_line)

	if _, err := out.Write([]byte(b.String())); err != nil {
		log.Fatal(err)
	}
}

// Separator: 3 Position Rune
func writeSeparator(out *os.File, b *strings.Builder, p int, c rune) {
	b.Reset()
	b.WriteRune(Separator)
	b.WriteRune(space)
	b.WriteString(strconv.Itoa(p))
	b.WriteRune(space)
	b.WriteRune(c)
	b.WriteRune(new_line)

	if _, err := out.Write([]byte(b.String())); err != nil {
		log.Fatal(err)
	}
}

// Keyword: 0 Position Type
func writeKeyword(out *os.File, b *strings.Builder, p int, t rune) {
	b.Reset()
	b.WriteRune(Keyword)
	b.WriteRune(space)
	b.WriteString(strconv.Itoa(p))
	b.WriteRune(space)
	b.WriteRune(t)
	b.WriteRune(new_line)

	if _, err := out.Write([]byte(b.String())); err != nil {
		log.Fatal(err)
	}
}

// ID: 0 Position ID
func writeIdentifier(out *os.File, b *strings.Builder, p int, st string) {
	b.Reset()
	b.WriteRune(Identifier)
	b.WriteRune(space)
	b.WriteString(strconv.Itoa(p))
	b.WriteRune(space)
	b.WriteString(st)
	b.WriteRune(new_line)

	if _, err := out.Write([]byte(b.String())); err != nil {
		log.Fatal(err)
	}
}

func resolveWord(out *os.File, b *strings.Builder, s *strings.Builder, p int) {
	if s.Len() > 0 {
		k := s.String()
		if v, ok := Keywords[k]; ok {
			writeKeyword(out, b, p, v)
		} else {
			writeIdentifier(out, b, p, k)
		}
	}
	s.Reset()
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
	n := strings.Split(path.Base(input), ".")[0] + ".bo"
	fp := path.Join(path.Dir(input), n)
	out, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			panic(err)
		}
	}()
	var word_builder strings.Builder
	var builder strings.Builder
	reader := bufio.NewReader(input_file)
	index := 0
	kw_pos := 0
	for {
		if current, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
				break
			}
		} else {
			if separator, is_separator := Separators[current]; is_separator {
				resolveWord(out, &builder, &word_builder, kw_pos)
				kw_pos = 0
				writeSeparator(out, &builder, index, separator)
			} else if operator, is_operator := Operators[current]; is_operator {
				resolveWord(out, &builder, &word_builder, kw_pos)
				kw_pos = 0
				writeOperator(out, &builder, index, operator)
			} else if unicode.IsLetter(current) {
				if word_builder.Len() == 0 {
					kw_pos = index
				}
				word_builder.WriteRune(current)
			} else if unicode.IsSpace(current) {
				resolveWord(out, &builder, &word_builder, kw_pos)
				kw_pos = 0
			} else if unicode.IsNumber(current) {
				if word_builder.Len() != 0 {
					word_builder.WriteRune(current)
				} else {
					resolveWord(out, &builder, &word_builder, kw_pos)
					kw_pos = 0
				}
			}
			index++
		}
	}
	resolveWord(out, &builder, &word_builder, index)
}
