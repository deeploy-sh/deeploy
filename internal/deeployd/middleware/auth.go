package middleware

import (
	"context"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/cookie"
	"github.com/deeploy-sh/deeploy/internal/deeployd/jwt"
	"github.com/deeploy-sh/deeploy/internal/deeployd/ui/pages"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
)

type AuthMiddleWare struct {
	userService service.UserServiceInterface
}

func NewAuthMiddleware(userService service.UserServiceInterface) *AuthMiddleWare {
	return &AuthMiddleWare{userService: userService}
}

func getToken(r *http.Request) string {
	// CLI token
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	// Web token
	return cookie.GetTokenFromCookie(r)
}

func (m *AuthMiddleWare) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		t, claims, err := jwt.ValidateToken(token)
		if err != nil || !t.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"].(string)
		user, err := m.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user.ToDTO())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAuth(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isCLI := r.URL.Query().Get("cli") == "true"

		token := getToken(r)
		if token == "" {
			path := "/"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}

			// CLI Auth need CLI params
			if isCLI {
				if r.URL.RawQuery != "" {
					path += "?" + r.URL.RawQuery // Behalte cli=true&port=xyz
				}
			}

			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RequireGuest(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isCLI := r.URL.Query().Get("cli") == "true"

		token := getToken(r)
		if token != "" {
			// CLI flow
			if isCLI {
				pages.CliAuthSuccess(
					r.URL.Query().Get("port"),
					token,
				).Render(r.Context(), w)
				return
			}

			// Web flow
			path := "/dashboard"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}
			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
