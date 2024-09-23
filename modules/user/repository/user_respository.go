package repository

import (
	"context"
	"errors"
	"log"
	"xyz-auth-service/modules/user/entity"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserRepositoryUseCase interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, uuid string) (*entity.User, error)
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindByEmail")
	defer span.End()

	var user entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [UserRepository - FindByEmail] Record not found for email", email)
			return nil, status.Errorf(codes.NotFound, "record not found for email %s", email)
		}
		log.Println("ERROR: [UserRepository - FindByEmail] Internal server error:", err)
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindById(ctx context.Context, uuid string) (*entity.User, error) {
	ctxSpan, span := trace.StartSpan(ctx, "UserRepository - FindById")
	defer span.End()

	var user entity.User
	if err := u.db.Debug().WithContext(ctxSpan).Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [UserRepository - FindById] Record not found for id", uuid)
			return nil, status.Errorf(codes.NotFound, "record not found for id %s", uuid)
		}
		log.Println("ERROR: [UserRepository - FindById] Internal server error:", err)
		return nil, err
	}

	return &user, nil
}
