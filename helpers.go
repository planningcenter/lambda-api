package api

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
)

// Abort immediately cancels an API request
func Abort() {
	runtime.Goexit()
}

// WriteJSON performs a JSON write into the writer
func WriteJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")

	var data []byte
	switch v := body.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		out, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		data = out
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))

	w.Write(data)
}
