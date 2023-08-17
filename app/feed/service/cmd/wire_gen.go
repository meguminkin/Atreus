// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"Atreus/app/feed/service/internal/biz"
	"Atreus/app/feed/service/internal/conf"
	"Atreus/app/feed/service/internal/data"
	"Atreus/app/feed/service/internal/server"
	"Atreus/app/feed/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, client *conf.Client, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	redisClient := data.NewRedisConn(confData)
	dataData, cleanup, err := data.NewData(db, redisClient, logger)
	if err != nil {
		return nil, nil, err
	}
	publishConn := server.NewPublishClient(client, logger)
	feedRepo := data.NewFeeedRepo(dataData, publishConn, logger)
	feedUsecase := biz.NewFeedUsecase(feedRepo, jwt, logger)
	feedService := service.NewFeedService(feedUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, feedService, logger)
	httpServer := server.NewHTTPServer(confServer, feedService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
