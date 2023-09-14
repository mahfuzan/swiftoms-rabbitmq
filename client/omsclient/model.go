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

type CreateInvoiceRequest struct {
	Invoice []OmsOrder `json:"invoice"`
}

type CreateInvoiceResponse []struct {
	Data    Invoice `json:"data"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

type Invoice struct {
	ChannelOrderIncrementId                   string  `json:"channel_order_increment_id"`
	ChannelCode                               string  `json:"channel_code"`
	CompanyCode                               string  `json:"company_code"`
	Items                                     []Item  `json:"items"`
	BaseCurrencyCode                          string  `json:"base_currency_code"`
	BaseDiscountAmount                        float64 `json:"base_discount_amount"`
	BaseGrandTotal                            float64 `json:"base_grand_total"`
	BaseDiscountTaxCompensationAmount         float64 `json:"base_discount_tax_compensation_amount"`
	BaseShippingAmount                        float64 `json:"base_shipping_amount"`
	BaseShippingDiscountTaxCompensationAmount float64 `json:"base_shipping_discount_tax_compensation_amnt"`
	BaseShippingIncludeTax                    float64 `json:"base_shipping_incl_tax"`
	BaseShippingTaxAmount                     float64 `json:"base_shipping_tax_amount"`
	BaseSubtotal                              float64 `json:"base_subtotal"`
	BaseSubtotalIncludeTax                    float64 `json:"base_subtotal_incl_tax"`
	BaseTaxAmount                             float64 `json:"base_tax_amount"`
}

type Item struct {
}
