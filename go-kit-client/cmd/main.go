package main

import (
	endpoint2 "go-kit-client/service/endpoint"
	. "go-kit-client/service/transport"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/url"
	"os"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd/consul"
	"strconv"
	"time"
)

func main() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client,_ := consulapi.NewClient(config)
	kitClient := consul.NewClient(client)
	logger:= log.NewLogfmtLogger(os.Stdout)
	tags := []string{"primary"}

	instance := consul.NewInstancer(kitClient, logger, "userservice", tags, true)

	factory := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
		tgt, _:=url.Parse("http://" + serviceUrl)
		return httptransport.NewClient("GET", tgt, EncodeRequest, DecodeResponse).Endpoint(), nil, nil
	}

	endpointer := sd.NewEndpointer(instance, factory, logger)
	endpoints,_ := endpointer.Endpoints()
	fmt.Println("服务有" + strconv.Itoa(len(endpoints)) + "条")

	mylb := lb.NewRoundRobin(endpointer)


	for{
		endpoint,_ := mylb.Endpoint()
		ctx:= context.Background()
		response,err := endpoint(ctx, endpoint2.UserRequest{Uid: 101})
		if err!=nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res := response.(endpoint2.UserResponse)
		fmt.Println(res.Result)
		time.Sleep(time.Second)
	}
}
