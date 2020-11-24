package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func (c *Client) Register(ctx context.Context, name string, login string, password string, phone string) (err error) {
	ctx, _ = context.WithTimeout(ctx, 55*time.Second)

	phoneInt, err := strconv.Atoi(phone)
	if err != nil {
		log.Printf("can't atoi: %s", err)
		return
	}

	requestData := UserTDO{
		Id:       0,
		Name:     name,
		Login:    login,
		Password: password,
		Phone:    phoneInt,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/users", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Println("Can't http.NewRequestWithContext, client.go 137")
		return fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("Can't http.DefaultClient.Do, client.go 143")
		return fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		log.Print("register 200 ok")
		return nil
	case 400:
		log.Print("can't bad request 400")
		return fmt.Errorf("can't bad request 400: %w", err)
	case 500:
		log.Print("can't server internal error 500")
		return fmt.Errorf("can't server internal error 500: %w", err)
	default:
		return fmt.Errorf("can't register: %s", response.StatusCode)
	}
}

type UserTDO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}
