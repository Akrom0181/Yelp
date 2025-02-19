package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/Akorm0181/yelp/pkg/logger"
	"github.com/Akorm0181/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type BusinessRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBusinessRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BusinessRepo {
	return &BusinessRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BusinessRepo) Create(ctx context.Context, req entity.Business) (entity.Business, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("businesses").
		Columns(`id, name, description, category_id, address, owner_id, latitude, longitude, contact_info, hours_of_operation`).
		Values(req.ID, req.Name, req.Description, req.CategoryID, req.Address, req.OwnerID, req.Latitude, req.Longitude, req.ContactInfo, req.HoursOfOperation).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) GetSingle(ctx context.Context, req entity.BusinessSingleRequest) (entity.Business, error) {
	response := entity.Business{}
	var (
		createdAt, updatedAt                       time.Time
		description, contactInfo, hoursOfOperation sql.NullString
		latitude, longitude                        sql.NullFloat64
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, description, category_id, address, latitude, longitude, contact_info, hours_of_operation, owner_id, created_at, updated_at`).
		From("businesses")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.OwnerID != "":
		qeuryBuilder = qeuryBuilder.Where("owner_id = ?", req.OwnerID)
	case req.CategoryID != "":
		qeuryBuilder = qeuryBuilder.Where("category_id = ?", req.CategoryID)
	default:
		return entity.Business{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.Name, &description, &response.CategoryID, &response.Address,
			&latitude, &longitude, &contactInfo, &hoursOfOperation, &response.OwnerID, &createdAt, &updatedAt)
	if err != nil {
		return entity.Business{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	if latitude.Valid {
		response.Latitude = latitude.Float64
	}
	if longitude.Valid {
		response.Longitude = longitude.Float64
	}
	if contactInfo.Valid {
		var contactInfoStruct entity.ContactInfo
		err := json.Unmarshal([]byte(contactInfo.String), &contactInfoStruct)
		if err != nil {
			return response, err
		}
		response.ContactInfo = contactInfoStruct
	}
	if hoursOfOperation.Valid {
		var hoursOfOperationStruct entity.HoursOfOperation
		err := json.Unmarshal([]byte(hoursOfOperation.String), &hoursOfOperationStruct)
		if err != nil {
			return response, err
		}
		response.HoursOfOperation = hoursOfOperationStruct
	}
	if description.Valid {
		response.Description = description.String
	}

	return response, nil
}

func (r *BusinessRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error) {
	var (
		response                                   = entity.BusinessList{}
		createdAt, updatedAt                       time.Time
		description, contactInfo, hoursOfOperation sql.NullString
		latitude, longitude                        sql.NullFloat64
	)

	// Fully qualify column names to avoid ambiguity
	queryBuilder := r.pg.Builder.
		Select(`
			b.id AS business_id, b.name, b.description, b.category_id, b.address, 
			b.latitude, b.longitude, b.contact_info, b.hours_of_operation, 
			b.owner_id, b.created_at, b.updated_at, 
			COALESCE(JSON_AGG(ba) FILTER (WHERE ba.id IS NOT NULL), '[]') AS attachments
		`).
		From("businesses b").
		LeftJoin("business_attachment AS ba ON b.id = ba.business_id")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.GroupBy("b.id").ToSql() // Group by business ID for proper aggregation
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item           entity.Business
			attachmentsRaw []byte // To hold the aggregated JSON array of attachments
		)

		err = rows.Scan(
			&item.ID, &item.Name, &description, &item.CategoryID, &item.Address,
			&latitude, &longitude, &contactInfo, &hoursOfOperation,
			&item.OwnerID, &createdAt, &updatedAt, &attachmentsRaw,
		)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		if latitude.Valid {
			item.Latitude = latitude.Float64
		}
		if longitude.Valid {
			item.Longitude = longitude.Float64
		}
		if contactInfo.Valid {
			var contactInfoStruct entity.ContactInfo
			err := json.Unmarshal([]byte(contactInfo.String), &contactInfoStruct)
			if err != nil {
				return response, err
			}
			item.ContactInfo = contactInfoStruct
		}
		if hoursOfOperation.Valid {
			var hoursOfOperationStruct entity.HoursOfOperation
			err := json.Unmarshal([]byte(hoursOfOperation.String), &hoursOfOperationStruct)
			if err != nil {
				return response, err
			}
			item.HoursOfOperation = hoursOfOperationStruct
		}
		if description.Valid {
			item.Description = description.String
		}

		err = json.Unmarshal(attachmentsRaw, &item.Attachments)
		if err != nil {
			return response, err
		}

		response.Items = append(response.Items, item)
	}

	// Count query
	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("businesses b").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BusinessRepo) Update(ctx context.Context, req entity.Business) (entity.Business, error) {
	mp := map[string]interface{}{
		"name":               req.Name,
		"description":        req.Description,
		"category_id":        req.CategoryID,
		"address":            req.Address,
		"latitude":           req.Latitude,
		"longitude":          req.Longitude,
		"contact_info":       req.ContactInfo,
		"hours_of_operation": req.HoursOfOperation,
		"owner_id":           req.OwnerID,
		"updated_at":         "now()",
	}

	qeury, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("businesses").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BusinessRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
