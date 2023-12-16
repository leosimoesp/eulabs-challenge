package repository

import "context"

type ProductRepositoryData struct {
	Reference string
	CreatedAt string
	ID        int64
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
}
