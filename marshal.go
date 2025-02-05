package engine

import (
	"bytes"
	"encoding/json"
)

func marshal(value interface{}) string {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(value); err != nil {
		panic(err)
	}

	return buf.String()
}
