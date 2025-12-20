package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/shared/errs"
)

// writeError logs the full error and sends a JSON response to the client.
// User-facing errors (ErrNotFound, etc.) show details, internal errors are hidden.
func writeError(w http.ResponseWriter, err error) {
	slog.Error("request failed", "error", err)

	status := http.StatusInternalServerError
	msg := "Internal error"

	switch {
	case errors.Is(err, errs.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, errs.ErrUnauthorized):
		status, msg = http.StatusUnauthorized, "Unauthorized"
	case errors.Is(err, errs.ErrInvalidCredentials):
		status, msg = http.StatusUnauthorized, "Invalid credentials"
	case errors.Is(err, errs.ErrInvalidInput):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, errs.ErrConflict), errors.Is(err, errs.ErrDuplicateEmail):
		status, msg = http.StatusConflict, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
