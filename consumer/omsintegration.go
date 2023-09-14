package consumer

import (
	"encoding/json"
	"log"

	"github.com/mahfuzan/swiftoms-rabbitmq/client/omsclient"
)

type OmsIntegrationService struct {
	omsClient omsclient.Client
}

func NewOmsService(omsClient omsclient.Client) *OmsIntegrationService {
	return &OmsIntegrationService{omsClient: omsClient}
}

const (
	SWIFTOMS_ORDER_QUEUE_NEW string = "swiftoms.order-queue.new"
)

// OmsIntegration process queue message to create order to OMS.
func (s *OmsIntegrationService) OmsIntegration(queueName string, messageBody []byte) (bool, error) {
	// Create a map that maps string keys to functions
	functions := map[string]func([]byte) (bool, error){
		SWIFTOMS_ORDER_QUEUE_NEW: s.createOrder,
	}

	if function, exists := functions[queueName]; exists {
		return function(messageBody)
	} else {
		log.Println("Invalid function name")
		return false, nil
	}
}

func (s *OmsIntegrationService) createOrder(messageBody []byte) (bool, error) {
	// mapping message
	var orderQueueMessage NewOrderQueueMessage
	err := json.Unmarshal(messageBody, &orderQueueMessage)
	if err != nil {
		log.Printf("Fail to unmarshal %v %s", err.Error(), messageBody)
		return false, nil
	}

	// mapping queue message to oms create order request
	orderRequests := make([]omsclient.OmsOrder, 0)
	for _, orderData := range orderQueueMessage.Orders {
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
	// if len(*response) > 0 {
	// 	for _, responseData := range *response {
	// 		log.Printf("Channel Order Increment ID: %v, Channel Code: %v", responseData.Data.ChannelOrderIncrementId, responseData.Data.ChannelCode)
	// 		log.Print(responseData.Message)
	// 	}
	// }

	return true, nil
}
