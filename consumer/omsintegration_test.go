package consumer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mahfuzan/swiftoms-rabbitmq/client/omsclient"
	"github.com/mahfuzan/swiftoms-rabbitmq/consumer"
	omsclient_mock "github.com/mahfuzan/swiftoms-rabbitmq/mock/omsclient"
	"github.com/stretchr/testify/require"
)

const (
	CASE_NORMAL       = "NormalCase"
	CASE_NIL_RESPONSE = "NilResponseCase"
)

func TestOmsIntegration(t *testing.T) {
	testCases := []struct {
		name          string
		order         consumer.NewOrderQueueMessage
		buildStubs    func(mockOmsClient *omsclient_mock.MockClient)
		checkResponse func(resp bool, err error)
	}{
		{
			name: CASE_NORMAL,
			order: consumer.NewOrderQueueMessage{
				Orders: []consumer.Order{
					{
						ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID001",
						ChannelCode:             "TESTCHANNELCODE",
					},
					{
						ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID002",
						ChannelCode:             "TESTCHANNELCODE",
					},
				},
			},
			buildStubs: func(mockOmsClient *omsclient_mock.MockClient) {
				mockOmsClient.EXPECT().
					CreateOrder(gomock.Any()).
					Times(1).
					Return(getCreateOrderResponse(CASE_NORMAL), nil)
			},
			checkResponse: func(resp bool, err error) {
				require.NoError(t, err)
				require.Equal(t, true, resp)
			},
		}, {
			name: CASE_NIL_RESPONSE,
			order: consumer.NewOrderQueueMessage{
				Orders: []consumer.Order{
					{
						ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID001",
						ChannelCode:             "TESTCHANNELCODE",
					},
					{
						ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID002",
						ChannelCode:             "TESTCHANNELCODE",
					},
				},
			},
			buildStubs: func(mockOmsClient *omsclient_mock.MockClient) {
				mockOmsClient.EXPECT().
					CreateOrder(gomock.Any()).
					Times(1).
					Return(getCreateOrderResponse(CASE_NIL_RESPONSE), nil)
			},
			checkResponse: func(resp bool, err error) {
				require.NoError(t, err)
				require.Equal(t, false, resp)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mock_controller := gomock.NewController(t)
			defer mock_controller.Finish()

			mockedOmsClient := omsclient_mock.NewMockClient(mock_controller)

			tc.buildStubs(mockedOmsClient)

			omsService := consumer.NewOmsService(mockedOmsClient)

			isSuccess, err := omsService.OmsIntegration(tc.order)

			tc.checkResponse(isSuccess, err)
		})
	}
}

func getCreateOrderResponse(caseName string) *omsclient.CreateOrderResponse {
	switch caseName {
	case CASE_NORMAL:
		return &omsclient.CreateOrderResponse{
			{
				Data: omsclient.OmsOrder{
					ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID001",
					ChannelCode:             "TESTCHANNELCODE",
				},
				Success: true,
				Message: "Success create order.",
			},
			{
				Data: omsclient.OmsOrder{
					ChannelOrderIncrementId: "TESTCHANNELORDERINCREMENTID002",
					ChannelCode:             "TESTCHANNELCODE",
				},
				Success: true,
				Message: "Success create order.",
			},
		}
	case CASE_NIL_RESPONSE:
		return nil
	}

	return &omsclient.CreateOrderResponse{}
}
