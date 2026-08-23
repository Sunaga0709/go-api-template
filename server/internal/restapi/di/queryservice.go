package di

import (
	"github.com/Sunaga0709/go-api-template/internal/database"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice"
	"github.com/Sunaga0709/go-api-template/internal/orchestration/queryservice/query/infra"
)

type QueryService struct {
	CustomerQueryService queryservice.CustomerQueryService
	BookQueryService     queryservice.BookQueryService
}

func NewQueryService(executorProvider database.ExecutorProvider) QueryService {
	customerQueryService := queryservice.NewCustomerQueryService(infra.NewCustomerQueryRepository(), executorProvider)
	bookQueryService := queryservice.NewBookQueryService(infra.NewBookQueryRepository(), executorProvider)

	return QueryService{
		CustomerQueryService: customerQueryService,
		BookQueryService:     bookQueryService,
	}
}
