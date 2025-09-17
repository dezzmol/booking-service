package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Host string
	Port string
	Grpc *GRPCConfig
}

type GRPCConfig struct {
	Port string
}

type DbConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type Config struct {
	app                *ApplicationConfig
	Db                 *DbConfig
	NotificationClient *ApplicationConfig
	PaymentClient      *ApplicationConfig
	Consul             *Consul
}

type Consul struct {
	host string
	port string
}

func (a *Application) InitConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("[InitConfig]: fatal error while reading config file: %s", err)
	}

	host := viper.GetString("application.host")
	port := viper.GetString("application.port")
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbPassword := viper.GetString("db.password")
	dbName := viper.GetString("db.name")
	dbUser := viper.GetString("db.user")
	grpcPort := viper.GetString("application.grpc.port")

	notificationHost := viper.GetString("clients.notification_client.host")
	notificationPort := viper.GetString("clients.notification_client.port")
	notificationGRPCPort := viper.GetString("clients.notification_client.grpc.port")

	paymentHost := viper.GetString("clients.payment_client.host")
	paymentPort := viper.GetString("clients.payment_client.port")
	paymentGRPCPort := viper.GetString("clients.payment_client.grpc.port")

	consulHost := viper.GetString("consul.host")
	consulPort := viper.GetString("consul.port")

	a.config = &Config{
		app: &ApplicationConfig{
			Host: host,
			Port: port,
			Grpc: &GRPCConfig{
				Port: grpcPort,
			},
		},
		Db: &DbConfig{
			Name:     dbName,
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
		},
		NotificationClient: &ApplicationConfig{
			Host: notificationHost,
			Port: notificationPort,
			Grpc: &GRPCConfig{
				Port: notificationGRPCPort,
			},
		},
		PaymentClient: &ApplicationConfig{
			Host: paymentHost,
			Port: paymentPort,
			Grpc: &GRPCConfig{
				Port: paymentGRPCPort,
			},
		},
		Consul: &Consul{
			host: consulHost,
			port: consulPort,
		},
	}

	return nil
}
