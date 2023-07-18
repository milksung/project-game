package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GetResult struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

func GetAndAuthorize(url, token string) GetResult {

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	call := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return GetResult{
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
		return GetResult{
			Status: http.StatusInternalServerError,
			Result: err.Error(),
		}
	}

	defer resp.Body.Close()

	var res interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Println(err)
		return GetResult{
			Status: resp.StatusCode,
			Result: res,
		}
	}

	return GetResult{
		Status: resp.StatusCode,
		Result: res,
	}
}
