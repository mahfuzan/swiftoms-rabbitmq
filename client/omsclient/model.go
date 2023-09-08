package omsclient

type CreateOrderRequest struct {
	Orders []OmsOrder `json:"orders"`
}

type OmsOrder struct {
	ChannelOrderIncrementId string `json:"channel_order_increment_id"`
	ChannelCode             string `json:"channel_code"`
}

type CreateOrderResponse []struct {
	Data    OmsOrder `json:"data"`
	Success bool     `json:"success"`
	Message string   `json:"message"`
}
