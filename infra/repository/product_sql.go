package repository

import (
	"database/sql"
	"log"
	"time"

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
	in repository.ProductRepositoryInput) (repository.ProductRepositoryData, error) {

	datetime := time.Now()
	datetimeFmt := datetime.Format(time.RFC3339)

	query := `INSERT INTO product (title, description, code, reference, price, 
	created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	insertResult, err := r.db.Exec(query, in.Title, in.Description, in.Code,
		in.Reference, in.PriceInCents, datetimeFmt, datetimeFmt)

	if err != nil {
		log.Fatalf("impossible insert product: %s", err)
		return repository.ProductRepositoryData{}, err
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted product id: %s", err)
	}

	return repository.ProductRepositoryData{
		ID:        id,
		Reference: in.Reference,
		CreatedAt: datetimeFmt,
	}, nil
}
