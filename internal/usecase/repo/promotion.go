package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/Akorm0181/yelp/pkg/logger"
	"github.com/Akorm0181/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type PromotionRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewPromotionRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *PromotionRepo {
	return &PromotionRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *PromotionRepo) Create(ctx context.Context, req entity.Promotion) (entity.Promotion, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("promotions").
		Columns(`id, user_id, title, description, discount_percentage, start_date, expires_at`).
		Values(req.ID, req.UserID, req.Title, req.Description, req.DiscountPercentage, req.StartedAt, req.ExpiresAt).ToSql()
	if err != nil {
		return entity.Promotion{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Promotion{}, err
	}

	return req, nil
}

func (r *PromotionRepo) GetSingle(ctx context.Context, req entity.PromotionSingleRequest) (entity.Promotion, error) {
	response := entity.Promotion{}
	var (
		createdAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_id, title, description, discount_percentage, start_date, expires_at, created_at`).
		From("promotions")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.Title != "":
		qeuryBuilder = qeuryBuilder.Where("title = ?", req.Title)
	case req.Description != "":
		qeuryBuilder = qeuryBuilder.Where("description = ?", req.Description)
	default:
		return entity.Promotion{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Promotion{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, response.UserID, &response.Title, &response.Description, &response.DiscountPercentage, &response.StartedAt, &response.ExpiresAt, &createdAt)
	if err != nil {
		return entity.Promotion{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)

	return response, nil
}

func (r *PromotionRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.PromotionGetList, error) {
	var (
		response  = entity.PromotionGetList{}
		createdAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_id, title, description, discount_percentage, start_date, expires_at, created_at`).
		From("promotions")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Promotion
		err = rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Description, &item.DiscountPercentage, &item.StartedAt, &item.ExpiresAt, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("promotions").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *PromotionRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("promotions").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
