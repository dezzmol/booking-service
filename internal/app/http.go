package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"booking-service/internal/generated"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *Application) initHTTP() {
	ctx := context.Background()
	mainMux := http.NewServeMux()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем gRPC endpoint для HTTP шлюза
	err := generated.RegisterBookingServiceHandlerFromEndpoint(
		ctx,
		mux,
		a.config.app.Host+":"+a.config.app.Grpc.Port,
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	mainMux.Handle("/", mux) // grpc-gateway маршруты

	swaggerJSON, err := os.ReadFile("internal/generated/api.swagger.json")
	if err != nil {
		log.Fatalf("failed to read swagger file: %v", err)
	}

	// Обработчик для самого swagger.json файла
	mainMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var swaggerSpec map[string]interface{}
		if err := json.Unmarshal(swaggerJSON, &swaggerSpec); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Добавляем базовый путь
		swaggerSpec["basePath"] = "/booking-service"

		// Сериализуем обратно в JSON
		modifiedJSON, err := json.Marshal(swaggerSpec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(modifiedJSON)
	})

	mainMux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))
	mainMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("HTTP gateway started on %s", a.config.app.Host+":"+a.config.app.Port)
	if err := http.ListenAndServe(a.config.app.Host+":"+a.config.app.Port, mainMux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
