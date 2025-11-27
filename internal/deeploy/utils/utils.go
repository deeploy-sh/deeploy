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

func Request(method, url string, data any) (*http.Response, error) {
	config, err := config.Load()
	if err != nil {
		return nil, err
	}

	var jsonData []byte

	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	r, err := http.NewRequest(method, config.Server+"/api"+url, bytes.NewBuffer(jsonData))
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

func IsOnline() bool {
	endpoints := []string{
		"https://www.google.com",
		"https://1.1.1.1", // Cloudflare
		"https://8.8.8.8", // Google DNS
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
	// TODO: probably should check errors here
	c, _ := config.Load()
	c.Token = ""
	_ = config.Save(c)
}

func DeleteCfgServer() {
	// TODO: probably should check errors here
	c, _ := config.Load()
	c.Server = ""
	_ = config.Save(c)
}
