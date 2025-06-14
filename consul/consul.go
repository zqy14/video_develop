package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	consulClient *api.Client
}
type ConsulRegister struct {
	Id      string
	Name    string
	Tags    []string
	Port    uint
	Address string
}

func ConsulConfig(host string, port uint) (*Consul, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Consul{client}, nil
}

func (c *Consul) ConsulRegister(register ConsulRegister) error {
	regs := api.AgentServiceRegistration{
		ID:      register.Id,
		Name:    register.Name,
		Tags:    register.Tags,
		Port:    int(register.Port),
		Address: register.Address,
	}
	return c.consulClient.Agent().ServiceRegister(&regs)

}
