package cards

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strings"
	"time"
)

type Url string

type Cards struct {
	Id      int    `json:"id"`
	Number  string `json:"number"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
	OwnerID int64  `json:"owner_id"`
}
type ModelTransferMoneyCardToCard struct {
	IdCardSender        int    `json:"id_card_sender"`
	NumberCardRecipient string `json:"number_card_recipient"`
	Count               int64  `json:"count"`
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

func (c *Card) AllCards(ctx context.Context) (model []Cards, err error) {
	ctx, _ = context.WithTimeout(ctx, time.Second)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/cards", c.url),
		bytes.NewBuffer(nil),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
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
	err = rest.ReadJSONBody(response.Request, &model)
	if err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}

	switch response.StatusCode {
	case 200:
			return model,nil
	case 400:

		return nil, ErrResponse
	default:
		return nil, ErrUnknown
	}
}
