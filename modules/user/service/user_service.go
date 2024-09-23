package service

import (
	"xyz-auth-service/common/config"
	"xyz-auth-service/modules/user/internal/repository"
)

type UserService struct {
	cfg config.Config
	userRepository repository.UserRepositoryUseCase
}

type UserServiceUseCase interface {
	