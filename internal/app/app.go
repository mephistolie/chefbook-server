package app

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mephistolie/chefbook-server/internal/config"
	delivery "github.com/mephistolie/chefbook-server/internal/delivery/http"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres"
	"github.com/mephistolie/chefbook-server/internal/server"
	"github.com/mephistolie/chefbook-server/internal/service"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	smtp "github.com/mephistolie/chefbook-server/pkg/mail"
	"github.com/siruspen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Errorf("failed to initialize config: %s", err.Error())
		return
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logger.Errorf("failed to initialize db: %s", err.Error())
	}

	memCache := cache.NewMemoryCache()
	hashManager := hash.NewBcryptManager(cfg.Auth.SaltCost)

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)
		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	repositories := repository.NewRepositories(db)
	services := service.NewServices(service.Dependencies{
		Repos:                  repositories,
		Cache:                  memCache,
		HashManager:            hashManager,
		TokenManager:           tokenManager,
		MailSender:            emailSender,
		MailConfig:            cfg.Mail,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
	})
	handler := delivery.NewHandler(services, tokenManager)

	srv := server.NewServer(cfg, handler.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
