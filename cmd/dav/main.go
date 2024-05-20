package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	mux "github.com/gabrielmoura/davServer/internal/http"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Carregar configurações")
	config.LoadConfig()

	log.Println("Iniciar banco de dados")
	data.InitDB()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: mux.InitMux(),
	}
	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Println("Servidor WebDAV iniciado em dav://localhost:8080/dav/")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao iniciar o servidor: %v\n", err)
		}
	}()

	// Block until a signal is received
	<-stop

	fmt.Println("\nIniciando o desligamento gracioso do servidor...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v\n", err)
	}

	fmt.Println("Servidor desligado com sucesso.")
}
