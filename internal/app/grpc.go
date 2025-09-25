package app

import (
	"fmt"
	"log"
	"net"

	"booking-service/internal/generated"

	"google.golang.org/grpc"
)

func (a *Application) initGRPC() {
	s := grpc.NewServer()
	generated.RegisterBookingServiceServer(s, a.Handlers.booking)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.config.app.Host, a.config.app.Grpc.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server started on port %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
