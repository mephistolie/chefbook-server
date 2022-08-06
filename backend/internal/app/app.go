package app

import (
	"chefbook-server/internal/app/dependencies/repository"
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/config"
	"chefbook-server/internal/delivery/http/router"
	"chefbook-server/internal/repository/postgres"
	"chefbook-server/internal/server"
	"chefbook-server/pkg/auth"
	"chefbook-server/pkg/hash"
	"chefbook-server/pkg/logger"
	smtp "chefbook-server/pkg/mail"
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {

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

	client, err := minio.New(cfg.S3.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Secure: true,
	})

	var firebaseApp *firebase.App = nil
	if cfg.Firebase.Enabled {
		firebaseKeyPath := fmt.Sprintf("%s/%s", configPath, cfg.Firebase.PrivateKeyFileName)
		opt := option.WithCredentialsFile(firebaseKeyPath)
		firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	repositories := repository.NewRepository(db, client, firebaseApp, cfg.Firebase.ApiKey)
	services := service.NewService(service.Dependencies{
		Repo:                  repositories,
		HashManager:           hashManager,
		TokenManager:          tokenManager,
		MailSender:            emailSender,
		MailConfig:            cfg.Mail,
		AccessTokenTTL:        cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:       cfg.Auth.JWT.RefreshTokenTTL,
		CacheTTL:              int64(cfg.CacheTTL.Seconds()),
		Environment:           cfg.Environment,
		Domain:                cfg.HTTP.Host,
		FirebaseImportEnabled: cfg.Firebase.Enabled,
	})
	handler := router.NewRouter(services, tokenManager)

	srv := server.NewServer(cfg, handler.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("server started")

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
