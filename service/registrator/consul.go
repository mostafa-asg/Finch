package registrator

import (
	"github.com/hashicorp/consul/api"
)

type consul struct {

}

func (c *consul) Register(name , uniqueID , server string , port int ) error {
	consulClient , err := api.NewClient( api.DefaultConfig() )
	if err != nil {
		return err
	}
	
	service := api.AgentServiceRegistration {
		Name:name,
		ID:uniqueID,
		Address:server,
		Port:port,
	}

	return consulClient.Agent().ServiceRegister( &service )
}

func (c *consul) Deregister(serviceID string) error {
	consulClient , err := api.NewClient( api.DefaultConfig() )
	if err != nil {
		return err
	}

	return consulClient.Agent().ServiceDeregister(serviceID)
}

//NewConsulServiceDiscovery registers services by HashiCorp Consul
func NewConsulServiceDiscovery() *consul {
	return &consul{}
}