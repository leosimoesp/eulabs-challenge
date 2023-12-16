package repository

import "context"

type ProductRepositoryData struct {
	Title        string
	Description  string
	Code         string
	Reference    string
	CreatedAt    string
	UpdatedAt    string
	PriceInCents int64
	ID           int64
}

type ProductRepositoryInput struct {
	Title        string
	Description  string
	Code         string
	Reference    string
	PriceInCents int64
}

type ProductRepository interface {
	Insert(ctx context.Context, in ProductRepositoryInput) (ProductRepositoryData, error)
	GetByCode(ctx context.Context, code string) (ProductRepositoryData, error)
}
