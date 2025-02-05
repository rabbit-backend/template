package engine

import "fmt"

type SqlParser struct {
	args  []any
	count int
}

func NewSqlParser() *SqlParser {
	return &SqlParser{
		args:  make([]any, 0),
		count: 0,
	}
}

func (parser *SqlParser) Parse(arg any) string {
	parser.args = append(parser.args, arg)
	parser.count += 1

	return fmt.Sprintf("$%d", parser.count)
}
