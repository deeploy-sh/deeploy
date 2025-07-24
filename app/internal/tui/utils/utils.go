package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/axadrn/deeploy/internal/errs"
	"github.com/axadrn/deeploy/internal/tui/config"
)

type RequestProps struct {
	Method string
	URL    string
	Data   interface{}
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

	r, err := http.NewRequest(p.Method, "http://"+config.Server+"/api"+p.URL, bytes.NewBuffer(jsonData))
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
