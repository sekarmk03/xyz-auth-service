package entity

import (
	"time"
	"xyz-auth-service/pb"
)

const (
	UserTableName = "users"
)

type User struct {
	Uuid      string    `json:"uuid"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      uint32    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserEntity(email, password string, role uint32) *User {
	return &User{
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) TableName() string {
	return UserTableName
}

func ConvertEntityToProto(u *User) *pb.User {
	return &pb.User{
		Uuid:      u.Uuid,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
