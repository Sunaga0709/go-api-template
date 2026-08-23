package di

import (
	bookuc "github.com/Sunaga0709/go-api-template/internal/domain/book/usecase"
	cfbuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	customeruc "github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/commandservice"
)

type CommandService struct {
	CustomerFavoriteBookCommandService commandservice.CustomerFavoriteBookCommandService
}

func NewCommandService(
	customerUsecase customeruc.CustomerUsecase,
	bookUsecase bookuc.BookUsecase,
	customerFavoriteBookUsecase cfbuc.CustomerFavoriteBookUsecase,
) CommandService {
	customerFavoriteBookCommandService := commandservice.NewCustomerFavoriteBookCommandService(
		customerUsecase,
		bookUsecase,
		customerFavoriteBookUsecase,
	)

	return CommandService{
		CustomerFavoriteBookCommandService: customerFavoriteBookCommandService,
	}
}
