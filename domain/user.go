package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	RegisterDate time.Time `json:"registerdate"`
}

//type UserRepository interface {
//	CreateUser(ctx context.Context, user *User) error
//	GetUserByID(ctx context.Context, id int64) (*User, error)
//	GetUserByUsername(ctx context.Context, username string) (*User, error)
//	GetUserByIIN(ctx context.Context, iin string) (*User, error)
//	GetAllUsers(ctx context.Context) ([]User, error)
//	UpgradeUserRepo(ctx context.Context, username string) error
//}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	ConfirmUser(ctx context.Context, token string) (JwtToken, error)
	LoginUser(ctx context.Context, user User) (JwtToken, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, psw string) (JwtToken, error)
	RefreshToken(ctx context.Context, token string) (JwtToken, error)
}
