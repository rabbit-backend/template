package engine

import "text/template"

type Engine struct {
	tempalte *template.Template
	parser   *SqlParser
}

func NewEngine() *Engine {
	parser := NewSqlParser()
	tmpl := template.New("__rabbit__").Funcs(template.FuncMap{
		"marshal":     marshal,
		"__sql_arg__": parser.Parse,
	})

	return &Engine{
		tempalte: tmpl,
		parser:   parser,
	}
}
