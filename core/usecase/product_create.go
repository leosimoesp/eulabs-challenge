package usecase

import (
	"context"
	"time"

	"github.com/lbsti/eulabs-challenge/core/entity"
	"github.com/lbsti/eulabs-challenge/core/repository"
)

const (
	ProductDefaultTimeout = time.Second * 30
)

type ProductCreate struct {
	repository repository.ProductRepository
}

type ProductInputDTO struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Code         string `json:"code"`
	Reference    string `json:"reference"`
	PriceInCents int64  `json:"priceInCents"`
}

type ProductOutputDTO struct {
	CreatedAt string `json:"createdAt"`
	Reference string `json:"reference"`
	ID        int64  `json:"id"`
}

func NewProductCreate(productRepo repository.ProductRepository) *ProductCreate {
	return &ProductCreate{
		repository: productRepo,
	}
}

func (p *ProductCreate) Execute(ctx context.Context, input ProductInputDTO) (ProductOutputDTO, error) {
	if e := validate(input); e != nil {
		return ProductOutputDTO{}, e
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ProductDefaultTimeout))
	defer cancel()

	productData, err := p.repository.Insert(ctxWithTimeout, repository.ProductRepositoryInput{
		Title:        input.Title,
		Description:  input.Description,
		Code:         input.Code,
		Reference:    input.Reference,
		PriceInCents: input.PriceInCents,
	})

	if err != nil {
		return ProductOutputDTO{}, err
	}

	return ProductOutputDTO{
		CreatedAt: productData.CreatedAt,
		Reference: productData.Reference,
		ID:        productData.ID,
	}, nil
}

func validate(input ProductInputDTO) error {
	product := entity.NewProduct()
	product.Title = input.Title
	product.Description = input.Description
	product.Code = input.Code
	product.Reference = input.Reference
	product.PriceInCents = input.PriceInCents

	if e := product.IsValid(); e != nil {
		return e
	}
	return nil
}
