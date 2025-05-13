package engine

import (
	"bytes"
	"os"
)

func (engine *Engine) Execute(name string, args any) (string, []any) {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	buf := new(bytes.Buffer)

	tmplStr, ok := engine.cache[name]
	if !ok {
		file, err := os.Open(name)
		if err != nil {
			panic(err)
		}

		file.WriteTo(buf)

		tmplStr = buf.String()
		engine.cache[name] = tmplStr
	}

	tmpl, err := engine.tempalte.Parse(tmplStr)
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
