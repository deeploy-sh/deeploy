package auth

import (
	"context"

	"github.com/axadrn/deeploy/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUser(ctx context.Context) *data.UserDTO {
	user, ok := ctx.Value("user").(*data.UserDTO)
	if ok {
		return user
	}
	return nil
}

func IsAuthenticated(ctx context.Context) bool {
	return GetUser(ctx) != nil
}

func IsOwner(dataUserID string, ctx context.Context) bool {
	return dataUserID == GetUser(ctx).ID
}
