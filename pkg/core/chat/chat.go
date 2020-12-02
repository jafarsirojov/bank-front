package chat

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

type Chat struct {
	url Url
}

func NewChat(url Url) *Chat {
	return &Chat{url: url}
}

var ErrUnknown = errors.New("unknown error")
var ErrResponse = errors.New("response error")

func (c *Chat) GetAllMessage(ctx context.Context, token string) (model []ModelMassage, err error) {
	log.Print("start func allhistory")
	ctx, _ = context.WithTimeout(ctx, 6666*time.Second)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/chat/message/all", c.url),
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
		log.Print("chat 200 ok")
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

type ModelMassage struct {
	ID            int       `json:"id"`
	SenderID      int       `json:"sender_id"`
	RecipientID   int       `json:"recipient_id"`
	RecipientName string    `json:"recipient_name"`
	Message       string    `json:"message"`
	Time          time.Time `json:"time"`
}