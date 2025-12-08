package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/shared/utils"
)

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrNoDeeployInstance = errors.New("no deeploy instance")
)

func IsOnline() bool {
	endpoints := []string{
		"https://www.google.com",
		"https://1.1.1.1",
		"https://8.8.8.8",
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest(http.MethodHead, endpoint, nil)
		_, err := client.Do(req)
		if err == nil {
			return true
		}
	}
	return false
}

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

func DeleteCfgToken() {
	c, _ := config.Load()
	c.Token = ""
	_ = config.Save(c)
}

func DeleteCfgServer() {
	c, _ := config.Load()
	c.Server = ""
	_ = config.Save(c)
}
