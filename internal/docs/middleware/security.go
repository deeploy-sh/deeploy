package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type contextKey string

const nonceKey contextKey = "csp-nonce"

// GenerateNonce creates a cryptographically secure random nonce
func GenerateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// GetNonce retrieves the nonce from context
func GetNonce(ctx context.Context) string {
	if nonce, ok := ctx.Value(nonceKey).(string); ok {
		return nonce
	}
	return ""
}

// NonceMiddleware generates a unique nonce for each request
func NonceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := GenerateNonce()
		ctx := context.WithValue(r.Context(), nonceKey, nonce)
		ctx = templ.WithNonce(ctx, nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := GetNonce(r.Context())

		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=()")

		cspPolicy := strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'nonce-" + nonce + "' https://unpkg.com https://cdnjs.cloudflare.com https://plausible.axeladrian.com",
			"style-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com",
			"img-src 'self' data: https:",
			"font-src 'self' data:",
			"connect-src 'self' https://plausible.axeladrian.com",
			"frame-ancestors 'none'",
			"base-uri 'self'",
			"form-action 'self'",
		}, "; ")

		w.Header().Set("Content-Security-Policy", cspPolicy)

		next.ServeHTTP(w, r)
	})
}
