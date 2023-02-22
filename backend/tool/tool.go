package tool

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadBody(r *http.Request) (jsonMap map[string]interface{}, err error) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	jsonMap = make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)

	return
}
