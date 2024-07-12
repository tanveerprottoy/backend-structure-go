package mock

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
)

// mock service can be used to mock hanlder's service
type Service struct {
	repository user.Repository
}

func NewService(repository user.Repository) *Service {
	return &Service{repository: repository}
}
