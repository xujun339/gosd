package discovery

import (
	"fmt"
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)
var ConsulClient *consulapi.Client
var serviceName string
var servicePort int
var serviceId string

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		panic(fmt.Errorf("regService client build fail"))
	}
	ConsulClient = client
	serviceId = serviceName+uuid.New().String()
}

func SetServiceNameAndPort(name string, port int)  {
	serviceName = name
	servicePort = port
}

func GetServicePort() int {
	return servicePort
}

func RegService()  {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = serviceId
	reg.Name = serviceName
	reg.Address = "127.0.0.1"
	reg.Port = servicePort
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = fmt.Sprintf("http://127.0.0.1:%d/health", servicePort)
	reg.Check = &check

	err := ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(fmt.Errorf("regService fail"))
	}
}

func UnRegService() {
	ConsulClient.Agent().ServiceDeregister(serviceId)
}
