package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deeploy-sh/deeploy/internal/shared/utils"
)

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrNoDeeployInstance = errors.New("no deeploy instance")
)

func ValidateServer(value string) error {
	if !utils.IsValidURL(value) {
		return ErrInvalidURL
	}

	url := fmt.Sprintf("%s/api/health", value)

	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return ErrInvalidURL
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK ||
		res.Header.Get("Content-Type") != "application/json" {
		return ErrNoDeeployInstance
	}

	var healthCheck struct {
		Service string
		Version string
	}

	if err := json.NewDecoder(res.Body).Decode(&healthCheck); err != nil || healthCheck.Service != "deeploy" {
		return ErrNoDeeployInstance
	}

	return nil
}
