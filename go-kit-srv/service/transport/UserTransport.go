package transport

import (
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
