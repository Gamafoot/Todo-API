package app

import (
	"context"
	"os"
	"os/signal"
	"root/internal/config"
	"root/internal/database"
	"root/internal/service"
	"root/internal/storage"
	v1 "root/internal/transport/http/v1"
	"root/pkg/jwt"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func Run() error {
	cfg := config.GetConfig()

	db, err := database.NewConnect(cfg.Database.URL)
	if err != nil {
		return pkgErrors.Errorf("failed to connect to database: %v", err)
	}

	storage := storage.NewPostgresStorage(db)

	tokenManager, err := jwt.NewManager(cfg.Jwt.SigningKey)
	if err != nil {
		return err
	}

	service := service.NewService(
		cfg,
		storage,
		tokenManager,
	)

	e := echo.New()

	handler := v1.NewHandler(cfg, service, tokenManager)
	handler.InitRoutes(e.Group(""))

	g := errgroup.Group{}

	g.Go(func() error {
		if err := e.Start(":" + cfg.Http.Port); err != nil {
			return pkgErrors.WithStack(err)
		}

		return nil
	})

	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		const timeout = 5 * time.Second

		ctx, shutdown := context.WithTimeout(context.Background(), timeout)
		defer shutdown()

		if err := e.Shutdown(ctx); err != nil {
			return pkgErrors.WithStack(err)
		}

		sqlDb, err := db.DB()
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		if err := sqlDb.Close(); err != nil {
			return pkgErrors.WithStack(err)
		}

		return nil
	})

	return g.Wait()
}
