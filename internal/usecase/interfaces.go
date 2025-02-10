// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/Akorm0181/yelp/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.User) (entity.User, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// BusinessRepo -.
	BusinessRepoI interface {
		Create(ctx context.Context, req entity.Business) (entity.Business, error)
		GetSingle(ctx context.Context, req entity.BusinessSingleRequest) (entity.Business, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error)
		Update(ctx context.Context, req entity.Business) (entity.Business, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// BusinessCategoryRepo -.
	BusinessCategoryRepoI interface {
		Create(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error)
		GetSingle(ctx context.Context, req entity.BusinessCategorySingleRequest) (entity.BusinessCategory, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessCategoryList, error)
		Update(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// BusinessRepo -.
	BusinessAttachmentRepoI interface {
		Create(ctx context.Context, req entity.BusinessAttachment) (entity.BusinessAttachment, error)
		MultipleUpsert(ctx context.Context, req entity.BusinessAttachmentMultipleInsertRequest) ([]entity.BusinessAttachment, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.BusinessAttachment, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessAttachmentList, error)
		Delete(ctx context.Context, req entity.Id) error
		Update(ctx context.Context, req entity.BusinessAttachment) (entity.BusinessAttachment, error)
	}

	// ReviewRepo -.
	ReviewRepoI interface {
		Create(ctx context.Context, req entity.Review) (entity.Review, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Review, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewList, error)
		Update(ctx context.Context, req entity.Review) (entity.Review, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// ReviewAttachmentRepo -.
	ReviewAttachmentRepoI interface {
		Create(ctx context.Context, req entity.ReviewAttachment) (entity.ReviewAttachment, error)
		MultipleUpsert(ctx context.Context, req entity.ReviewAttachmentMultipleInsertRequest) ([]entity.ReviewAttachment, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.ReviewAttachment, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewAttachmentList, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// ReportRepo -.
	ReportRepoI interface {
		Create(ctx context.Context, req entity.Report) (entity.Report, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Report, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReportList, error)
		Update(ctx context.Context, req entity.Report) (entity.Report, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// NotificationRepo -.
	NotificationRepoI interface {
		Create(ctx context.Context, req entity.Notification) (entity.Notification, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Notification, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.NotificationList, error)
		Update(ctx context.Context, req entity.Notification) (entity.Notification, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateStatus(ctx context.Context, req entity.Notification) (entity.Notification, error)
	}

	// EventRepo -.
	EventRepoI interface {
		Create(ctx context.Context, req entity.Event) (entity.Event, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Event, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.EventList, error)
		Update(ctx context.Context, req entity.Event) (entity.Event, error)
		Delete(ctx context.Context, req entity.Id) error
		AddParticipant(ctx context.Context, req entity.EventParticipant) (entity.EventParticipant, error)
		RemoveParticipant(ctx context.Context, req entity.EventParticipant) error
		GetParticipants(ctx context.Context, req entity.GetListFilter) (entity.EventParticipantList, error)
	}

	// BookmarkRepo -.
	BookmarkRepoI interface {
		Create(ctx context.Context, req entity.Bookmark) (entity.Bookmark, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Bookmark, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BookmarksList, error)
		Update(ctx context.Context, req entity.Bookmark) (entity.Bookmark, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// PromotionRepo -.
	PromotionRepoI interface {
		Create(ctx context.Context, req entity.Promotion) (entity.Promotion, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.PromotionGetList, error)
		GetSingle(ctx context.Context, req entity.PromotionSingleRequest) (entity.Promotion, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	TagRepoI interface {
		Create(ctx context.Context, req entity.Tag) (entity.Tag, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Tag, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.TagList, error)
		Update(ctx context.Context, req entity.Tag) (entity.Tag, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
		GetDefaultTags(ctx context.Context) ([]string, error)
	}

	// UserTagRepo -.
	UserTagRepoI interface {
		Create(ctx context.Context, req entity.UserTag) (entity.UserTag, error)
		Delete(ctx context.Context, req entity.Id) error
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserTagList, error)
	}

	// FollowerRepo -.
	FollowerRepoI interface {
		UpsertOrRemove(ctx context.Context, req entity.Follower) (entity.Follower, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
	}
)
