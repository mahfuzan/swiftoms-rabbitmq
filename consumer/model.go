package consumer

type NewOrderQueueMessage struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ChannelOrderIncrementId string `json:"channel_order_increment_id"`
	ChannelCode             string `json:"channel_code"`
}
