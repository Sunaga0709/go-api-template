package di

import (
	"github.com/Sunaga0709/go-api-template/internal/database"
	bookinfra "github.com/Sunaga0709/go-api-template/internal/domain/book/infra"
	bookuc "github.com/Sunaga0709/go-api-template/internal/domain/book/usecase"
	cfbinfra "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/infra"
	cfbuc "github.com/Sunaga0709/go-api-template/internal/domain/customer-favorite-book/usecase"
	customerinfra "github.com/Sunaga0709/go-api-template/internal/domain/customer/infra"
	customeruc "github.com/Sunaga0709/go-api-template/internal/domain/customer/usecase"
)

type Usecase struct {
	CustomerUsecase             customeruc.CustomerUsecase
	BookUsecase                 bookuc.BookUsecase
	CustomerFavoriteBookUsecase cfbuc.CustomerFavoriteBookUsecase
}

func NewUsecase(
	executorProvider database.ExecutorProvider,
	txManager database.TxManager,
) Usecase {
	customerUsecase := customeruc.NewCustomerUsecase(customerinfra.NewCustomerRepository(), executorProvider)
	bookUsecase := bookuc.NewBookUsecase(bookinfra.NewBookRepository(), executorProvider)
	customerFavoriteBookUsecase := cfbuc.NewCustomerFavoriteBookUsecase(cfbinfra.NewCustomerFavoriteBookRepository(), txManager)

	return Usecase{
		CustomerUsecase:             customerUsecase,
		BookUsecase:                 bookUsecase,
		CustomerFavoriteBookUsecase: customerFavoriteBookUsecase,
	}
}
