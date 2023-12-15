package repository

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
	Insert(in ProductRepositoryInput) (ProductRepositoryData, error)
}
