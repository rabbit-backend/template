package engine

import (
	"bytes"
	"os"
)

func (engine *Engine) Execute(name string, args any) (string, []any) {
	engine.lock.Lock()
	defer engine.lock.Unlock()

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

	sqlArgs := engine.parser.args
	engine.parser.Reset()

	return out.String(), sqlArgs
}
