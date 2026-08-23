package controller

import (
	"context"
	"net/http"

	cfbuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/commandservice"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
	"github.com/Sunaga0709/go-api-template/internal/restapi/middleware"
)

//nolint:revive // generated openapi
func (c *Controller) PostV1CustomersCustomerIdFavoriteBooks(
	ctx context.Context,
	req gen.PostV1CustomersCustomerIdFavoriteBooksRequestObject,
) (gen.PostV1CustomersCustomerIdFavoriteBooksResponseObject, error) {
	if req.Body == nil {
		return nil, requestBodyRequiredError()
	}

	if err := c.customerFavoriteBookCommandService.Favorite(ctx, commandservice.FavoriteInput{
		CustomerID: req.CustomerId,
		BookID:     req.Body.BookId,
	}); err != nil {
		if cerr, ok := commandservice.ParseError(err); ok {
			var status int
			switch cerr.Kind {
			case commandservice.ErrorKindInvalidValue:
				status = http.StatusBadRequest
			case commandservice.ErrorKindNotFound:
				status = http.StatusNotFound
			case commandservice.ErrorKindAlreadyExists:
				status = http.StatusConflict
			default:
				status = http.StatusInternalServerError
			}

			e := middleware.NewErrorDetail(status, err)
			if cerr.ClientMessage != nil {
				return nil, e.WithClientMessage(*cerr.ClientMessage)
			}

			return nil, e
		}

		return nil, unknownError(err)
	}

	return gen.PostV1CustomersCustomerIdFavoriteBooks200JSONResponse{
		PostV1CustomersCustomerIdFavoriteBooksResponseJSONResponse: gen.PostV1CustomersCustomerIdFavoriteBooksResponseJSONResponse{},
	}, nil
}

//nolint:revive // generated openapi
func (c *Controller) PutV1CustomersCustomerIdFavoriteBooks(
	ctx context.Context,
	req gen.PutV1CustomersCustomerIdFavoriteBooksRequestObject,
) (gen.PutV1CustomersCustomerIdFavoriteBooksResponseObject, error) {
	if req.Body == nil {
		return nil, requestBodyRequiredError()
	}

	if err := c.customerFavoriteBookUsecase.Unfavorite(ctx, cfbuc.UnfavoriteInput{
		CustomerID: req.CustomerId,
		BookID:     req.Body.BookId,
	}); err != nil {
		if uerr, ok := cfbuc.ParseError(err); ok {
			var status int
			switch uerr.Kind {
			case cfbuc.ErrorKindInvalidValue:
				status = http.StatusBadRequest
			case cfbuc.ErrorKindNotFavorited:
				status = http.StatusNotFound
			default:
				status = http.StatusInternalServerError
			}

			e := middleware.NewErrorDetail(status, err)
			if uerr.ClientMessage != nil {
				return nil, e.WithClientMessage(*uerr.ClientMessage)
			}

			return nil, e
		}

		return nil, unknownError(err)
	}

	return gen.PutV1CustomersCustomerIdFavoriteBooks200JSONResponse{
		PutV1CustomersCustomerIdFavoriteBooksResponseJSONResponse: gen.PutV1CustomersCustomerIdFavoriteBooksResponseJSONResponse{},
	}, nil
}
