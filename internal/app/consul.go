package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func (a *Application) initConsul() error {
	cfg := api.DefaultConfig()
	cfg.Address = a.config.Consul.host + ":" + a.config.Consul.port

	client, err := api.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("consul init err: %v\n", err)
	}
	port, err := strconv.Atoi(a.config.app.Port)
	if err != nil {
		return fmt.Errorf("consul init err: %v\n", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      "booking_service",
		Name:    "booking_service",
		Port:    port,
		Address: "booking_service",
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + "booking_service" + ":" + a.config.app.Port + "/health",
			Interval: "10s",
			Timeout:  "10s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("consul init err: %v\n", err)
	}

	log.Println("consul init success.")
	return nil
}
