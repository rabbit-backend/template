package main

import (
	"fmt"

	engine "github.com/rabbit-backend/template"
)

func main() {
	sqlEngine := engine.NewEngineWithPlaceHolder(engine.NewSqlitePlaceholder())

	query, args := sqlEngine.Execute(
		"test/app.sql",
		map[string]map[string]string{
			"args": {"id": "123"},
			"data": {"name": "template engine", "version": "1.0.0"},
		},
	)

	fmt.Println(query, args)
}
