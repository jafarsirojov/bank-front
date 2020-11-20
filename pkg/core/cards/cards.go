package cards

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

type Cards struct {
	Id      int    `json:"id"`
	Number  string `json:"number"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	OwnerID int    `json:"owner_id"`
}
type ModelTransferMoneyCardToCard struct {
	IdCardSender        int    `json:"id_card_sender"`
	NumberCardRecipient string `json:"number_card_recipient"`
	Count               int    `json:"count"`
}

type ModelBlockCard struct {
	Id     int    `json:"id"`
	Number string `json:"number"`
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

type Card struct {
	url Url
}

func NewCard(url Url) *Card {
	return &Card{url: url}
}

func (c *Card) AllCards(ctx context.Context, token string) (model []Cards, err error) {
	ctx, _ = context.WithTimeout(ctx, 55*time.Second)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/cards", c.url),
		bytes.NewBuffer(nil),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Fatalf("can't close response body: %d", err)
		}
	}()
	err = ReadJSONBody2(response, &model)
	if err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}

	switch response.StatusCode {
	case 200:
		return model, nil
	case 400:

		return nil, ErrResponse
	default:
		return nil, ErrUnknown
	}
}

func ReadJSONBody2(response *http.Response, dto interface{}) error {
	if response.Header.Get("Content-Type") != "application/json" {
		return errors.New("error: incorrect Content-Type")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("error")
	}
	defer response.Body.Close()

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return errors.New("error")
	}

	return nil
}

//-----------------------

func (c *Card) Transfer(ctx context.Context, numberCardRecipient string, idCardSender string, count string, token string) (err error) {
	// add timeout to context
	ctx, _ = context.WithTimeout(ctx, 55*time.Second)

	idCardSenderInt, err := strconv.Atoi(idCardSender)
	if err != nil {
		log.Printf("can't atoi: %s", err)
		return
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		log.Printf("can't atoi: %s", err)
		return
	}

	requestData := ModelTransferMoneyCardToCard{
		NumberCardRecipient: numberCardRecipient,
		IdCardSender:        idCardSenderInt,
		Count:               countInt,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/cards/transmoney", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		// context.Canceled
		// context.DeadlineExceeded
		return fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		log.Print("transfer 200 ok")
		return nil
	case 400:
		log.Print("can't bad request 400")
		return fmt.Errorf("can't bad request 400: %w", err)
	case 500:
		log.Print("can't server internal error 500")
		return fmt.Errorf("can't server internal error 500: %w", err)
	default:
		return fmt.Errorf("can't transfer money: %s", response.StatusCode)
	}

}

func (c *Card) AddCard(ctx context.Context, name string, balanceStr string, owneridStr string, token string) (err error) {
	ctx, _ = context.WithTimeout(ctx, 55*time.Second)

	balance, err := strconv.Atoi(balanceStr)
	if err != nil {
		log.Printf("can't atoi: %s", err)
		return
	}

	ownerid, err := strconv.Atoi(owneridStr)
	if err != nil {
		log.Printf("can't atoi: %s", err)
		return
	}

	requestData := Cards{
		Id:      0,
		Name:    name,
		Balance: balance,
		OwnerID: ownerid,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/cards", c.url),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		log.Print("transfer 200 ok")
		return nil
	case 400:
		log.Print("can't bad request 400")
		return fmt.Errorf("can't bad request 400: %w", err)
	case 500:
		log.Print("can't server internal error 500")
		return fmt.Errorf("can't server internal error 500: %w", err)
	default:
		return fmt.Errorf("can't transfer money: %s", response.StatusCode)
	}
}
