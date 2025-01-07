package usecase

import (
	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/usecase/repo"
	"github.com/Akorm0181/yelp/pkg/logger"
	"github.com/Akorm0181/yelp/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo    UserRepoI
	SessionRepo SessionRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:    repo.NewUserRepo(pg, config, logger),
		SessionRepo: repo.NewSessionRepo(pg, config, logger),
	}
}
