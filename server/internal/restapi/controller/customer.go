package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice"
	"github.com/Sunaga0709/go-api-template/internal/restapi/controller/openapitype"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
	"github.com/Sunaga0709/go-api-template/internal/restapi/middleware"
)

func (c *Controller) GetV1Customers(
	ctx context.Context,
	request gen.GetV1CustomersRequestObject,
) (gen.GetV1CustomersResponseObject, error) {
	var (
		count  int
		offset int
		order  string
	)
	if request.Params.Count != nil {
		count = *request.Params.Count
	}
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
	}
	if request.Params.Order != nil {
		order = string(*request.Params.Order)
	}

	cs, err := c.customerQueryService.List(ctx, queryservice.ListCustomerInput{
		Count:  count,
		Offset: offset,
		Order:  order,
	})
	if err != nil {
		return nil, middleware.NewErrorDetail(http.StatusInternalServerError, err)
	}

	customers := make([]gen.Customer, 0, len(cs.Customers))
	for _, v := range cs.Customers {
		customers = append(customers, gen.Customer{
			CustomerId: v.CustomerID,
			Nickname:   v.Nickname,
			Birthday:   v.Birthday.String(),
			Location:   v.Location,
		})
	}

	return gen.GetV1Customers200JSONResponse{
		GetV1CustomersResponseJSONResponse: gen.GetV1CustomersResponseJSONResponse{
			Customers: customers,
		},
	}, nil
}

//nolint:revive // generated openapi
func (c *Controller) GetV1CustomersCustomerId(
	ctx context.Context,
	request gen.GetV1CustomersCustomerIdRequestObject,
) (gen.GetV1CustomersCustomerIdResponseObject, error) {
	output, err := c.customerQueryService.Get(ctx, queryservice.GetCustomerInput{
		CustomerID: request.CustomerId,
		Count:      request.Params.Count,
		Offset:     request.Params.Offset,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if qsErr, ok := errors.AsType[queryservice.Error](err); ok && qsErr.Kind == queryservice.ErrorKindNotFound {
			status = http.StatusNotFound
		}

		return nil, middleware.NewErrorDetail(status, err)
	}

	books := make([]gen.Book, 0, len(output.FavoriteBooks))
	for _, v := range output.FavoriteBooks {
		books = append(books, gen.Book{
			BookId:          v.BookID,
			Title:           v.Title,
			Summary:         v.Summary,
			Author:          v.Author,
			Price:           v.Price,
			PublicationDate: openapitype.DateToOpenAPI(v.PublicationDate),
		})
	}

	return gen.GetV1CustomersCustomerId200JSONResponse{
		GetV1CustomersCustomerIdResponseJSONResponse: gen.GetV1CustomersCustomerIdResponseJSONResponse{
			CustomerId:    output.CustomerID,
			Nickname:      output.Nickname,
			Birthday:      output.Birthday,
			Location:      output.Location,
			FavoriteBooks: books,
		},
	}, nil
}

func (c *Controller) PostV1Customers(
	ctx context.Context,
	_ gen.PostV1CustomersRequestObject,
) (gen.PostV1CustomersResponseObject, error) {
	output, err := c.customerUsecase.Create(ctx)
	if err != nil {
		return nil, middleware.NewErrorDetail(http.StatusInternalServerError, err)
	}

	return gen.PostV1Customers200JSONResponse{
		PostV1CustomersResponseJSONResponse: gen.PostV1CustomersResponseJSONResponse{
			CustomerId: output.CustomerID.String(),
		},
	}, nil
}

//nolint:revive // generated openapi
func (c *Controller) PutV1CustomersCustomerId(
	ctx context.Context,
	request gen.PutV1CustomersCustomerIdRequestObject,
) (gen.PutV1CustomersCustomerIdResponseObject, error) {
	if err := c.customerUsecase.Update(ctx, usecase.UpdateInput{
		CustomerID: request.CustomerId,
		Nickname:   request.Body.Nickname,
		Birthday:   request.Body.Birthday,
		Location:   request.Body.Location,
	}); err != nil {
		status := http.StatusInternalServerError

		if ucErr, ok := errors.AsType[usecase.Error](err); ok {
			switch ucErr.Kind {
			case usecase.ErrorKindInvalidValue:
				status = http.StatusBadRequest
			case usecase.ErrorKindNotFound:
				status = http.StatusNotFound
			default:
			}
		}

		return nil, middleware.NewErrorDetail(status, err)
	}

	return gen.PutV1CustomersCustomerId200JSONResponse{
		PutV1CustomersCustomerIdResponseJSONResponse: gen.PutV1CustomersCustomerIdResponseJSONResponse{},
	}, nil
}
