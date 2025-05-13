package engine

import (
	"sync"
	"text/template"
)

type Engine struct {
	tempalte *template.Template
	parser   *SqlParser
	lock     sync.RWMutex
	cache    map[string]string
}

func createEngine() *Engine {
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

func NewEngine() *Engine {
	engine := createEngine()
	engine.parser.placeHolder = &PostgresPlaceholder{
		index: 0,
	}

	return engine
}

func NewEngineWithPlaceHolder(placeHolder PlaceHolder) *Engine {
	engine := createEngine()
	engine.parser.placeHolder = placeHolder

	return engine
}
