package di

import (
	"github.com/Sunaga0709/go-api-template/internal/database"
)

type DI struct {
	CommandService CommandService
	QueryService   QueryService
	Usecase        Usecase
}

func New(
	executorProvider database.ExecutorProvider,
	txManager database.TxManager,
) DI {
	uc := NewUsecase(executorProvider, txManager)
	queryService := NewQueryService(executorProvider)
	commandService := NewCommandService(uc.CustomerUsecase, uc.BookUsecase, uc.CustomerFavoriteBookUsecase)

	return DI{
		CommandService: commandService,
		QueryService:   queryService,
		Usecase:        uc,
	}
}
