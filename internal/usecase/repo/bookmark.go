package repo

import (
	"context"
	"time"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/Akorm0181/yelp/pkg/logger"
	"github.com/Akorm0181/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type BookmarkRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBookmarkRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BookmarkRepo {
	return &BookmarkRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BookmarkRepo) Create(ctx context.Context, req entity.Bookmark) (entity.Bookmark, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("bookmarks").
		Columns(`id, user_id, business_id`).
		Values(req.ID, req.UserID, req.BusinessID).ToSql()
	if err != nil {
		return entity.Bookmark{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Bookmark{}, err
	}

	return req, nil
}

func (r *BookmarkRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Bookmark, error) {
	response := entity.Bookmark{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, created_at, updated_at`).
		From("bookmarks").Where("id = ?", req.ID)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Bookmark{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserID, &response.BusinessID, &createdAt, &updatedAt)
	if err != nil {
		return entity.Bookmark{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BookmarkRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BookmarksList, error) {
	var (
		response             = entity.BookmarksList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, created_at, updated_at`).
		From("bookmarks")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Bookmark
		err = rows.Scan(&item.ID, &item.UserID, &item.BusinessID, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("bookmarks").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BookmarkRepo) Update(ctx context.Context, req entity.Bookmark) (entity.Bookmark, error) {
	mp := map[string]interface{}{
		"business_id": req.BusinessID,
		"updated_at":  "now()",
	}

	query, args, err := r.pg.Builder.Update("bookmarks").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Bookmark{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Bookmark{}, err
	}

	return req, nil
}

func (r *BookmarkRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("bookmarks").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BookmarkRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	query, args, err := r.pg.Builder.Update("bookmarks").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
