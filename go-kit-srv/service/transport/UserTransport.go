package transport

import (
	"go-kit-srv/service"
	"go-kit-srv/service/endpoint"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if uid,ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return endpoint.UserRequest{
			Uid:uid,
		}, nil
	}
	return nil, errors.New("参数错误")
}

func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-Type", contentType)
	statusCode := http.StatusInternalServerError
	if merr,ok := err.(*service.MyError); ok {
		statusCode = merr.Code
		body = []byte(merr.Msg)
	}
	w.WriteHeader(statusCode)
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}
	w.Write(body)
}
