package consumer

import (
	"log"

	"github.com/mahfuzan/swiftoms-rabbitmq/client/omsclient"
	"github.com/mahfuzan/swiftoms-rabbitmq/config"
)

type OmsIntegrationService struct {
	conf      config.Config
	omsClient omsclient.Client
}

func NewOmsService(
	conf config.Config,
	omsClient omsclient.Client) (*OmsIntegrationService, error) {
	return &OmsIntegrationService{
		conf:      conf,
		omsClient: omsClient}, nil
}

// OmsIntegration process queue message to create order to OMS.
func (s *OmsIntegrationService) OmsIntegration(orderMessage NewOrderQueueMessage) (bool, error) {
	// mapping queue message to oms create order request
	orderRequests := make([]omsclient.OmsOrder, 0)
	for _, orderData := range orderMessage.Orders {
		orderRequests = append(orderRequests, omsclient.OmsOrder{
			ChannelOrderIncrementId: orderData.ChannelOrderIncrementId,
			ChannelCode:             orderData.ChannelCode,
		})
	}

	// create order to OMS
	response, err := s.omsClient.CreateOrder(omsclient.CreateOrderRequest{
		Orders: orderRequests,
	})
	if err != nil {
		return false, err
	}
	if response == nil {
		return false, nil
	}

	// process response
	if len(*response) > 0 {
		for _, responseData := range *response {
			log.Printf("Channel Order Increment ID: %v, Channel Code: %v", responseData.Data.ChannelOrderIncrementId, responseData.Data.ChannelCode)
			log.Print(responseData.Message)
		}
	}

	return true, nil
}
