package controller

import (
	favoritebookuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	customeruc "github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/commandservice"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice"
)

type Controller struct {
	customerQueryService queryservice.CustomerQueryService
	customerUsecase      customeruc.CustomerUsecase

	customerFavoriteBookCommandService commandservice.CustomerFavoriteBookCommandService
	customerFavoriteBookUsecase        favoritebookuc.CustomerFavoriteBookUsecase

	bookQueryService queryservice.BookQueryService
}

func NewController(
	customerQueryService queryservice.CustomerQueryService,
	customerUsecase customeruc.CustomerUsecase,
	customerFavoriteBookCommandService commandservice.CustomerFavoriteBookCommandService,
	customerFavoriteBookUsecase favoritebookuc.CustomerFavoriteBookUsecase,
	bookQueryService queryservice.BookQueryService,
) *Controller {
	return &Controller{
		customerQueryService:               customerQueryService,
		customerUsecase:                    customerUsecase,
		customerFavoriteBookCommandService: customerFavoriteBookCommandService,
		customerFavoriteBookUsecase:        customerFavoriteBookUsecase,
		bookQueryService:                   bookQueryService,
	}
}
