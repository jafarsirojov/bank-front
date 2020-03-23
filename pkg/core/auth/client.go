package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Url string

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// errors are part API
var ErrUnknown = errors.New("unknown error")
var ErrResponse = errors.New("response error")

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	return strings.Join(e.Errors, ", ")
}

// for errors.Is
func (e *ErrorResponse) Unwrap() error {
	return ErrResponse
}

type Client struct {
	url Url
}

func NewClient(url Url) *Client {
	return &Client{url: url}
}

func (c *Client) Login(ctx context.Context, login string, password string) (token string, err error) {
	// add timeout to context
	ctx, _ = context.WithTimeout(ctx, time.Second)

	requestData := TokenRequest{
		Username: login,
		Password: password,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/tokens", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	// in other request
	// request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		// context.Canceled
		// context.DeadlineExceeded
		return "", fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("can't parse response: %w", err)
	}

	switch response.StatusCode {
	case 200:
		var responseData *TokenResponse
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			return "", fmt.Errorf("can't decode response: %w", err)
		}
		return responseData.Token, nil
	case 400:
		var responseData *ErrorResponse
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			return "", fmt.Errorf("can't decode response: %w", err)
		}
		return "", responseData
	default:
		return "", ErrUnknown
	}

}
