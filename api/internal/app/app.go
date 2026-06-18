// nolint: funlen
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ultrathreads/pkg/database"
	"ultrathreads/pkg/email/smtp"

	"ultrathreads/pkg/storage"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"ultrathreads/internal/config"
	"ultrathreads/internal/delivery"
	"ultrathreads/internal/repository"
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/cache"
	"ultrathreads/pkg/hash"
	"ultrathreads/pkg/logger"
	"ultrathreads/pkg/otp"
)

type Server struct {
	httpServer *http.Server
}

func newServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// @title UltraThreads API
// @version 1.0
// @description REST API for UltraThreads App

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies
	db, err := database.NewClient(cfg.MySQL.DSN)
	if err != nil {
		logger.Error(err)

		return
	}

	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)

		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	otpGenerator := otp.NewGOTPGenerator()

	storageProvider, err := newStorageProvider(cfg)
	if err != nil {
		logger.Warnf("file storage not initialized: %s", err.Error())
	}

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		SchoolsRepo:            repos.Schools,
		StudentsRepo:           repos.Students,
		StudentLessonsRepo:     repos.StudentLessons,
		CoursesRepo:            repos.Courses,
		ModulesRepo:            repos.Modules,
		LessonContentRepo:      repos.LessonContent,
		PackagesRepo:           repos.Packages,
		OffersRepo:             repos.Offers,
		PromoCodesRepo:         repos.PromoCodes,
		OrdersRepo:             repos.Orders,
		AdminsRepo:             repos.Admins,
		UsersRepo:              repos.Users,
		FilesRepo:              repos.Files,
		SurveyResultsRepo:      repos.SurveyResults,
		Cache:                  memCache,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		EmailSender:            emailSender,
		EmailConfig:            cfg.Email,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		FondyCallbackURL:       cfg.Payment.FondyCallbackURL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		OtpGenerator:           otpGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		StorageProvider:        storageProvider,
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
	})
	handlers := handler.NewHandler(services, tokenManager)

	services.Files.InitStorageUploaderWorkers(context.Background())

	// HTTP Server
	srv := newServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	// Close MySQL connection
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Error(err)
	}
}

func newStorageProvider(cfg *config.Config) (storage.Provider, error) {
	client, err := minio.New(cfg.FileStorage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.FileStorage.AccessKey, cfg.FileStorage.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	provider := storage.NewFileStorage(client, cfg.FileStorage.Bucket, cfg.FileStorage.Endpoint)

	return provider, nil
}
