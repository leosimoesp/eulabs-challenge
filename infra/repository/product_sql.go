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

func (r ProductRepositorySQL) GetByCode(ctx context.Context,
	code string) (repository.ProductRepositoryData, error) {

	query := `SELECT p.id, p.title, p.description, p.price_in_cents,
	p.reference, 
	CAST(p.created_at AS CHAR) created_at,
	CAST(p.updated_at AS CHAR) updated_at 
	FROM products p WHERE LOWER(p.code) = ?`

	codeWithoutSpace := strings.ReplaceAll(code, " ", "")
	codeLowerCase := strings.ToLower(codeWithoutSpace)

	var id, price int64
	var title, description, reference, createdAt, updatedAt string

	if e := r.db.QueryRowContext(ctx, query, codeLowerCase).Scan(&id, &title,
		&description, &price, &reference, &createdAt, &updatedAt); e != nil {
		if e == sql.ErrNoRows {
			return repository.ProductRepositoryData{}, entity.ProductNotFoundErr
		}
		log.Default().Printf("impossible to retrieve product: %s", e)
		return repository.ProductRepositoryData{}, e
	}

	return repository.ProductRepositoryData{
		Title:        title,
		Description:  description,
		Code:         codeWithoutSpace,
		Reference:    reference,
		PriceInCents: price,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		ID:           id,
	}, nil
}

func (r ProductRepositorySQL) DeleteByCode(ctx context.Context, code string) (bool, error) {
	codeWithoutSpace := strings.ReplaceAll(code, " ", "")
	codeLowerCase := strings.ToLower(codeWithoutSpace)

	query := `DELETE FROM products WHERE LOWER(code) = ?`
	result, e := r.db.ExecContext(ctx, query, codeLowerCase)

	if affectedRows, err := result.RowsAffected(); err == nil {
		if affectedRows > 0 {
			return affectedRows > 0, nil
		}
		return false, entity.ProductNotFoundErr
	}
	return e == nil, e
}
