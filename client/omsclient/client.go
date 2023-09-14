package omsclient

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	CREATE_ORDER_URL = "/rest/V1/order/create/"
)

type Client interface {
	CreateOrder(request CreateOrderRequest) (*CreateOrderResponse, error)
}

type client struct {
	restyClient *resty.Client
}

func NewClient(baseUrl, token string) Client {
	return &client{
		restyClient: resty.New().
			SetBaseURL(baseUrl).
			SetAuthToken(token),
	}
}

func (c *client) CreateOrder(request CreateOrderRequest) (*CreateOrderResponse, error) {
	resp, err := c.restyClient.R().
		SetBody(request).
		SetResult(&CreateOrderResponse{}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Post(CREATE_ORDER_URL)
	if err != nil {
		err = errors.Wrap(err, "Execute API call")
		return nil, err
	}

	log.Printf("URL: %s", resp.Request.URL)
	log.Printf("API call response time: %v seconds", resp.Time().Seconds())

	respModel, ok := resp.Result().(*CreateOrderResponse)
	if !ok {
		err = errors.New("Can't parse expected result")
		return nil, err
	}
	return respModel, nil
}
