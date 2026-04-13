package main

import (
	"github.com/trustfab/gapchain/backend/internal/config"
	httpserver "github.com/trustfab/gapchain/backend/internal/handler/http"
	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
	repo "github.com/trustfab/gapchain/backend/internal/repository/fabric"
	usecase "github.com/trustfab/gapchain/backend/internal/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			infra.NewGatewayRegistry,

			// Repositories
			repo.NewLohangRepo,
			repo.NewNhatkyRepo,
			repo.NewGiaodichRepo,

			// Usecases
			usecase.NewLohangUsecase,
			usecase.NewNhatkyUsecase,
			usecase.NewGiaodichUsecase,

			// HTTP Handlers
			httpserver.NewLohangHandler,
			httpserver.NewNhatkyHandler,
			httpserver.NewGiaodichHandler,
			httpserver.NewAuthHandler,
		),
		fx.Invoke(
			httpserver.SetupRouter,
		),
	).Run()
}
