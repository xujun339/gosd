package endpoint

import (
	"errors"
	"go-kit-srv/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

// 定义request response格式

type UserRequest struct {
	Uid int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}

func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("too many requests")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEndpoint(userService service.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		return UserResponse{Result: userService.GetName(r.Uid)},nil
	}
}

