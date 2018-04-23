package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Responder interface {
	StatusCode(i int)
	StatusString(s string)
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  error       `json:"error,omitempty"`
}

func (r *Response) StatusCode(i int) {
	r.Code = i
}
func (r *Response) StatusString(str string) {
	r.Status = str
}

func NewResponse(data interface{}, err error) *Response {
	return &Response{
		Data:  data,
		Error: err,
	}
}

func GetStatusString(code int) string {
	str := http.StatusText(code)
	if str == "" {
		return "Unknown status code"
	}
	return str
}

func JSONResponse(w http.ResponseWriter, code int, v Responder) {
	w.Header().Set("Content-Type", "application/json")
	if v == nil {
		HandleError(w, ErrorEmptyBody, nil)
		return
	}
	buf := &bytes.Buffer{}

	v.StatusCode(code)
	v.StatusString(GetStatusString(code))
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		HandleError(w, ErrorJSONEncoding, err)
		return
	}

	w.WriteHeader(code)
	w.Write(buf.Bytes())
}
