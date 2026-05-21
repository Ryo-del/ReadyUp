package main

import (
	"ReadyUp/internal/auth"
	"ReadyUp/internal/config"
	"ReadyUp/internal/db"
	apphttp "ReadyUp/internal/http"
	"ReadyUp/internal/repository"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("configs/local.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var cfg config.Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)
	pool, err := db.NewPostgres(connString)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	// repositories
	userRepo := repository.NewUserRepository(pool)
	jwtManager, err := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireDuration())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ database connected")
	port := ":" + viper.GetString("server.port")
	r := apphttp.NewRouter(userRepo, jwtManager)

	slog.Info("🚀 Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
