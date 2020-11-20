package history

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Url string

type History struct {
	url Url
}

func NewHistory(url Url) *History {
	return &History{url: url}
}

var ErrUnknown = errors.New("unknown error")
var ErrResponse = errors.New("response error")

func (c *History) AllHistory(ctx context.Context, token string) (model []ModelOperationsLog, err error) {
	log.Print("start func allhistory")
	ctx, _ = context.WithTimeout(ctx, 6666*time.Second)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/history", c.url),
		bytes.NewBuffer(nil),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s",token))
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
		log.Print("allhistory 200 ok")
		return model,nil
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

type ModelOperationsLog struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	RecipientSender string `json:"recipientsender"`
	Count           int64  `json:"count"`
	BalanceOld      int64  `json:"balanceold"`
	BalanceNew      int64  `json:"balancenew"`
	Time            int64  `json:"time"`
	OwnerID         int64  `json:"ownerid"`
}
