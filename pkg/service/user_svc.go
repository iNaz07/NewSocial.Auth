package service

import (
	"auth/domain"
	"auth/ent"
	"auth/ent/user"
	utils "auth/utils"
	"context"
	"net/http"
	"time"
)

type userUsecase struct {
	db             *ent.Client
	jwt            domain.JwtService
	timeoutContext time.Duration
}

func NewUserService(repo *ent.Client, jwt domain.JwtService, time time.Duration) domain.UserService {
	return &userUsecase{
		db:             repo,
		jwt:            jwt,
		timeoutContext: time}
}

func (u *userUsecase) CreateUser(ctx context.Context, newUser *domain.User) error {

	_, err := u.db.User.Query().Where(user.Username(newUser.Username)).Only(ctx)
	if err != nil {
		if err := utils.ValidateCreds(newUser.Username, newUser.Password); err != nil {
			return &domain.LogError{err.Error(), err, http.StatusBadRequest}
		}
		hashedPassword := utils.GenerateHash(newUser.Password)
		newUser.Password = hashedPassword
		newUser.RegisterDate = time.Now()

		_, err := u.db.User.Create().SetEmail(newUser.Email).SetUsername(newUser.Username).SetPassword(newUser.Password).Save(ctx)
		if err != nil {
			return &domain.LogError{"registration error", err, http.StatusInternalServerError}
		}
		return nil
	}
	return &domain.LogError{"user already registered by iin", err, http.StatusBadRequest}
}

func (u *userUsecase) ConfirmUser(ctx context.Context, token string) (domain.JwtToken, error) {
	//TODO implement me что он делает?
	panic("implement me")
}

func (u *userUsecase) LoginUser(ctx context.Context, newUser domain.User) (domain.JwtToken, error) {

	var creds struct {
		id  int64
		psw string
	}

	err := u.db.User.Query().Where(user.Username(newUser.Email)).Select(user.FieldID, user.FieldPassword).Scan(ctx, &creds)
	if err != nil {
		return domain.JwtToken{}, &domain.LogError{"user not found error", err, http.StatusBadRequest}
	}
	if !utils.ComparePasswordHash(newUser.Password, creds.psw) {
		return domain.JwtToken{}, &domain.LogError{"invalid password", err, http.StatusBadRequest}
	}

	token, err := u.jwt.GenerateToken(creds.id, "", newUser.Email)
	if err != nil {
		return domain.JwtToken{}, &domain.LogError{"cannot generate token", err, http.StatusInternalServerError}
	}
	if err := u.jwt.InsertToken(creds.id, token); err != nil {
		return domain.JwtToken{}, &domain.LogError{"cannot insert token", err, http.StatusInternalServerError}
	}
	return domain.JwtToken{
		AccessToken: token,
	}, nil
}

func (u *userUsecase) ForgotPassword(ctx context.Context, email string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUsecase) ResetPassword(ctx context.Context, token, psw string) (domain.JwtToken, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userUsecase) RefreshToken(ctx context.Context, token string) (domain.JwtToken, error) {
	//TODO implement me
	panic("implement me")
}
