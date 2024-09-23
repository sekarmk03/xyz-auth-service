package builder

import (
	"xyz-auth-service/common/config"
	commonJwt "xyz-auth-service/common/jwt"
	"xyz-auth-service/modules/auth/internal/handler"
	userRepo "xyz-auth-service/modules/user/repository"
	userSvc "xyz-auth-service/modules/user/service"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func BuildAuthHandler(cfg config.Config, db *gorm.DB, jwtManager *commonJwt.JWT, grpcConn *grpc.ClientConn) *handler.AuthHandler {
	userRepository := userRepo.NewUserRepository(db)
	userSvc := userSvc.NewUserService(cfg, userRepository)

	return handler.NewAuthHandler(cfg, userSvc, jwtManager)
}
