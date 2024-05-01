package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"github.com/lbsti/eulabs-challenge/internal/core/entity"
	"github.com/lbsti/eulabs-challenge/internal/core/repository"
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
		slog.Error("impossible insert product", slog.Any("msg", err))
		return repository.ProductRepositoryData{}, err
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		slog.Error("impossible to retrieve last inserted product id", slog.Any("msg", err))
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

	if err := r.db.QueryRowContext(ctx, query, codeLowerCase).Scan(&id, &title,
		&description, &price, &reference, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return repository.ProductRepositoryData{}, entity.ProductNotFoundErr
		}
		slog.Error("impossible to retrieve product", slog.Any("msg", err))
		return repository.ProductRepositoryData{}, err
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
	result, err := r.db.ExecContext(ctx, query, codeLowerCase)

	if err != nil {
		slog.Error("impossible to delete product", slog.Any("msg", err))
		return false, err
	}

	if affectedRows, err := result.RowsAffected(); err == nil {
		if affectedRows > 0 {
			return affectedRows > 0, nil
		}
		return false, entity.ProductNotFoundErr
	} else {
		slog.Error("impossible to delete product", slog.Any("msg", err))
	}

	return true, nil
}

func (r ProductRepositorySQL) Update(ctx context.Context,
	in repository.ProductRepositoryInput) error {
	codeWithoutSpace := strings.ReplaceAll(in.Code, " ", "")
	codeLowerCase := strings.ToLower(codeWithoutSpace)

	query := `UPDATE products SET title = ?, description = ?, reference = ?,
	 price_in_cents = ? 
	 WHERE LOWER(code) = ?`

	result, err := r.db.ExecContext(ctx, query, in.Title, in.Description,
		in.Reference, in.PriceInCents, codeLowerCase)

	if affectedRows, err := result.RowsAffected(); err == nil {
		if affectedRows > 0 {
			return nil
		}
	} else {
		slog.Error("impossible to update product", slog.Any("msg", err))
	}
	return err
}
