package mock

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
)

// mock service can be used to mock hanlder's service
type Service struct {
	repository product.Repository
}

func NewService(repository product.Repository) *Service {
	return &Service{repository: repository}
}
