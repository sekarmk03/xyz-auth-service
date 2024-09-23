package handler

import (
	"xyz-auth-service/common/config"
	"xyz-auth-service/pb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	config config.Config
	userSvc 