package app

import (
	"log"

	"booking-service/internal/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *App) initClients() {
	a.initPaymentClient()
}

func (a *App) initPaymentClient() {
	conn, err := grpc.NewClient("dns:///"+a.config.PaymentClient.Host+":"+a.config.PaymentClient.Grpc.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	a.Clients.payment = generated.NewPaymentServiceClient(conn)
	log.Printf("Payment client initialized on %v\n",
		a.config.PaymentClient.Host+":"+a.config.PaymentClient.Grpc.Port)
}
