package ctxkeys

import "context"

type contextKey string

const (
	GitHubStars contextKey = "github_stars"
	URLPathKey  contextKey = "url_path"
)

// URLPath returns the URL path from context
func URLPath(ctx context.Context) string {
	if path, ok := ctx.Value(URLPathKey).(string); ok {
		return path
	}
	return ""
}
