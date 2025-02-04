package usecase

import (
	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/usecase/repo"
	"github.com/Akorm0181/yelp/pkg/logger"
	"github.com/Akorm0181/yelp/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo               UserRepoI
	SessionRepo            SessionRepoI
	BookmarkRepo           BookmarkRepoI
	BusinessRepo           BusinessRepoI
	BusinessCategoryRepo   BusinessCategoryRepoI
	BusinessAttachmentRepo BusinessAttachmentRepoI
	ReviewRepo             ReviewRepoI
	ReviewAttachmentRepo   ReviewAttachmentRepoI
	ReportRepo             ReportRepoI
	NotificationRepo       NotificationRepoI
	EventRepo              EventRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:               repo.NewUserRepo(pg, config, logger),
		BookmarkRepo:           repo.NewBookmarkRepo(pg, config, logger),
		SessionRepo:            repo.NewSessionRepo(pg, config, logger),
		BusinessRepo:           repo.NewBusinessRepo(pg, config, logger),
		BusinessCategoryRepo:   repo.NewBusinessCategoryRepo(pg, config, logger),
		BusinessAttachmentRepo: repo.NewBusinessAttachmentRepo(pg, config, logger),
		ReviewRepo:             repo.NewReviewRepo(pg, config, logger),
		ReviewAttachmentRepo:   repo.NewReviewAttachmentRepo(pg, config, logger),
		ReportRepo:             repo.NewReportRepo(pg, config, logger),
		NotificationRepo:       repo.NewNotificationRepo(pg, config, logger),
		EventRepo:              repo.NewEventRepo(pg, config, logger),
	}
}
