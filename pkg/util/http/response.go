package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  error       `json:"error,omitempty"`
}

func NewResponse(data interface{}, err error) *Response {
	return &Response{
		Data:  data,
		Error: err,
	}
}

func JSONResponse(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if v == nil {
		HandleError(w, ErrorEmptyBody, nil)
		return
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		HandleError(w, ErrorJSONEncoding, err)
		return
	}

	w.WriteHeader(code)
	w.Write(buf.Bytes())
}
