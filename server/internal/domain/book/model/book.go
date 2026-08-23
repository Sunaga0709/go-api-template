package model

type Book struct {
	BookID          BookID
	Title           BookTitle
	Summary         BookSummary
	Author          BookAuthor
	Price           BookPrice
	PublicationDate BookPublicationDate
}

func NewBook(
	bookID BookID,
	title BookTitle,
	summary BookSummary,
	author BookAuthor,
	price BookPrice,
	publicationDate BookPublicationDate,
) *Book {
	return &Book{
		BookID:          bookID,
		Title:           title,
		Summary:         summary,
		Author:          author,
		Price:           price,
		PublicationDate: publicationDate,
	}
}
