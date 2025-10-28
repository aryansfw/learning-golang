package domain

type Product struct {
	ID    int64
	Name  string
	Price	int 
	Stock int

	UserID int64 // Who this product belongs to
}
