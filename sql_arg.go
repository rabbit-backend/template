package engine

type SqlParser struct {
	args        []any
	placeHolder PlaceHolder
}

func NewSqlParser() *SqlParser {
	return &SqlParser{
		args: make([]any, 0),
	}
}

func (parser *SqlParser) Parse(arg any) string {
	parser.args = append(parser.args, arg)
	return parser.placeHolder.NextToken()
}

func (parser *SqlParser) Reset() {
	parser.placeHolder.Reset()
	parser.args = make([]any, 0)
}
