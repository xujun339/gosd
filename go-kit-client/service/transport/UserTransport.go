package transport

import (
	"go-kit-client/service/endpoint"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func EncodeRequest(_ context.Context,req *http.Request,r interface{}) error{
	userRequest := r.(endpoint.UserRequest)
	req.URL.Path += "/user/" + strconv.Itoa(userRequest.Uid)
	return nil
}

func DecodeResponse(_ context.Context,res *http.Response) (response interface{}, err error) {
	if res.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var userResponse endpoint.UserResponse
	err =json.NewDecoder(res.Body).Decode(&userResponse)
	if err!=nil {
		return nil, err
	}
	return userResponse, nil
}
