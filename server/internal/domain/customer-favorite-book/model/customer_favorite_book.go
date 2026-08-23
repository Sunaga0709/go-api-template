package model

type CustomerFavoriteBook struct {
	CustomerID CustomerID
	BookIDs    BookIDs
}

func NewCustomerFavoriteBook(customerID CustomerID, bookIDs BookIDs) *CustomerFavoriteBook {
	return &CustomerFavoriteBook{
		CustomerID: customerID,
		BookIDs:    bookIDs,
	}
}

func (c *CustomerFavoriteBook) Add(bookID BookID) error {
	return c.BookIDs.add(bookID)
}

func (c *CustomerFavoriteBook) Remove(bookID BookID) error {
	return c.BookIDs.remove(bookID)
}
