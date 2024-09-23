package auth

import (
	"xyz-auth-service/common/config"
	commonJwt "xyz-auth-service/common/jwt"
	"xyz-auth-service/modules/auth/internal/builder"
	"xyz-auth-service/pb"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitGrpc(server *grpc.Server, cfg config.Config, db *gorm.DB, jwtManager *commonJwt.JWT, grpcConn *grpc.ClientConn) {
	auth := builder.BuildAuthHandler(cfg, db, jwtManager, grpcConn)
	pb.RegisterAuthServiceServer(server, auth)
}
