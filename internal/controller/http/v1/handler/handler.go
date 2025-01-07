package handler

import (
	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/usecase"
	"github.com/Akorm0181/yelp/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
	}
}
