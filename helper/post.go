package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PostResult struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

func Post(url string, body interface{}) (interface{}, error) {

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	jsonBody, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(jsonBody)

	call := http.Client{}
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := call.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func PostAndAuthorize(url, token string, body interface{}) PostResult {

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	jsonBody, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(jsonBody)

	call := http.Client{}
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		fmt.Println(err)
		return PostResult{
			Status: http.StatusInternalServerError,
			Result: err.Error(),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := call.Do(req)
	if err != nil {
		fmt.Println(err)
		return PostResult{
			Status: http.StatusInternalServerError,
			Result: err.Error(),
		}
	}

	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Println(err)
		return PostResult{
			Status: resp.StatusCode,
			Result: res,
		}
	}

	return PostResult{
		Status: resp.StatusCode,
		Result: res,
	}
}
