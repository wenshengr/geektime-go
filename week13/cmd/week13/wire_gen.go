//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"week13/internal/biz"
	"week13/internal/conf"
	"week13/internal/data"
	"week13/internal/server"
	"week13/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, error) {
	dataData, err := data.NewData(confData)
	if err != nil {
		return nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService)
	grpcServer := server.NewGRPCServer(confServer, greeterService)
	app := newApp(logger, httpServer, grpcServer)
	return app, nil
}
