package auth

import (
	"context"

	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
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

func GetUser(ctx context.Context) *repo.UserDTO {
	user, ok := ctx.Value("user").(*repo.UserDTO)
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
