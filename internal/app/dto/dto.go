package dto

type SourceProduct struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Brand    string `json:"brand"`
	Category string `json:"category"`
	Price    string `json:"price"`
	Stock    int    `json:"stock"`
}

type SourceClient struct {
	ID        int   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Products  []int `json:"products"`
}

type Stats struct {
	Products  int `json:"products"`
	Clients   int `json:"clients"`
	Brands    int `json:"brands"`
	Categories int `json:"categories"`
}