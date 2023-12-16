package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/repository"
)

type ProductRepositorySQL struct {
	db *sql.DB
}

func NewProductRepositorySQL(db *sql.DB) repository.ProductRepository {
	return ProductRepositorySQL{
		db: db,
	}
}

func (r ProductRepositorySQL) Insert(
	ctx context.Context,
	in repository.ProductRepositoryInput) (repository.ProductRepositoryData, error) {

	datetime := time.Now()
	datetimeFmt := datetime.Format(time.RFC3339)

	query := `INSERT INTO products (title, description, code, reference, price_in_cents, 
	created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	insertResult, err := r.db.ExecContext(ctx, query, in.Title, in.Description,
		in.Code, in.Reference, in.PriceInCents, datetimeFmt, datetimeFmt)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return repository.ProductRepositoryData{}, entity.DuplicatedProductCodeErr
		}
		log.Default().Printf("impossible insert product: %s", err)
		return repository.ProductRepositoryData{}, err
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Default().Printf("impossible to retrieve last inserted product id: %s", err)
	}

	return repository.ProductRepositoryData{
		ID:        id,
		Reference: in.Reference,
		CreatedAt: datetimeFmt,
	}, nil
}
