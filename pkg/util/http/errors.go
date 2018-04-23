package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errMsg string

func (e errMsg) Error() string {
	return string(e)
}

const (
	ErrorJSONDecoding          errMsg = "[JSON] decoding json body"
	ErrorJSONEncoding          errMsg = "[JSON] encoding json body"
	ErrorEmptyBody             errMsg = "[HTTP] empty message body"
	ErrorDecodingPathUserId    errMsg = "[URLPath] error decoding userId path"
	ErrorDecodingPathTokenId   errMsg = "[URLPath] error decoding tokenId path"
	ErrorDecodingPathTopicId   errMsg = "[URLPath] error decoding topicId path"
	ErrorDecodingPathMessageId errMsg = "[URLPath] error decoding messageId path"
	ErrorCreatingResource      errMsg = "[Resource] error creating resource"
	ErrorDeletingResource      errMsg = "[Resource] error deleting resource"
	ErrorUpdatingResource      errMsg = "[Resource] error updating resource"
	ErrorGetResource           errMsg = "[Resource] error geting resource"
	ErrorParsingQueryParams    errMsg = "[Query] error parsing query parameters"
)

var errCodeMap = map[errMsg]int{
	ErrorJSONDecoding:          http.StatusInternalServerError,
	ErrorJSONEncoding:          http.StatusInternalServerError,
	ErrorEmptyBody:             http.StatusBadRequest,
	ErrorDecodingPathTokenId:   http.StatusBadRequest,
	ErrorDecodingPathUserId:    http.StatusBadRequest,
	ErrorDecodingPathTopicId:   http.StatusBadRequest,
	ErrorDecodingPathMessageId: http.StatusBadRequest,
	ErrorCreatingResource:      http.StatusInternalServerError,
	ErrorDeletingResource:      http.StatusInternalServerError,
	ErrorUpdatingResource:      http.StatusInternalServerError,
	ErrorGetResource:           http.StatusInternalServerError,
	ErrorParsingQueryParams:    http.StatusBadRequest,
}

type errorMsg struct {
	Msg          errMsg `json:"msg"`
	OrginalError error  `json:"-"`
}

func (e errorMsg) Error() string {
	return fmt.Sprintf("%s -- %s", string(e.Msg), e.OrginalError)
}

func HandleError(w http.ResponseWriter, eCode errMsg, uerr error) {
	e := &errorMsg{Msg: eCode, OrginalError: uerr}
	var code int
	var ok bool
	code, ok = errCodeMap[e.Msg]
	fmt.Println("My Code", code)
	if !ok {
		http.Error(w, fmt.Sprintf("Unknown error: %s", string(e.Msg)), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&Response{Code: code, Data: nil, Error: eCode}); err != nil {
		http.Error(w, string(ErrorJSONEncoding), http.StatusInternalServerError)
		return
	}
}
