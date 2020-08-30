package endpoint

import (
	"go-kit-srv/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

// 定义request response格式

type UserRequest struct {
	Uid int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		return UserResponse{Result: userService.GetName(r.Uid)},nil
	}
}

