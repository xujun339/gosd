package main

// 用户微服务提供者
import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go-kit-srv/service"
	"go-kit-srv/service/discovery"
	"go-kit-srv/service/endpoint"
	"go-kit-srv/service/transport"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)
func main() {

	name := flag.String("name", "", "服务名")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()

	if *name == "" {
		panic(fmt.Errorf("请指定服务名"))
	}
	if *port == 0 {
		panic(fmt.Errorf("请指定端口"))
	}
	discovery.SetServiceNameAndPort(*name, *port)

	user:=service.UserService{}
	limit := rate.NewLimiter(1,  3)
	userEndpoint:= endpoint.RateLimit(limit)(endpoint.GenUserEndpoint(user))
	serverHandler := httptransport.NewServer(userEndpoint, transport.DecodeUserRequest, transport.EncodeUserResponse)
	r := mux.NewRouter()
	r.Handle("/user/{uid:\\d+}", serverHandler)
	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		w.Write([]byte("ok!"))
	})
	errChan := make(chan error)
	go func() {
		discovery.RegService()
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r)
		if err!=nil {
			log.Println(err)
			errChan <- err
		}
	}()

	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan<-fmt.Errorf("%s", <-sigChan)
	}()

	getError := <-errChan

	discovery.UnRegService()
	log.Println(getError)



}
