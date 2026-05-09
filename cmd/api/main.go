package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	viper.SetConfigFile("config/local.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	port := ":" + viper.GetString("server.port")

	slog.Info("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
