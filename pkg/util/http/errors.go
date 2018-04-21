package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errMsg string

const (
	ErrorJSONDecoding       errMsg = "[JSON] decoding json body"
	ErrorJSONEncoding       errMsg = "[JSON] encoding json body"
	ErrorEmptyBody          errMsg = "[HTTP] empty message body"
	ErrorDecodingPathUserId errMsg = "[URLPath] error decoding userId path"
	ErrorCreatingResource   errMsg = "[Resource] error creating resource"
	ErrorDeletingResource   errMsg = "[Resource] error deleting resource"
	ErrorUpdatingResource   errMsg = "[Resource] error updating resource"
	ErrorGetResource        errMsg = "[Resource] error geting resource"
	ErrorParsingQueryParams errMsg = "[Query] error parsing query parameters"
)

var errCodeMap = map[errMsg]int{
	ErrorJSONDecoding:       http.StatusInternalServerError,
	ErrorJSONEncoding:       http.StatusInternalServerError,
	ErrorEmptyBody:          http.StatusBadRequest,
	ErrorCreatingResource:   http.StatusInternalServerError,
	ErrorDeletingResource:   http.StatusInternalServerError,
	ErrorUpdatingResource:   http.StatusInternalServerError,
	ErrorGetResource:        http.StatusInternalServerError,
	ErrorParsingQueryParams: http.StatusBadRequest,
}

type errorMsg struct {
	Msg          errMsg `json:"msg"`
	OrginalError error  `json:"-"`
}

func (e errorMsg) Error() string {
	return fmt.Sprintf("%s -- %s", string(e.Msg), e.OrginalError)
}

func HandleError(w http.ResponseWriter, eCode errMsg, err error) {
	e := &errorMsg{Msg: eCode, OrginalError: err}
	var code int
	var ok bool
	code, ok = errCodeMap[e.Msg]
	if !ok {
		http.Error(w, fmt.Sprintf("Unknown error: %s", string(e.Msg)), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&Response{Status: code, Data: nil, Error: err}); err != nil {
		http.Error(w, string(ErrorJSONEncoding), http.StatusInternalServerError)
		return
	}
}