package server

import (
	"api/internal/application/handlers"
	"api/internal/application/services"

	"api/internal/infrastructure/repository"
	"api/pkg/broker"
	"api/pkg/broker/producer"
	"api/pkg/storage"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
)

type Server struct {
	httpServer     *http.Server
	brokerProducer *producer.Producer
	store          storage.Storager
}

func NewServer(addr string, store storage.Storager) *Server {
	prometheus.MustRegister(requestsTotal)

	router := mux.NewRouter()

	brokerConfig := broker.LoadConfig()
	brokerProducer := producer.NewProducer(brokerConfig.Brokers)

	userRepo := repository.NewUserRepository(store)
	userService := services.NewUserService(userRepo, *brokerProducer)
	userHandler := handlers.NewUserHandler(userService)

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestsTotal.Inc()
		w.Write([]byte("Hello"))
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer:     httpServer,
		brokerProducer: brokerProducer,
		store:          store,
	}
}

func (s *Server) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API Server started on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down API Server...")

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if s.brokerProducer != nil {
		s.brokerProducer.Close()
	}

	if s.store != nil {
		s.store.Close()
	}

	log.Println("API Server exited gracefully")
	return nil
}
