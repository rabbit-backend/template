package engine

import (
	"sync"
	"text/template"
)

type Engine struct {
	tempalte       *template.Template
	parser         *SqlParser
	lock           sync.RWMutex
	cache          map[string]string
	isCacheEnabled bool
}

func createEngine() *Engine {
	parser := NewSqlParser()
	tmpl := template.New("__rabbit__").Funcs(template.FuncMap{
		"marshal":                 marshal,
		"__sql_arg__":             parser.Parse,
		"__default_placeholder__": DefaultPlaceHolder,
	})

	return &Engine{
		tempalte:       tmpl,
		parser:         parser,
		cache:          make(map[string]string, 0),
		isCacheEnabled: true,
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
