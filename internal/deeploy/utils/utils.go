package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/utils"
)

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrNoDeeployInstance = errors.New("no deeploy instance")
)

type RequestProps struct {
	Method string
	URL    string
	Data   any
}

func Request(p RequestProps) (*http.Response, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	var jsonData []byte

	if p.Data != nil {
		jsonData, err = json.Marshal(p.Data)
		if err != nil {
			return nil, err
		}
	}

	r, err := http.NewRequest(p.Method, config.Server+"/api"+p.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+config.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrUnauthorized
	}

	return res, nil
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
