package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"xyz-auth-service/common/config"
	commonErr "xyz-auth-service/common/error"
	commonJwt "xyz-auth-service/common/jwt"
	"xyz-auth-service/common/utils"
	"xyz-auth-service/modules/user/service"
	"xyz-auth-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	config     config.Config
	userSvc    service.UserServiceUseCase
	jwtManager *commonJwt.JWT
}

func NewAuthHandler(config config.Config, userSvc service.UserServiceUseCase, jwtManager *commonJwt.JWT) *AuthHandler {
	return &AuthHandler{
		config:     config,
		userSvc:    userSvc,
		jwtManager: jwtManager,
	}
}

func (ah *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := ah.userSvc.FindByEmail(ctx, req.Email)
	if err != nil {
		if user == nil {
			log.Println("ERROR: [AuthHandler - Login] User not found for email", req.Email)
			return &pb.LoginResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "User not found for email " + req.Email,
			}, status.Errorf(codes.NotFound, "user not found for email %s", req.Email)
		}
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [AuthHandler - Login] Error while find user by email: ", parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(codes.Internal, parseError.Message)
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		log.Println("WARNING: [AuthHandler - LoginUser] Invalid credentials")
		return &pb.LoginResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "invalid credentials",
		}, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	token, err := ah.jwtManager.GenerateToken(user.Uuid, user.Role)

	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUser] Error while generating token:", parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: "token failed to generate: " + parseError.Message,
		}, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
	}

	return &pb.LoginResponse{
		Code:    uint32(http.StatusOK),
		Message: "login success",
		Token:   token,
	}, nil
}

func (ah *AuthHandler) GetCurrentUser(ctx context.Context, req *emptypb.Empty) (*pb.UserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Metadata not found")
		return &pb.UserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "metadata not found",
		}, status.Errorf(codes.InvalidArgument, "metadata not found")
	}

	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Authorization token not found")
		return &pb.UserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "authorization token not found",
		}, status.Errorf(codes.InvalidArgument, "authorization token not found")
	}

	authHeader := values[0]
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Invalid authorization header")
		return &pb.UserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "invalid authorization header",
		}, status.Errorf(codes.InvalidArgument, "invalid authorization header")
	}

	accessToken := parts[1]
	claims, err := ah.jwtManager.Verify(accessToken)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Error while verifying token:", parseError.Message)
		return &pb.UserResponse{
			Code:    uint32(http.StatusUnauthorized),
			Message: parseError.Message,
		}, status.Errorf(codes.Unauthenticated, parseError.Message)
	}

	user, err := ah.userSvc.FindById(ctx, claims.Cred)
	if err != nil {
		if user == nil {
			log.Println("ERROR: [AuthHandler - GetCurrentUser] User not found for ID", claims.Cred)
			return &pb.UserResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "user not found for ID " + claims.Cred,
			}, status.Errorf(codes.NotFound, "user not found for id %s", claims.Cred)
		}
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Error while find user by ID: ", parseError.Message)
		return &pb.UserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(codes.Internal, parseError.Message)
	}

	userProto := &pb.User{
		Uuid:  user.Uuid,
		Email: user.Email,
		Role:  user.Role,
	}

	return &pb.UserResponse{
		Code:    uint32(http.StatusOK),
		Message: "get current user success",
		Data:    userProto,
	}, nil
}
