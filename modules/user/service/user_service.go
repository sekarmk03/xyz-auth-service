package service

import (
	"context"
	"log"
	"xyz-auth-service/common/config"
	commonErr "xyz-auth-service/common/error"
	"xyz-auth-service/modules/user/entity"
	"xyz-auth-service/modules/user/repository"
)

type UserService struct {
	cfg            config.Config
	userRepository repository.UserRepositoryUseCase
}

type UserServiceUseCase interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, uuid string) (*entity.User, error)
}

func NewUserService(cfg config.Config, userRepository repository.UserRepositoryUseCase) *UserService {
	return &UserService{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

func (svc *UserService) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	res, err := svc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [UserService - FindByEmail] Error while find user by email: ", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *UserService) FindById(ctx context.Context, uuid string) (*entity.User, error) {
	res, err := svc.userRepository.FindById(ctx, uuid)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [UserService - FindById] Error while find user by ID: ", parseError.Message)
		return nil, err
	}

	return res, nil
}
