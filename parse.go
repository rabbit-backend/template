package engine

import (
	"bytes"
	"os"
)

func Execute(name string, args any) (string, []any) {
	engine := NewEngine()
	buf := new(bytes.Buffer)

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	file.WriteTo(buf)

	tmpl, err := engine.tempalte.Parse(buf.String())
	if err != nil {
		panic(err)
	}

	out := new(bytes.Buffer)
	if err := tmpl.Execute(out, args); err != nil {
		panic(err)
	}

	return out.String(), engine.parser.args
}
