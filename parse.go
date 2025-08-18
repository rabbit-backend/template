package engine

import (
	"bytes"
	"os"
)

func (engine *Engine) Execute(name string, args any) (string, []any) {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	buf := new(bytes.Buffer)

	var tmplStr string
	var ok bool

	if engine.isCacheEnabled {
		tmplStr, ok = engine.cache[name]
		if !ok {
			file, err := os.Open(name)
			if err != nil {
				panic(err)
			}

			file.WriteTo(buf)

			tmplStr = buf.String()
			engine.cache[name] = tmplStr
		}
	} else {
		file, err := os.Open(name)
		if err != nil {
			panic(err)
		}

		file.WriteTo(buf)

		tmplStr = buf.String()
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

func (engine *Engine) SetCache(cacheStatus bool) {
	engine.isCacheEnabled = cacheStatus
}
