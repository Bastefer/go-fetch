package domain

type Brand struct {
	ID   int
	Name string
}

type Category struct {
	ID   int
	Name string
}

type Product struct {
	ID         int
	Name       string
	BrandID    int
	CategoryID int
	Price      int64
	Stock      int
}