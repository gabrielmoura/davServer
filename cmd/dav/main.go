package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/backup"
	"github.com/gabrielmoura/davServer/internal/data"
	mux "github.com/gabrielmoura/davServer/internal/http"
	"github.com/gabrielmoura/davServer/internal/i2p"
	"github.com/gabrielmoura/davServer/internal/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.InitLogger()

	log.Logger.Info("Carregar configurações")
	err := config.LoadConfig()
	if err != nil {
		return
	}

	log.Logger.Info("Iniciar banco de dados")
	err = data.InitDB()
	if err != nil {
		return
	}
	err = backup.HandleBackup()
	if err != nil {
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: mux.InitMux(),
	}
	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Logger.Info("Iniciando o servidor...")
		if config.Conf.I2PCfg.Enabled {
			ls, err := i2p.InitI2P()
			if err != nil {
				panic(fmt.Sprintf("Erro ao iniciar o servidor I2P: %v\n", err))
			}
			log.Logger.Info(fmt.Sprintf("Servidor WebDAV iniciado em dav://%s/dav/", ls.Addr()))
			if err := server.Serve(ls); err != nil {
				panic(fmt.Sprintf("Erro ao iniciar o servidor: %v\n", err))
			}
		} else {
			log.Logger.Info(fmt.Sprintf("Servidor WebDAV iniciado em dav://localhost:%d/dav/", config.Conf.Port))
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(fmt.Sprintf("Erro ao iniciar o servidor: %v\n", err))
			}
		}
	}()

	// Block until a signal is received
	<-stop

	log.Logger.Info("Desligando o servidor...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Erro ao desligar o servidor: %v\n", err))
	}

	log.Logger.Info("Servidor desligado com sucesso.")
}
