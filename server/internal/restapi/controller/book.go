package controller

import (
	"context"
	"net/http"

	"github.com/Sunaga0709/go-api-template/internal/restapi/controller/openapitype"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
	"github.com/Sunaga0709/go-api-template/internal/restapi/middleware"
)

func (c *Controller) GetV1Books(
	ctx context.Context,
	_ gen.GetV1BooksRequestObject,
) (gen.GetV1BooksResponseObject, error) {
	result, err := c.bookQueryService.List(ctx)
	if err != nil {
		return nil, middleware.NewErrorDetail(http.StatusInternalServerError, err)
	}

	books := make([]gen.Book, 0, len(result))
	for _, v := range result {
		books = append(books, gen.Book{
			BookId:          v.BookID,
			Title:           v.Title,
			Summary:         v.Summary,
			Author:          v.Author,
			Price:           v.Price,
			PublicationDate: openapitype.DateToOpenAPI(v.PublicationDate),
		})
	}

	return gen.GetV1Books200JSONResponse{
		GetV1BooksResponseJSONResponse: gen.GetV1BooksResponseJSONResponse{
			Books: books,
		},
	}, nil
}
