package domain

type Client struct {
	ID        int
	FirstName string
	LastName  string
}

type ClientProduct struct {
	ClientID  int
	ProductID int
}


